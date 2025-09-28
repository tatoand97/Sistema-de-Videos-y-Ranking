package ranking

import (
	"context"
	"database/sql"
	"fmt"
)

type RankItem struct {
	UserID   int64  `json:"user_id"`
	Position int    `json:"position"`
	Username string `json:"username"`
	City     string `json:"city"`
	Votes    int64  `json:"votes"`
}

type Computer interface {
	Compute(ctx context.Context, city *string, page, size int) ([]RankItem, error)
}

type computer struct{ db *sql.DB }

func NewRankComputer(db *sql.DB) Computer { return &computer{db: db} }

// SQL alineado con la API:
// - Cuenta votos por usuario sobre videos publicados y procesados (processed_file no nulo)
// - username: split_part(u.email,'@',1)
// - city: c.name
// - Agrupa por u.user_id, u.email, c.name
// - Orden: votes DESC, u.user_id ASC (desempate estable)
const baseSQL = `
SELECT
  u.user_id,
  split_part(u.email,'@',1) AS username,
  COALESCE(c.name, '')      AS city,
  COUNT(vt.vote_id)         AS votes
FROM users u
JOIN video v       ON v.user_id   = u.user_id
LEFT JOIN vote vt  ON vt.video_id = v.video_id
LEFT JOIN city c   ON c.city_id   = u.city_id
WHERE v.status = 'PUBLISHED' AND v.processed_file IS NOT NULL
%s
GROUP BY u.user_id, u.email, c.name
ORDER BY votes DESC, u.user_id ASC
LIMIT $1 OFFSET $2
`

func (c *computer) Compute(ctx context.Context, city *string, page, size int) ([]RankItem, error) {
	offset := (page - 1) * size

	var rows *sql.Rows
	var err error

	if city != nil && *city != "" {
		// Filtro por nombre de ciudad:
		// - Si existe immutable_unaccent(text), lo usamos para ignorar tildes
		// - De lo contrario, fallback a lower() simple
		where := `
 AND (
   to_regprocedure('immutable_unaccent(text)') IS NOT NULL
   AND immutable_unaccent(lower(COALESCE(c.name,''))) = immutable_unaccent(lower($3))
 )
 OR (
   to_regprocedure('immutable_unaccent(text)') IS NULL
   AND lower(COALESCE(c.name,'')) = lower($3)
 )`
		rows, err = c.db.QueryContext(ctx, fmt.Sprintf(baseSQL, where), size, offset, *city)
	} else {
		rows, err = c.db.QueryContext(ctx, fmt.Sprintf(baseSQL, ""), size, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []RankItem
	pos := 0
	for rows.Next() {
		pos++
		var it RankItem
		if err := rows.Scan(&it.UserID, &it.Username, &it.City, &it.Votes); err != nil {
			return nil, err
		}
		it.Position = pos
		res = append(res, it)
	}
	return res, rows.Err()
}
