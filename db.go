package catalogue

import(
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	GetSock(ctx context.Context, id string) (Sock, error)
	GetSocks(ctx context.Context, tags []string, order string) ([]Sock, error)
	GetTags(context.Context) ([]string, error)
	CountSocks(ctx context.Context, tags []string) (int, error)
	Ping() error
	Close() error
}

type SqlxDb interface {
	Ping() error
	Close() error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Prepare(ctx context.Context, query string) (StmtMiddleware, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// SqlxDb meets the Database interface requirements
type db struct {
	// db is a reference for the underlying database implementation
	//db     *sqlx.DB
	db     SqlxDb
	logger log.Logger
}

type sqlxdb struct {
	db *sqlx.DB
}

func NewDatabase(sdb *sqlx.DB, logger log.Logger) Database {
	return &db {
		db:     SqlxDbTracingMiddleware()(&sqlxdb{db: sdb}),
		logger: logger,
	}
}

// ErrNotFound is returned when there is no sock for a given ID.
var ErrNotFound = errors.New("not found")

// ErrDBConnection is returned when connection with the database fails.
var ErrDBConnection = errors.New("database connection error")

var baseQuery = "SELECT sock.sock_id AS id, sock.name, sock.description, sock.price, sock.count, sock.image_url_1, sock.image_url_2, GROUP_CONCAT(tag.name) AS tag_name FROM sock JOIN sock_tag ON sock.sock_id=sock_tag.sock_id JOIN tag ON sock_tag.tag_id=tag.tag_id"

func (db *db) GetSock(ctx context.Context, id string) (Sock, error) {
	query := baseQuery + " WHERE sock.sock_id =? GROUP BY sock.sock_id;"

	var sock Sock
	err := db.db.Get(ctx, &sock, query, id)
	if err != nil {
		db.logger.Log("database error", err)
		return Sock{}, ErrNotFound
	}

	sock.ImageURL = []string{sock.ImageURL_1, sock.ImageURL_2}
	sock.Tags = strings.Split(sock.TagString, ",")
	return sock, nil
}

func (db *db) GetSocks(ctx context.Context, tags []string, order string) ([]Sock, error) {
	var socks []Sock
	query := baseQuery

	var args []interface{}

	for i, t := range tags {
		if i == 0 {
			query += " WHERE tag.name=?"
			args = append(args, t)
		} else {
			query += " OR tag.name=?"
			args = append(args, t)
		}
	}

	query += " GROUP BY id"

	if order != "" {
		query += " ORDER BY ?"
		args = append(args, order)
	}

	query += ";"

	err := db.db.Select(ctx, &socks, query, args...)
	if err != nil {
		db.logger.Log("database error", err)
		return []Sock{}, ErrDBConnection
	}
	for i, s := range socks {
		socks[i].ImageURL = []string{s.ImageURL_1, s.ImageURL_2}
		socks[i].Tags = strings.Split(s.TagString, ",")
	}

	return socks, nil
}

func (db *db) CountSocks(ctx context.Context, tags []string) (int, error) {
	query := "SELECT COUNT(DISTINCT sock.sock_id) FROM sock JOIN sock_tag ON sock.sock_id=sock_tag.sock_id JOIN tag ON sock_tag.tag_id=tag.tag_id"

	var args []interface{}

	for i, t := range tags {
		if i == 0 {
			query += " WHERE tag.name=?"
			args = append(args, t)
		} else {
			query += " OR tag.name=?"
			args = append(args, t)
		}
	}

	query += ";"

	sel, err := db.db.Prepare(ctx, query)

	if err != nil {
		db.logger.Log("database error", err)
		return 0, ErrDBConnection
	}
	defer sel.Close()

	var count int
	err = sel.QueryRow(ctx, args...).Scan(&count)

	if err != nil {
		db.logger.Log("database error", err)
		return 0, ErrDBConnection
	}
	return count, nil
}

func (db *db) GetTags(ctx context.Context) ([]string, error) {
	var tags []string
	query := "SELECT name FROM tag;"
	rows, err := db.db.Query(ctx, query)
	if err != nil {
		return []string{}, ErrDBConnection
	}
	var tag string
	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			db.logger.Log("database error", err)
			continue
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (db *db) Ping() error {
	return db.db.Ping()
}

func (db *db) Close() error {
	return db.db.Close()
}

func (sqlxdb *sqlxdb) Ping() error {
	return sqlxdb.db.Ping()
}

func (sqlxdb *sqlxdb) Close() error {
	return sqlxdb.db.Close()
}

func (sqlxdb *sqlxdb) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return sqlxdb.db.Select(dest, query, args...)
}

func (sqlxdb *sqlxdb) Prepare(ctx context.Context, query string) (StmtMiddleware, error) {
	sel, err := sqlxdb.db.Prepare(query)
	return StmtMiddleware{next: sel}, err
}

func (sqlxdb *sqlxdb) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return sqlxdb.db.Get(dest, query, args...)
}

func (sqlxdb *sqlxdb) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return sqlxdb.db.Query(query, args...)
}
