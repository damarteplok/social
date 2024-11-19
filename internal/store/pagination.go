package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedQuery struct {
	Limit  int    `json:"limit" validate:"required,gte=1,lte=150"`
	Page   int    `json:"page" validate:"required,gte=1"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	Search string `json:"search" validate:"max=100"`
	Since  string `json:"since"`
	Until  string `json:"until"`
}

type PaginatedFeedQuery struct {
	PaginatedQuery
	Tags []string `json:"tags" validate:"max=5"`
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}

func (pq *PaginatedQuery) Parse(r *http.Request) error {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return err
		}
		pq.Limit = l
	}

	page := qs.Get("page")
	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			return err
		}
		pq.Page = p
	}

	pq.Offset = pq.Limit * (pq.Page - 1)

	sort := qs.Get("sort")
	if sort != "" {
		pq.Sort = sort
	}

	search := qs.Get("search")
	if search != "" {
		pq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		pq.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		pq.Until = parseTime(until)
	}

	return nil
}

func (pfq *PaginatedFeedQuery) Parse(r *http.Request) error {
	if err := pfq.PaginatedQuery.Parse(r); err != nil {
		return err
	}

	qs := r.URL.Query()

	tags := qs.Get("tags")
	if tags != "" {
		pfq.Tags = strings.Split(tags, ",")
	}

	return nil
}
