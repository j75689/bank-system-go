package model

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Pagination struct {
	Page    int64 `json:"page,omitempty"`
	PerPage int64 `json:"per_page,omitempty"`
}

func (p Pagination) LimitAndOffset(db *gorm.DB) *gorm.DB {
	if p.PerPage != 0 || p.Offset() != 0 {
		db = db.Limit(int(p.PerPage)).Offset(int(p.Offset()))
	}
	return db
}

func (p Pagination) Offset() int64 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}

type SortOrder string

func (field SortOrder) String() string {
	return string(field)
}

const (
	SortASC  SortOrder = "ASC"
	SortDESC SortOrder = "DESC"
)

type SortField struct {
	Field string    `json:"sort_field,omitempty"`
	Order SortOrder `json:"sort_order,omitempty"`
}

type Sorting []SortField

func (s Sorting) Sort(db *gorm.DB) *gorm.DB {
	sortfield := []string{}
	for _, sort := range s {
		if len(sort.Field) != 0 && len(sort.Order) != 0 {
			sortfield = append(sortfield, fmt.Sprintf("%s %s", sort.Field, sort.Order))
		}
	}

	if len(sortfield) > 0 {
		db = db.Order(strings.Join(sortfield, ","))
	}

	return db
}

type BaseFilter struct {
	CreateAtGte time.Time `json:"create_at_gte,omitempty"`
	CreateAtLte time.Time `json:"create_at_lte,omitempty"`
}

var _emptyTime = time.Time{}

func (f BaseFilter) Filter(db *gorm.DB) *gorm.DB {
	if f.CreateAtGte != _emptyTime {
		db.Where("created_at >=?", f.CreateAtGte)
	}
	if f.CreateAtLte != _emptyTime {
		db.Where("created_at <=?", f.CreateAtLte)
	}
	return db
}
