package catalogue

import (
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware decorates a service.
type Middleware func(Service) Service

// LoggingMiddleware logs method calls, parameters, results, and elapsed time.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) List(tags []string, order string, pageNum, pageSize int) (socks []Sock, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "List",
			"tags", strings.Join(tags, ", "),
			"order", order,
			"pageNum", pageNum,
			"pageSize", pageSize,
			"result", len(socks),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.List(tags, order, pageNum, pageSize)
}

func (mw loggingMiddleware) Count(tags []string) (n int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Count",
			"tags", strings.Join(tags, ", "),
			"result", n,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Count(tags)
}

func (mw loggingMiddleware) Get(id string) (s Sock, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Get",
			"id", id,
			"sock", s.ID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Get(id)
}

func (mw loggingMiddleware) Tags() (tags []string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Tags",
			"result", len(tags),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Tags()
}

func (mw loggingMiddleware) Health() (health []Health) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Health",
			"result", len(health),
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Health()
}

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		Service:        s,
	}
}

func (s *instrumentingService) List(tags []string, order string, pageNum, pageSize int) ([]Sock, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list").Add(1)
		s.requestLatency.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.List(tags, order, pageNum, pageSize)
}

func (s *instrumentingService) Count(tags []string) (int, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "count").Add(1)
		s.requestLatency.With("method", "count").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Count(tags)
}

func (s *instrumentingService) Get(id string) (Sock, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get").Add(1)
		s.requestLatency.With("method", "get").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Get(id)
}

func (s *instrumentingService) Tags() ([]string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "tags").Add(1)
		s.requestLatency.With("method", "tags").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.Tags()
}

func (s *instrumentingService) Health() []Health {
	defer func(begin time.Time) {
		s.requestCount.With("method", "health").Add(1)
		s.requestLatency.With("method", "health").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Health()
}
