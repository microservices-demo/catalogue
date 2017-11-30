package catalogue

// service.go contains the definition and implementation (business logic) of the
// catalogue service. Everything here is agnostic to the transport (HTTP).

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Service is the catalogue service, providing read operations on a saleable
// catalogue of sock products.
type Service interface {
	List(ctx context.Context, tags []string, order string, pageNum, pageSize int) ([]Sock, error) // GET /catalogue
	Count(ctx context.Context, tags []string) (int, error)                                        // GET /catalogue/size
	Get(ctx context.Context, id string) (Sock, error)                                             // GET /catalogue/{id}
	Tags(ctx context.Context) ([]string, error)                                                 // GET /tags
	Health() []Health                                                        // GET /health
}

// Middleware decorates a Service.
type Middleware func(Service) Service

// Sock describes the thing on offer in the catalogue.
type Sock struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	ImageURL    []string `json:"imageUrl" db:"-"`
	ImageURL_1  string   `json:"-" db:"image_url_1"`
	ImageURL_2  string   `json:"-" db:"image_url_2"`
	Price       float32  `json:"price" db:"price"`
	Count       int      `json:"count" db:"count"`
	Tags        []string `json:"tag" db:"-"`
	TagString   string   `json:"-" db:"tag_name"`
}

// Health describes the health of a service
type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

// NewCatalogueService returns an implementation of the Service interface,
// with connection to an SQL database.
func NewCatalogueService(db Database, logger log.Logger) Service {
	return &catalogueService{
		db:     db,
		logger: logger,
	}
}

type catalogueService struct {
	db     Database
	logger log.Logger
}

func (s *catalogueService) List(ctx context.Context, tags []string, order string, pageNum, pageSize int) ([]Sock, error) {
	socks, err := s.db.GetSocks(ctx, tags, order)
	if err != nil {
		return []Sock{}, ErrDBConnection
	}

	// DEMO: Change 0 to 850
	time.Sleep(0 * time.Millisecond)

	socks = cut(socks, pageNum, pageSize)

	return socks, nil
}

func (s *catalogueService) Count(ctx context.Context, tags []string) (int, error) {
	count, err := s.db.CountSocks(ctx, tags)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *catalogueService) Get(ctx context.Context, id string) (Sock, error) {
	sock, err := s.db.GetSock(ctx, id)
	if err != nil {
		return Sock{}, err
	}

	return sock, nil
}

func (s *catalogueService) Health() []Health {
	var health []Health
	dbstatus := "OK"

	err := s.db.Ping()
	if err != nil {
		dbstatus = "err"
	}

	app := Health{"catalogue", "OK", time.Now().String()}
	db := Health{"catalogue-db", dbstatus, time.Now().String()}

	health = append(health, app)
	health = append(health, db)

	return health
}

func (s *catalogueService) Tags(ctx context.Context) ([]string, error) {
	tags, err := s.db.GetTags(ctx)
	if err != nil {
		return []string{}, err
	}
	return tags, nil
}

func cut(socks []Sock, pageNum, pageSize int) []Sock {
	if pageNum == 0 || pageSize == 0 {
		return []Sock{} // pageNum is 1-indexed
	}
	start := (pageNum * pageSize) - pageSize
	if start > len(socks) {
		return []Sock{}
	}
	end := (pageNum * pageSize)
	if end > len(socks) {
		end = len(socks)
	}
	return socks[start:end]
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
