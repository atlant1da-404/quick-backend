package cache

import "github.com/dgraph-io/ristretto"

type Repository struct {
	r *ristretto.Cache
}

// New creates a new Dragonfly repository.
func New(r *ristretto.Cache) *Repository {
	return &Repository{r: r}
}
