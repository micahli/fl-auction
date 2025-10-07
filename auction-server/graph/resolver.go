// graph/resolver.go
package graph

import (
	"github.com/micahli/fl-auction/auction-server/internal/service"
	"github.com/micahli/fl-auction/auction-server/internal/store"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service *service.AuctionService
	store   *store.AuctionStore
}

// NewResolver creates a new root resolver
func NewResolver(svc *service.AuctionService, st *store.AuctionStore) *Resolver {
	return &Resolver{
		service: svc,
		store:   st,
	}
}
