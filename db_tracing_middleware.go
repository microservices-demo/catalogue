package catalogue

import (
	"context"
	"database/sql"
	"unsafe"

	otext "github.com/opentracing/opentracing-go/ext"
	stdopentracing "github.com/opentracing/opentracing-go"
)

// Middleware decorates a database.
type DbMiddleware func(Database) Database

// DbTracingMiddleware traces database calls.
func DbTracingMiddleware() DbMiddleware {
	return func(next Database) Database {
		return dbTracingMiddleware{
			next: next,
		}
	}
}

type dbTracingMiddleware struct {
	next Database
}

type StmtMiddleware struct {
	next *sql.Stmt
}

func (stmt StmtMiddleware) Close() error {
	return stmt.next.Close()
}

func (stmt StmtMiddleware) QueryRow(ctx context.Context, args ...interface{}) *sql.Row {
	span := startSpan(ctx, "rows from database")
	rows := stmt.next.QueryRow(args...)
	finishSpan(span, unsafe.Sizeof(rows))
	return rows
}

func (mw dbTracingMiddleware) Close() error {
	return mw.next.Close()
}

func (mw dbTracingMiddleware) Ping() error {
	return mw.next.Ping()
}

func (mw dbTracingMiddleware) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	span := startSpan(ctx, "socks from database")
	err := mw.next.Select(ctx, dest, query, args...)
	finishSpan(span, unsafe.Sizeof(dest))
	return err
}

func (mw dbTracingMiddleware) Prepare(query string) (StmtMiddleware, error) {
	return mw.next.Prepare(query)
}

func (mw dbTracingMiddleware) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	span := startSpan(ctx, "get from database")
	err := mw.next.Get(ctx, dest, query, args...)
	finishSpan(span, unsafe.Sizeof(dest))
	return err
}

func (mw dbTracingMiddleware) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span := startSpan(ctx, "query from database")
	rows, err := mw.next.Query(ctx, query, args...)
	finishSpan(span, unsafe.Sizeof(rows))
	return rows, err
}

func startSpan(ctx context.Context, n string) stdopentracing.Span {
	var span stdopentracing.Span
	span, ctx = stdopentracing.StartSpanFromContext(ctx, n)
	otext.SpanKindRPCClient.Set(span)
	span.SetTag("db.type", "mysql")
	span.SetTag("peer.address", "catalogue-db:3306")
	return span
}

func finishSpan(span stdopentracing.Span, size uintptr) {
	span.SetTag("db.query.result.size", size)
	span.Finish()
}
