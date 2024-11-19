package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/damarteplok/social/internal/store"
)

func (app *application) handleRequestError(w http.ResponseWriter, r *http.Request, err error) {
	switch err {
	case store.ErrNotFound:
		app.notFoundResponse(w, r, err)
	case store.ErrBadRequest:
		app.badRequestResponse(w, r, err)
	case store.ErrMethodNotAllowed:
		app.methodNotAllowedResponse(w, r, err)
	default:
		app.internalServerError(w, r, err)
	}
}

func StringPtr(s string) *string {
	return &s
}

func GetPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}

func GetUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

func getFlowNodeQueryParams(r *http.Request) (*FlowNodeQueryParams, error) {
	size := r.URL.Query().Get("size")
	if size == "" {
		size = "50"
	}

	order := r.URL.Query().Get("order")
	if order == "" {
		order = "DESC"
	} else {
		order = strings.ToUpper(order)
		if order != "DESC" && order != "ASC" {
			return nil, fmt.Errorf("invalid order value")
		}
	}

	sort := r.URL.Query().Get("sort")

	typeStr := r.URL.Query().Get("type")

	state := r.URL.Query().Get("state")

	searchAfter := r.URL.Query().Get("searchAfter")
	if searchAfter != "" {
		searchAfter = fmt.Sprintf("[%s]", searchAfter)
	}

	searchBefore := r.URL.Query().Get("searchBefore")
	if searchBefore != "" {
		searchBefore = fmt.Sprintf("[%s]", searchBefore)
	}

	return &FlowNodeQueryParams{
		Size:         size,
		Order:        order,
		Sort:         sort,
		SearchAfter:  searchAfter,
		SearchBefore: searchBefore,
		Type:         typeStr,
		State:        state,
	}, nil
}

func getTaskListQueryParams(r *http.Request) (*TaskListQueryParams, error) {
	sizeStr := r.URL.Query().Get("size")
	var size int32
	if sizeStr == "" {
		size = 50
	} else {
		size64, err := strconv.ParseInt(sizeStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid size value")
		}
		size = int32(size64)
	}

	order := r.URL.Query().Get("order")
	if order == "" {
		order = "DESC"
	} else {
		order = strings.ToUpper(order)
		if order != "DESC" && order != "ASC" {
			return nil, fmt.Errorf("invalid order value")
		}
	}

	sort := r.URL.Query().Get("sort")

	state := r.URL.Query().Get("state")

	taskDefinitionId := r.URL.Query().Get("taskDefinitionId")

	processInstanceKey := r.URL.Query().Get("processInstanceKey")

	searchAfter := r.URL.Query().Get("searchAfter")
	if searchAfter != "" {
		searchAfter = fmt.Sprintf("[%s]", searchAfter)
	}

	searchBefore := r.URL.Query().Get("searchBefore")
	if searchBefore != "" {
		searchBefore = fmt.Sprintf("[%s]", searchBefore)
	}

	return &TaskListQueryParams{
		Size:               size,
		Order:              order,
		Sort:               sort,
		SearchAfter:        searchAfter,
		SearchBefore:       searchBefore,
		State:              state,
		TaskDefinitionId:   taskDefinitionId,
		ProcessInstanceKey: processInstanceKey,
	}, nil
}
