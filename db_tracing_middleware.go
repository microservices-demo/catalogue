package catalogue

import (
	"fmt"
	"context"
	"database/sql"

	otext "github.com/opentracing/opentracing-go/ext"
	stdopentracing "github.com/opentracing/opentracing-go"
)

// DbMiddleware decorates a Database.
type DbMiddleware func(Database) Database

// SqlxDbMiddleware decorates a SqlxDb.
type SqlxDbMiddleware func(SqlxDb) SqlxDb

// DbTracingMiddleware returns middleware for tracing app level db access.
func DbTracingMiddleware() DbMiddleware {
	return func(next Database) Database {
		return dbTracingMiddleware{
			next: next,
		}
	}
}

// SqlxDbTracingMiddleware returns middleware for tracing low level db access.
func SqlxDbTracingMiddleware() SqlxDbMiddleware {
	return func(next SqlxDb) SqlxDb {
		return sqlxDbTracingMiddleware{
			next: next,
		}
	}
}

// dbTracingMiddleware meets the Database interface.
type dbTracingMiddleware struct {
	next Database
}

// sqlxDbTracingMiddleware meets the SqlxDb interface.
type sqlxDbTracingMiddleware struct {
	next SqlxDb
}

type StmtMiddleware struct {
	next *sql.Stmt
}

func (stmt StmtMiddleware) Close() error {
	return stmt.next.Close()
}

func (stmt StmtMiddleware) QueryRow(ctx context.Context, args ...interface{}) *sql.Row {
	return stmt.next.QueryRow(args...)
}

func (mw dbTracingMiddleware) GetSock(ctx context.Context, id string) (Sock, error) {
	span, ctx := startSpan(ctx, "sock from database")
	sock, err := mw.next.GetSock(ctx, id)
	finishSpan(span, len(fmt.Sprintf("%#v", sock)))
	return sock, err
}

func (mw dbTracingMiddleware) GetSocks(ctx context.Context, tags []string, order string) ([]Sock, error) {
	span, ctx := startSpan(ctx, "socks from database")
	socks, err := mw.next.GetSocks(ctx, tags, order)
	finishSpan(span, len(fmt.Sprintf("%#v", tags)))
	return socks, err
}

func (mw dbTracingMiddleware) CountSocks(ctx context.Context, tags []string) (int, error) {
	span, ctx := startSpan(ctx, "count socks from database")
	count, err := mw.next.CountSocks(ctx, tags)
	finishSpan(span, 170)
	return count, err
}

func (mw dbTracingMiddleware) GetTags(ctx context.Context) ([]string, error) {
	span, ctx := startSpan(ctx, "tags from database")
	tags, err := mw.next.GetTags(ctx)
	finishSpan(span, len(fmt.Sprintf("%#v", tags)))
	return tags, err
}

func (mw dbTracingMiddleware) Ping() error {
	return mw.next.Ping()
}

func (mw dbTracingMiddleware) Close() error {
	return mw.next.Close()
}

func (mw sqlxDbTracingMiddleware) Ping() error {
	return mw.next.Ping()
}

func (mw sqlxDbTracingMiddleware) Close() error {
	return mw.next.Close()
}

func (mw sqlxDbTracingMiddleware) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	span := stdopentracing.SpanFromContext(ctx)
	span.SetTag("db.query.size", len(query))
	return mw.next.Select(ctx, dest, query, args...)
}

func (mw sqlxDbTracingMiddleware) Prepare(ctx context.Context, query string) (StmtMiddleware, error) {
	span := stdopentracing.SpanFromContext(ctx)
	span.SetTag("db.query.size", len(query))
	return mw.next.Prepare(ctx, query)
}

func (mw sqlxDbTracingMiddleware) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	span := stdopentracing.SpanFromContext(ctx)
	span.SetTag("db.query.size", len(query) + len(fmt.Sprintf("%#v", args)))
	return mw.next.Get(ctx, dest, query, args...)
}

func (mw sqlxDbTracingMiddleware) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span := stdopentracing.SpanFromContext(ctx)
	span.SetTag("db.query.size", len(query) + len(fmt.Sprintf("%#v", args)))
	return mw.next.Query(ctx, query, args...)
}

func startSpan(ctx context.Context, n string) (stdopentracing.Span, context.Context) {
	var span stdopentracing.Span
	span, ctx = stdopentracing.StartSpanFromContext(ctx, n)
	otext.SpanKindRPCClient.Set(span)
	span.SetTag("db.type", "mysql")
	span.SetTag("peer.address", "catalogue-db:3306")
	return span, ctx
}

func finishSpan(span stdopentracing.Span, size int) {
	span.SetTag("db.query.result.size", size)
	span.Finish()
}
