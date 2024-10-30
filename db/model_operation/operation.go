package model_operation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/devilpython/devil-db/db/sql_interface"
	"github.com/devilpython/devil-db/operation_map"

	"github.com/devilpython/devil-tools/cache"
	"github.com/devilpython/devil-tools/utils"
)

// 数据操作器
type DataOperator struct {
	Operator interface{}
}

// 数据屏蔽器
type DataShield struct {
	Name        string //字段名
	OperateType int    //支持的数据库操作类型
}

// 数据修改器
type DataReviser struct {
	Name        string //字段名
	Method      string //修改器方法
	OperateType int    //支持的数据库操作类型
}

// 数据填充器
type DataPadding struct {
	Name        string //字段名
	Method      string //填充器方法
	Param       string //填充器参数
	OperateType int    //支持的数据库操作类型
}

// 数据转换器
type DataExchanger struct {
	Name           string                 //字段名
	ConditionField string                 //条件字段
	ConditionValue string                 //条件值
	ExchangeData   map[string]interface{} //需要转换的数据
	OperateType    int                    //支持的数据库操作类型
}

// 验证器接口
type DataOperatorInterface interface {
	//operate(dataMap map[string]interface{}, primaryKey string, operationType int) (string, bool)
	Operate(dataMap map[string]interface{}, ai sql_interface.SqlActionInterface, nick, primaryKey string, operationType int)
}

// 数据操作器的操作方法
func (operator DataOperator) Operate(dataMap map[string]interface{}, ai sql_interface.SqlActionInterface, nick, primaryKey string, operationType int) {
	switch operatorObj := operator.Operator.(type) {
	case DataShield:
		if operatorObj.OperateType&operationType > 0 {
			maskData(dataMap, operatorObj)
		}
	case DataReviser:
		if operatorObj.OperateType&operationType > 0 {
			reviseData(ai, nick, dataMap, operatorObj)
		}
	case DataPadding:
		if operatorObj.OperateType&operationType > 0 {
			paddingData(dataMap, operatorObj)
		}
	case DataExchanger:
		if operationType&operatorObj.OperateType == sql_interface.ModelPermissionsOperationTypeSave {
			exchangeDataToDB(dataMap, operatorObj)
		} else if operationType&operatorObj.OperateType == sql_interface.ModelPermissionsOperationTypeQuery {
			exchangeDataFromDB(primaryKey, dataMap, operatorObj)
		}
	}
}

// 屏蔽数据
func maskData(dataMap map[string]interface{}, operator DataShield) {
	_, hasData := dataMap[operator.Name]
	if hasData {
		delete(dataMap, operator.Name)
	}
}

// 修改数据
func reviseData(ai sql_interface.SqlActionInterface, nick string, dataMap map[string]interface{}, operator DataReviser) {
	funcObj, hasFunc := operation_map.GetOperationFunc(operator.Method)
	if hasFunc {
		funcObj(ai, nick, operator.Name, dataMap)
	}
}

// 填充数据
func paddingData(dataMap map[string]interface{}, operator DataPadding) {
	if operator.Method == "redis" && len(operator.Param) > 0 {
		redisKey := operator.Param
		for key := range dataMap {
			redisKey = strings.ReplaceAll(redisKey, fmt.Sprintf("{%s}", key), fmt.Sprintf("%v", dataMap[key]))
		}
		data, hasData := cache.Get(redisKey)
		if hasData {
			dataMap[operator.Name] = data
		}
	}
}

// 转换数据到数据库
func exchangeDataToDB(dataMap map[string]interface{}, operator DataExchanger) {
	changeMap(operator.ExchangeData, dataMap)
	jsonData, err := utils.ConvertDataToJson(dataMap)
	if err == nil {
		dataMap[operator.Name] = jsonData
	}
}

// 从数据库取出数据并转换
func exchangeDataFromDB(primaryKey string, dataMap map[string]interface{}, operator DataExchanger) map[string]interface{} {
	jsonData, hasJsonData := dataMap[operator.Name]
	if hasJsonData {
		strJsonData, jsonOk := jsonData.(string)
		if jsonOk {
			resultMap := make(map[string]interface{})
			err := utils.ConvertJsonToData(strJsonData, &resultMap)
			if err == nil {
				if len(primaryKey) > 0 {
					pkValue, hasPrimaryKey := dataMap[primaryKey]
					if hasPrimaryKey {
						resultMap[primaryKey] = pkValue
					}
				}
				return resultMap
			}
		}
	}
	return operator.ExchangeData
}

// 转换数据
func changeMap(srcMap map[string]interface{}, destMap map[string]interface{}) {
	for key, value := range srcMap {
		data, hasData := destMap[key]
		if hasData {
			//类型相等，检查是否数组，是数组则复制数组
			//类型不等，赋值默认数据
			if reflect.TypeOf(data).Kind() == reflect.TypeOf(value).Kind() {
				if reflect.TypeOf(data).Kind() == reflect.Slice {
					srcArray, srcOk := value.([]interface{})
					destArray, destOk := data.([]interface{})
					if srcOk && destOk {
						if len(destArray) > 0 && len(srcArray) > 0 {
							changeArray(srcArray, destArray)
						} else {
							var nilArray []interface{}
							destMap[key] = nilArray
						}
					}
				}
			} else {
				setDefaultValue(destMap, key, value)
			}
		} else {
			setDefaultValue(destMap, key, value)
		}
	}
}

// 转换数组
func changeArray(srcArray []interface{}, destArray []interface{}) {
	for index := range destArray {
		srcMap, srcOk := srcArray[0].(map[string]interface{})
		destMap, destOk := destArray[index].(map[string]interface{})
		if srcOk && destOk {
			changeMap(srcMap, destMap)
		}
	}
}

// 设置默认值
func setDefaultValue(destMap map[string]interface{}, key string, value interface{}) {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		var nilArray []interface{}
		destMap[key] = nilArray
	} else {
		destMap[key] = value
	}
}
