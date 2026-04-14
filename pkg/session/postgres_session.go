//ff:type feature=pkg-session type=model
//ff:what PostgreSQL 기반 세션 구조체
package session

import "database/sql"

type postgresSession struct {
	db *sql.DB
}
