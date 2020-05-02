package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/my1562/api/utils"
)

func ApplyListQueryParams(db *gorm.DB, model interface{}, params *utils.ListQueryParams) *gorm.DB {

	if params.HasRange {
		db = db.Offset(params.Range.From).Limit(params.Range.To - params.Range.From + 1)
	}
	if params.HasSort {
		scope := db.NewScope(model)
		if field, ok := scope.FieldByName(params.Sort.Field); ok {
			db = db.Order(field.DBName + " " + string(params.Sort.Direction))
		}
	}
	return db
}

func GetResultRange(db *gorm.DB, model interface{}, params *utils.ListQueryParams, resourceName string) (string, error) {
	var total int
	if err := db.Model(model).Count(&total).Error; err != nil {
		return "", err
	}

	if params.HasRange {
		contentRange := fmt.Sprintf(
			"%s %d-%d/%d",
			resourceName,
			params.Range.From,
			params.Range.To,
			total,
		)
		return contentRange, nil
	}
	return "", nil

}
