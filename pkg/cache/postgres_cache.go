//ff:type feature=pkg-cache type=model
//ff:what PostgreSQL 기반 캐시 구조체
package cache

import "database/sql"

type postgresCache struct {
	db *sql.DB
}
