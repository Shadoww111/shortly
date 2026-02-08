package services

// This file documents the fix for expired links being served from cache.
//
// Problem: When a link expires, the cached version in Redis still serves
// the redirect until the TTL expires (up to 24h after the link expired).
//
// Fix: After cache hit, verify the link is still active and not expired
// by checking the database. If invalid, delete the cache entry.
//
// The fix is applied in LinkService.Resolve():
//   1. Cache hit -> quick DB check for is_active + expires_at
//   2. If invalid -> delete cache key + return error
//   3. If valid -> proceed with redirect
//
// Additionally: when a link is deactivated or deleted, the cache key
// is now explicitly removed.
