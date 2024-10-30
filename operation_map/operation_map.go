package operation_map

import (
	"sync"

	"github.com/devilpython/devil-db/db/sql_interface"
)

var funcMap map[string]func(sql_interface.SqlActionInterface, string, string, map[string]interface{})
var once sync.Once

// 初始化操作器函数
func InitOperationFunc(userFuncMap map[string]func(sql_interface.SqlActionInterface, string, string, map[string]interface{})) {
	once.Do(func() {
		funcMap = make(map[string]func(sql_interface.SqlActionInterface, string, string, map[string]interface{}))
		funcMap["md5"] = md5
		funcMap["timestamp"] = timestamp
		//funcMap["RemoveInvalidCorpus"] = flow_utils.RemoveInvalidCorpus
	})
	if userFuncMap != nil {
		for key, value := range userFuncMap {
			funcMap[key] = value
		}
	}
}

// 获得操作器函数
func GetOperationFunc(funcName string) (func(sql_interface.SqlActionInterface, string, string, map[string]interface{}), bool) {
	InitOperationFunc(nil)
	funcObj, hasFunc := funcMap[funcName]
	return funcObj, hasFunc
}
