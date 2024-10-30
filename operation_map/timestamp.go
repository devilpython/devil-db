package operation_map

import (
	"fmt"

	"github.com/devilpython/devil-db/db/sql_interface"

	"github.com/devilpython/devil-tools/utils"
)

// 时间戳操作
func timestamp(ai sql_interface.SqlActionInterface, nick, fieldName string, dataMap map[string]interface{}) {
	value, hasData := dataMap[fieldName]
	if hasData {
		strValue := fmt.Sprintf("%v", value)
		dataMap[fieldName] = utils.ConvertTimeToTimestamp(strValue)
	}
}
