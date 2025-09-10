package interfaces

import (
	"context"
)

// LeaderboardEntry represents a member and its score in a leaderboard.
type LeaderboardEntry struct {
	Member string
	Score  int64
}

// Stats aggregates for a pollId.
type Stats struct {
	Total   int64
	LastTS  int64
	Version int64
	Uniques int64 // approximate via HLL if enabled
}

// Aggregates defines operations to update/read leaderboards and stats in Redis.
// Keys follow the scheme using hash-tags for cluster co-location:
//   - ZSET:   lb:{votes:<pollId>}           (member=id, score=exact total)
//   - HASH:   stats:{votes:<pollId>}        (fields: total, version)
//   - STRING: stats:{votes:<pollId>}.last_ts
//   - HLL:    hll:{votes:<pollId>}          (optional)
type Aggregates interface {
	// UpdateAfterVote applies the standard updates after a successful DB write.
	// It increments the leaderboard score for member, increments total, sets last_ts,
	// optionally adds to HLL for uniques and increments a version field.
	UpdateAfterVote(ctx context.Context, pollId, member string, uniqueMember *string, tsUnix int64, incVersion bool) error

	// GetLeaderboard returns a slice of entries from start..stop (inclusive), ordered desc by score.
	GetLeaderboard(ctx context.Context, pollId string, start, stop int64) ([]LeaderboardEntry, error)

	// GetStats returns aggregated totals and uniques (if HLL present).
	GetStats(ctx context.Context, pollId string) (Stats, error)

	// GetScore returns the score for a given member in a poll. If absent, 0 and nil error.
	GetScore(ctx context.Context, pollId, member string) (int64, error)
}
