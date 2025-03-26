# shortly

url shortener with click analytics, qr code generation, and link management. built with go.

## features

- **shorten urls** — random or custom short codes
- **click tracking** — ip, user agent, referer, device, browser, os
- **analytics** — clicks by day, top referrers, country breakdown, device stats
- **qr codes** — generate png qr codes for any short link
- **link management** — expiration dates, max click limits, tags
- **auth** — jwt with rate-limited register/login
- **caching** — redis for fast redirects
- **pagination** — paginated link listing

## tech stack

| | |
|---|---|
| language | go 1.22 |
| router | chi v5 |
| database | postgresql 16 |
| cache | redis 7 |
| auth | jwt (golang-jwt) + bcrypt |
| qr | go-qrcode |
| container | docker + multi-stage build |

## setup

```bash
# start postgres + redis
docker-compose up -d db redis

# run
cp .env.example .env
make run
```

or full docker:

```bash
make docker
```

api on http://localhost:8080

## api

### auth
| method | route | description |
|--------|-------|-------------|
| POST | /api/auth/register | register (rate limited) |
| POST | /api/auth/login | login |

### links (auth required)
| method | route | description |
|--------|-------|-------------|
| POST | /api/links | create short link |
| GET | /api/links | list your links (paginated) |
| DELETE | /api/links/{id} | delete link |
| GET | /api/links/{id}/stats | click analytics |

### public
| method | route | description |
|--------|-------|-------------|
| GET | /{code} | redirect to original url |
| GET | /qr/{code}?size=256 | get qr code png |

### create link payload

```json
{
  "url": "https://example.com/very/long/path",
  "title": "my link",
  "custom_code": "mylink",
  "expires_in": 30,
  "max_clicks": 1000,
  "tags": ["marketing", "social"]
}
```

### analytics response

```json
{
  "total_clicks": 1523,
  "unique_clicks": 891,
  "clicks_by_day": [{"date": "2025-03-20", "count": 45}],
  "top_referrers": [{"name": "twitter.com", "count": 320}],
  "top_countries": [{"name": "US", "count": 612}],
  "top_browsers": [{"name": "Chrome", "count": 890}],
  "top_devices": [{"name": "mobile", "count": 723}]
}
```

## license

MIT
