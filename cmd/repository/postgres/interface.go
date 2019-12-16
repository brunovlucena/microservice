package postgres

import (
	"errors"
	"net/url"
	"strings"

	. "github.com/brunovlucena/microservice/cmd/data"
)

//Repository repository interface
type Repository interface {
	Create(config *Config) (*Config, error)
	Find(name string) (*Config, error)
	FindAll() ([]*Config, error)
	Update(config *Config) (*Config, error)
	Remove(name string) (*Config, error)
	Search(params url.Values) ([]*Config, error)
}

// NewRepository returns a selected repository.
func NewRepository(name, host, port, user, pass, dbname string) (Repository, error) {
	switch strings.ToLower(name) {
	case "postgres":
		return NewPostgres(host, port, user, pass, dbname)
	}
	return nil, errors.New("Invalid base given")
}
