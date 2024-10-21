package main

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"oneof=asc desc"`
	Offset int    `json:"offset" validate:"oneof=asc desc"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}
