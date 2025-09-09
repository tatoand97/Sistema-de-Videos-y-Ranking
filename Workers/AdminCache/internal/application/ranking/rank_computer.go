package ranking

import (
	"context"
	"database/sql"
	"fmt"
)

type RankItem struct {
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

// Usa solo columnas que existen seguro:
// - "username": derivado de email -> split_part(u.email,'@',1)
// - "city": de city.name vía users.city_id
// Tablas singulares: video / vote (según tu esquema)
// Si quieres filtrar por estado del video, descomenta WHERE vi.status = 'PROCESSED'
const baseSQL = `
SELECT
  split_part(u.email,'@',1) AS username,
  COALESCE(c.name, '')      AS city,
  COUNT(vt.vote_id)         AS votes
FROM video vi
JOIN users u       ON u.user_id   = vi.user_id
LEFT JOIN vote vt  ON vt.video_id = vi.video_id
LEFT JOIN city c   ON c.city_id   = u.city_id
/* WHERE vi.status = 'PROCESSED' */
%s
GROUP BY 1, 2
ORDER BY votes DESC, 1 ASC
LIMIT $1 OFFSET $2
`

func (c *computer) Compute(ctx context.Context, city *string, page, size int) ([]RankItem, error) {
	offset := (page - 1) * size

	var rows *sql.Rows
	var err error

	if city != nil && *city != "" {
		// Filtro por nombre de ciudad usando c.name.
		// Si tienes la extensión unaccent, lo hacemos insensible a acentos.
		where := `
 AND (
   (select count(*) from pg_extension where extname='unaccent') > 0
   AND lower(unaccent(COALESCE(c.name,''))) = lower(unaccent($3))
 )
 OR (
   (select count(*) from pg_extension where extname='unaccent') = 0
   AND lower(COALESCE(c.name,'')) = lower($3)
 )`
		rows, err = c.db.QueryContext(ctx, fmt.Sprintf(baseSQL, where), size, offset, *city)
	} else {
		rows, err = c.db.QueryContext(ctx, fmt.Sprintf(baseSQL, ""), size, offset)
	}
	if err != nil { return nil, err }
	defer rows.Close()

	var res []RankItem
	pos := 0
	for rows.Next() {
		pos++
		var it RankItem
		if err := rows.Scan(&it.Username, &it.City, &it.Votes); err != nil {
			return nil, err
		}
		it.Position = pos
		res = append(res, it)
	}
	return res, rows.Err()
}