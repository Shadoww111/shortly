package services

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shortly/internal/models"
	"github.com/shortly/internal/utils"
)

type ClickService struct {
	db  *pgxpool.Pool
	geo *GeoService
}

func NewClickService(db *pgxpool.Pool, geo *GeoService) *ClickService {
	return &ClickService{db: db, geo: geo}
}

func (s *ClickService) Record(ctx context.Context, linkID int, ip, userAgent, referer string) error {
	device, browser, os := utils.ParseUserAgent(userAgent)

	var country, city string
	if s.geo != nil {
		if result, err := s.geo.Lookup(ip); err == nil {
			country = result.Country
			city = result.City
		} else {
			log.Printf("geo lookup failed for %s: %v", ip, err)
		}
	}

	_, err := s.db.Exec(ctx,
		`INSERT INTO clicks (link_id, ip_address, user_agent, referer, country, city, device, browser, os)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		linkID, ip, userAgent, referer, country, city, device, browser, os,
	)
	return err
}

func (s *ClickService) GetStats(ctx context.Context, linkID int, days int) (*models.ClickStats, error) {
	stats := &models.ClickStats{}

	s.db.QueryRow(ctx, "SELECT COUNT(*) FROM clicks WHERE link_id=$1", linkID).Scan(&stats.TotalClicks)
	s.db.QueryRow(ctx, "SELECT COUNT(DISTINCT ip_address) FROM clicks WHERE link_id=$1", linkID).Scan(&stats.UniqueClicks)

	rows, _ := s.db.Query(ctx,
		`SELECT DATE(created_at) as day, COUNT(*) FROM clicks
		 WHERE link_id=$1 AND created_at >= NOW() - INTERVAL '1 day' * $2
		 GROUP BY day ORDER BY day`, linkID, days)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var dc models.DayCount
			rows.Scan(&dc.Date, &dc.Count)
			stats.ClicksByDay = append(stats.ClicksByDay, dc)
		}
	}

	stats.TopReferrers = s.getTopN(ctx, linkID, "referer", 10)
	stats.TopCountries = s.getTopN(ctx, linkID, "country", 10)
	stats.TopBrowsers = s.getTopN(ctx, linkID, "browser", 5)
	stats.TopDevices = s.getTopN(ctx, linkID, "device", 5)
	stats.TopOS = s.getTopN(ctx, linkID, "os", 5)

	return stats, nil
}

func (s *ClickService) getTopN(ctx context.Context, linkID int, column string, limit int) []models.NameCount {
	var results []models.NameCount
	rows, err := s.db.Query(ctx,
		`SELECT COALESCE(NULLIF(`+column+`, ''), 'direct') as name, COUNT(*) as count
		 FROM clicks WHERE link_id=$1 GROUP BY name ORDER BY count DESC LIMIT $2`,
		linkID, limit,
	)
	if err != nil {
		return results
	}
	defer rows.Close()
	for rows.Next() {
		var nc models.NameCount
		rows.Scan(&nc.Name, &nc.Count)
		results = append(results, nc)
	}
	return results
}
