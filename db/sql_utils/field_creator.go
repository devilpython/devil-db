package sql_utils

import (
	"github.com/devilpython/devil-db/db/model"

	"github.com/devilpython/devil-tools/utils"
)

// 创建字段值
func createFieldValue(fieldValue interface{}, hasValue bool, field model.Field, dataMap map[string]interface{}) (interface{}, bool) {
	if !hasValue && fieldValue == nil && field.Create == "md5" {
		idString := utils.CreateId()
		dataMap[field.Name] = idString
		return idString, true
	}
	return fieldValue, hasValue
}
