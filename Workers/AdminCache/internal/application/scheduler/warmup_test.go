package scheduler

import (
	"testing"

	"github.com/stretchr/testify/require"

	"admincache/internal/application/ranking"
)

func TestNormalizeRankingLimit(t *testing.T) {
	items := []ranking.RankItem{
		{UserID: 1, Username: "alice", Votes: 50},
		{UserID: 2, Username: "bob", Votes: 75},
		{UserID: 3, Username: "carol", Votes: 25},
	}

	normalized, err := normalizeRanking(items, 2)
	require.NoError(t, err)
	require.Len(t, normalized, 2)
	require.Equal(t, int64(2), normalized[0].UserID)
	require.Equal(t, 1, normalized[0].Rank)
	require.Equal(t, int64(1), normalized[1].UserID)
	require.Equal(t, 2, normalized[1].Rank)
}

func TestNormalizeRankingDetectsDuplicateUser(t *testing.T) {
	items := []ranking.RankItem{
		{UserID: 1, Username: "alice", Votes: 50},
		{UserID: 1, Username: "alice", Votes: 30},
	}

	_, err := normalizeRanking(items, 10)
	require.Error(t, err)
}

func TestNormalizeRankingDetectsDuplicateUsernameWhenNoID(t *testing.T) {
	items := []ranking.RankItem{
		{UserID: 0, Username: "anon", Votes: 50},
		{UserID: 0, Username: "anon", Votes: 40},
	}

	_, err := normalizeRanking(items, 10)
	require.Error(t, err)
}
