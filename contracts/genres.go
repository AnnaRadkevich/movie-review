package contracts

type Genre struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

type GetGenreRequest struct {
	GenreId int `param:"genreId" validate:"nonzero"`
}

type CreateGenreRequest struct {
	Name string `json:"Name" validate:"min=3,max=32" `
}

type UpdateGenreRequest struct {
	GenreId int    `param:"genreId" validate:"nonzero"`
	Name    string `json:"Name" validate:"min=3,max=32" `
}

type DeleteGenreRequest struct {
	GenreId int `param:"genreId" validate:"nonzero"`
}
