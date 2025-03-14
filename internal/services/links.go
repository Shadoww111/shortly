package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shortly/internal/cache"
	"github.com/shortly/internal/config"
	"github.com/shortly/internal/models"
	"github.com/shortly/internal/utils"
)

type LinkService struct {
	db    *pgxpool.Pool
	cache *cache.RedisCache
	cfg   *config.Config
}

func NewLinkService(db *pgxpool.Pool, cache *cache.RedisCache, cfg *config.Config) *LinkService {
	return &LinkService{db: db, cache: cache, cfg: cfg}
}

func (s *LinkService) Create(ctx context.Context, userID int, req models.CreateLinkRequest) (*models.Link, error) {
	if !utils.IsValidURL(req.URL) {
		return nil, errors.New("invalid url")
	}

	var code string
	var err error

	if req.CustomCode != "" {
		if !utils.IsValidCustomCode(req.CustomCode) {
			return nil, errors.New("invalid custom code: 3-20 alphanumeric chars, hyphens, underscores")
		}
		// check if taken
		var exists bool
		s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM links WHERE short_code=$1)", req.CustomCode).Scan(&exists)
		if exists {
			return nil, errors.New("short code already taken")
		}
		code = req.CustomCode
	} else {
		for i := 0; i < 5; i++ {
			code, err = utils.GenerateShortCode(s.cfg.ShortCodeLength)
			if err != nil {
				return nil, err
			}
			var exists bool
			s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM links WHERE short_code=$1)", code).Scan(&exists)
			if !exists {
				break
			}
			if i == 4 {
				return nil, errors.New("could not generate unique code")
			}
		}
	}

	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(req.ExpiresIn) * 24 * time.Hour)
		expiresAt = &t
	}

	link := &models.Link{}
	err = s.db.QueryRow(ctx,
		`INSERT INTO links (short_code, original_url, title, user_id, expires_at, max_clicks)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, short_code, original_url, title, user_id, is_active, expires_at, max_clicks, created_at, updated_at`,
		code, req.URL, req.Title, userID, expiresAt, req.MaxClicks,
	).Scan(&link.ID, &link.ShortCode, &link.OriginalURL, &link.Title, &link.UserID,
		&link.IsActive, &link.ExpiresAt, &link.MaxClicks, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert link: %w", err)
	}

	link.ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, link.ShortCode)

	// handle tags
	if len(req.Tags) > 0 {
		for _, tagName := range req.Tags {
			var tagID int
			err := s.db.QueryRow(ctx,
				`INSERT INTO tags (name, user_id) VALUES ($1, $2) ON CONFLICT (name, user_id) DO UPDATE SET name=$1 RETURNING id`,
				tagName, userID,
			).Scan(&tagID)
			if err != nil {
				continue
			}
			s.db.Exec(ctx, "INSERT INTO link_tags (link_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", link.ID, tagID)
		}
	}

	// cache the redirect
	_ = s.cache.Set(ctx, "link:"+code, link.OriginalURL, 24*time.Hour)

	return link, nil
}

func (s *LinkService) Resolve(ctx context.Context, code string) (string, int, error) {
	// try cache first
	var cachedURL string
	if err := s.cache.Get(ctx, "link:"+code, &cachedURL); err == nil {
		// get link id for click tracking
		var linkID int
		s.db.QueryRow(ctx, "SELECT id FROM links WHERE short_code=$1", code).Scan(&linkID)
		return cachedURL, linkID, nil
	}

	var link models.Link
	err := s.db.QueryRow(ctx,
		"SELECT id, original_url, is_active, expires_at, max_clicks FROM links WHERE short_code=$1",
		code,
	).Scan(&link.ID, &link.OriginalURL, &link.IsActive, &link.ExpiresAt, &link.MaxClicks)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", 0, errors.New("not found")
		}
		return "", 0, err
	}

	if !link.IsActive {
		return "", 0, errors.New("link disabled")
	}
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return "", 0, errors.New("link expired")
	}
	if link.MaxClicks != nil {
		var count int
		s.db.QueryRow(ctx, "SELECT COUNT(*) FROM clicks WHERE link_id=$1", link.ID).Scan(&count)
		if count >= *link.MaxClicks {
			return "", 0, errors.New("click limit reached")
		}
	}

	_ = s.cache.Set(ctx, "link:"+code, link.OriginalURL, 24*time.Hour)
	return link.OriginalURL, link.ID, nil
}

func (s *LinkService) ListByUser(ctx context.Context, userID, page, perPage int) (*models.LinkListResponse, error) {
	offset := (page - 1) * perPage

	var total int
	s.db.QueryRow(ctx, "SELECT COUNT(*) FROM links WHERE user_id=$1", userID).Scan(&total)

	rows, err := s.db.Query(ctx,
		`SELECT l.id, l.short_code, l.original_url, l.title, l.user_id, l.is_active,
		        l.expires_at, l.max_clicks, l.created_at, l.updated_at,
		        COALESCE((SELECT COUNT(*) FROM clicks WHERE link_id=l.id), 0) as click_count
		 FROM links l WHERE l.user_id=$1 ORDER BY l.created_at DESC LIMIT $2 OFFSET $3`,
		userID, perPage, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		err := rows.Scan(&l.ID, &l.ShortCode, &l.OriginalURL, &l.Title, &l.UserID,
			&l.IsActive, &l.ExpiresAt, &l.MaxClicks, &l.CreatedAt, &l.UpdatedAt, &l.ClickCount)
		if err != nil {
			continue
		}
		l.ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, l.ShortCode)
		links = append(links, l)
	}

	return &models.LinkListResponse{Links: links, Total: total, Page: page, PerPage: perPage}, nil
}

func (s *LinkService) Delete(ctx context.Context, linkID, userID int) error {
	res, err := s.db.Exec(ctx, "DELETE FROM links WHERE id=$1 AND user_id=$2", linkID, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}
