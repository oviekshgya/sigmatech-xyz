package models

import (
	"gorm.io/gorm"
)

func Paginate(pageSize, page int) func(qs *gorm.DB) *gorm.DB {

	return func(qs *gorm.DB) *gorm.DB {

		if pageSize == 0 {
			pageSize = 10
		}

		if page == 0 {
			page = 1
		}
		offset := (page - 1) * pageSize

		return qs.Limit(pageSize).Offset(offset)
	}
}

func Pagination(pageSize, page int) func(qs *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
