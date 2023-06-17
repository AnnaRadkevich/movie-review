package genres

import "github.com/cloudmachinery/movie-reviews/internal/dbx"

type Genre struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

var _ dbx.Keyer = MovieGenreRelation{}

type MovieGenreRelation struct {
	MovieID int
	GenreID int
	OrderNo int
}

func (m MovieGenreRelation) Key() any {
	type MovieGenreRelationKey struct {
		MovieID, GenreID int
	}
	return MovieGenreRelationKey{
		MovieID: m.MovieID,
		GenreID: m.GenreID,
	}
}
