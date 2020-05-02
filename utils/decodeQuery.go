package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type SortDirection string

const (
	SortDirectionAscending  SortDirection = "ASC"
	SortDirectionDescending               = "DESC"
)

type ListQueryParams struct {
	HasRange bool
	HasSort  bool
	Range    struct {
		From int64
		To   int64
	}
	Sort struct {
		Field     string
		Direction SortDirection
	}
}

type ReactAdminListRequest struct {
	Range string `form:"range"`
	Sort  string `form:"sort"`
}

func DecodeReactAdminQueryParams(c *gin.Context) (*ListQueryParams, error) {
	req := ReactAdminListRequest{}

	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, err
	}

	params := &ListQueryParams{}

	if req.Range != "" {
		rawRange := []int64{}
		err := json.Unmarshal([]byte(req.Range), &rawRange)
		if err == nil {
			params.HasRange = true
			params.Range.From = rawRange[0]
			params.Range.To = rawRange[1]
		}
	}
	if req.Sort != "" {
		rawSort := []string{}
		err := json.Unmarshal([]byte(req.Sort), &rawSort)
		if err == nil {
			params.HasSort = true
			params.Sort.Field = rawSort[0]
			if rawSort[1] == "DESC" {
				params.Sort.Direction = SortDirectionDescending
			} else {
				params.Sort.Direction = SortDirectionAscending
			}
		}
	}

	return params, nil
}
