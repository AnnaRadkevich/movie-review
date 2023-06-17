package movies

import (
	"github.com/cloudmachinery/movie-reviews/internal/config"
	"github.com/cloudmachinery/movie-reviews/internal/modules/genres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewModule(db *pgxpool.Pool, paginationConfig config.PaginationConfig, genresModule *genres.Module) *Module {
	repo := NewRepository(db, genresModule.Repository)
	service := NewService(repo, genresModule.Service)
	handler := NewHandler(service, paginationConfig)
	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repo,
	}
}
