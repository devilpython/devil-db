package operation_map

import (
	"fmt"

	"github.com/devilpython/devil-db/db/sql_interface"

	"github.com/devilpython/devil-tools/utils"
)

// MD5操作
func md5(ai sql_interface.SqlActionInterface, nick, fieldName string, dataMap map[string]interface{}) {
	value, hasData := dataMap[fieldName]
	if hasData {
		strValue := fmt.Sprintf("%v", value)
		dataMap[fieldName] = utils.Md5(strValue)
	}
}
