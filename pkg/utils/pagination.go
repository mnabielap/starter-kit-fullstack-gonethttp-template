package utils

import (
	"math"

	"gorm.io/gorm"
)

// PaginationScope holds the parameters for pagination
type PaginationScope struct {
	Page  int
	Limit int
	Sort  string
}

// PaginationResult matches the JSON structure of the PHP project
type PaginationResult struct {
	Results      interface{} `json:"results"`
	Page         int         `json:"page"`
	Limit        interface{} `json:"limit"` // Interface because it can be "all" or int
	TotalPages   int         `json:"totalPages"`
	TotalResults int64       `json:"totalResults"`
}

// Paginate returns a GORM scope to handle offset and limit
func (p *PaginationScope) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Limit <= 0 {
			// If limit is 0 or negative (effectively "all" in logic handling), don't limit query
			return db
		}
		offset := (p.Page - 1) * p.Limit
		return db.Offset(offset).Limit(p.Limit)
	}
}

// GetPaginationResult constructs the response object
func GetPaginationResult(totalRows int64, page int, limit int, results interface{}) PaginationResult {
	var totalPages int
	var limitVal interface{} = limit

	if limit > 0 {
		totalPages = int(math.Ceil(float64(totalRows) / float64(limit)))
	} else {
		totalPages = 1
		limitVal = "all"
	}

	return PaginationResult{
		Results:      results,
		Page:         page,
		Limit:        limitVal,
		TotalPages:   totalPages,
		TotalResults: totalRows,
	}
}