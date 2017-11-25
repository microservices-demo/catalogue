package catalogue

import(
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Database interface {
	Close() error
	Ping() error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Prepare(query string) (StmtMiddleware, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// SqlxDb meets the Database interface requirements
type SqlxDb struct {
	// db is a reference for the underlying database implementation
	Db *sqlx.DB
}

func (sqlxdb *SqlxDb) Close() error {
	return sqlxdb.Db.Close()
}

func (sqlxdb *SqlxDb) Ping() error {
	return sqlxdb.Db.Ping()
}

func (sqlxdb *SqlxDb) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return sqlxdb.Db.Select(dest, query, args...)
}

func (sqlxdb *SqlxDb) Prepare(query string) (StmtMiddleware, error) {
	sel, err := sqlxdb.Db.Prepare(query)
	return StmtMiddleware{next: sel}, err
}

func (sqlxdb *SqlxDb) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return sqlxdb.Db.Get(dest, query, args...)
}

func (sqlxdb *SqlxDb) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return sqlxdb.Db.Query(query, args...)
}
