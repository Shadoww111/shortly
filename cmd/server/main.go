package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"

	"github.com/shortly/internal/cache"
	"github.com/shortly/internal/config"
	"github.com/shortly/internal/database"
	"github.com/shortly/internal/handlers"
	"github.com/shortly/internal/middleware"
	"github.com/shortly/internal/services"
)

func main() {
	cfg := config.Load()

	// database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("db:", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("migrations:", err)
	}
	log.Println("db connected + migrated")

	// redis
	rdb, err := cache.NewRedisCache(cfg.RedisURL)
	if err != nil {
		log.Println("warning: redis unavailable, running without cache:", err)
		rdb = nil
	} else {
		defer rdb.Close()
		log.Println("redis connected")
	}

	// services
	authSvc := services.NewAuthService(db, cfg.JWTSecret)
	linkSvc := services.NewLinkService(db, rdb, cfg)
	clickSvc := services.NewClickService(db)

	// handlers
	authH := handlers.NewAuthHandler(authSvc)
	linkH := handlers.NewLinkHandler(linkSvc, clickSvc)
	qrH := handlers.NewQRHandler(cfg)

	// router
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// public
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/{code}", linkH.Redirect)
	r.Get("/qr/{code}", qrH.Generate)

	// auth
	r.Route("/api/auth", func(r chi.Router) {
		r.Use(httprate.LimitByIP(20, time.Minute))
		r.Post("/register", authH.Register)
		r.Post("/login", authH.Login)
	})

	// protected
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))

		r.Post("/links", linkH.Create)
		r.Get("/links", linkH.List)
		r.Delete("/links/{id}", linkH.Delete)
		r.Get("/links/{id}/stats", linkH.GetStats)
	})

	log.Printf("shortly running on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
