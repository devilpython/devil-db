package model_action

import (
	"fmt"

	"github.com/devilpython/devil-db/db/model"
	"github.com/devilpython/devil-db/db/sql_interface"
	"github.com/devilpython/devil-db/db/sql_utils"
)

// 保存数据
func SaveData(nick string, dataMap map[string]interface{}) (string, bool) {
	message := ""
	successful := true
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		action := sql_utils.Action{}
		message, successful = executeValidater(modelObj, action, sql_interface.ModelPermissionsOperationTypeSave, dataMap)
		executeOperator(modelObj, action, sql_interface.ModelPermissionsOperationTypeSave, dataMap)
		if successful {
			sqlMap := make(map[string]string)
			sql_utils.SetSaveSqlMap(nick, sqlMap, dataMap)
			if !action.SaveModel(sqlMap) {
				message = "Database error"
				successful = false
			}
		}
	} else {
		message = "The specified model does not exist"
		successful = false
	}
	return message, successful
}

// 插入数据
func InsertData(nick string, dataMap map[string]interface{}) (string, bool) {
	message := ""
	successful := true
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		action := sql_utils.Action{}
		message, successful = executeValidater(modelObj, action, sql_interface.ModelPermissionsOperationTypeSave, dataMap)
		executeOperator(modelObj, action, sql_interface.ModelPermissionsOperationTypeSave, dataMap)
		if successful {
			sqlMap := make(map[string]string)
			sql_utils.SetInsertSqlMap(nick, sqlMap, dataMap)
			if !action.SaveModel(sqlMap) {
				message = "Database error"
				successful = false
			}
		}
	} else {
		message = "The specified model does not exist"
		successful = false
	}
	return message, successful
}

// 插入数据
func UpdateData(nick string, dataMap map[string]interface{}) (string, bool) {
	message := ""
	successful := true
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		action := sql_utils.Action{}
		message, successful = executeValidater(modelObj, action, sql_interface.ModelPermissionsOperationTypeSave, dataMap)
		executeOperator(modelObj, action, sql_interface.ModelPermissionsOperationTypeSave, dataMap)
		if successful {
			sqlMap := make(map[string]string)
			sql_utils.SetUpdateSqlMap(nick, sqlMap, dataMap)
			if !action.SaveModel(sqlMap) {
				message = "Database error"
				successful = false
			}
		}
	} else {
		message = "The specified model does not exist"
		successful = false
	}
	return message, successful
}

// 删除数据
func RemoveData(nick string, dataMap map[string]interface{}) (string, bool) {
	message := ""
	successful := true
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		action := sql_utils.Action{}
		message, successful = executeValidater(modelObj, action, sql_interface.ModelPermissionsOperationTypeRemove, dataMap)
		executeOperator(modelObj, action, sql_interface.ModelPermissionsOperationTypeRemove, dataMap)
		if successful && !action.DeleteModel(nick, dataMap) {
			message = "Database error"
			successful = false
		}
	} else {
		message = "The specified model does not exist"
		successful = false
	}
	return message, successful
}

// 查询数据
func QueryData(nick string, dataMap map[string]interface{}) ([]map[string]interface{}, string, bool) {
	message := ""
	successful := true
	modelObj, hasModel := model.GetModel(nick)
	var dataArray []map[string]interface{}
	if hasModel {
		dataArray, message, successful = executeQuery(modelObj, dataMap)
	} else {
		message = "The specified model does not exist"
		successful = false
	}
	return dataArray, message, successful
}

// 通过主键获取数据
func GetData(nick string, dataMap map[string]interface{}) (map[string]interface{}, string, bool) {
	message := ""
	successful := true
	var resultData map[string]interface{}
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		if modelObj.PrimaryKey.IsPrimaryKey { //有主键
			//检查主键
			_, hasPrimaryKey := dataMap[modelObj.PrimaryKey.Name]
			if hasPrimaryKey { //提交的数据有主键
				var dataArray []map[string]interface{}
				dataArray, message, successful = executeQuery(modelObj, dataMap)
				if dataArray != nil && len(dataArray) == 1 { //如果找到唯一数据，设置为返回结果
					resultData = dataArray[0]
				} else {
					message = model.GetQueryParameterErrorMessage()
					successful = false
				}
			} else {
				message = model.GetMissingPrimaryKeyMessage()
				successful = false
			}
		} else { //没有主键
			message = "Primary key does not exist"
		}
	} else {
		message = fmt.Sprintf("The specified model[%s] does not exist", nick)
		successful = false
	}
	return resultData, message, successful
}

// 执行操作器
func executeOperator(modelObj model.Model, action sql_utils.Action, operationType int, dataMap map[string]interface{}) {
	//执行操作器
	for operatorIndex := range modelObj.OperationArray {
		modelObj.OperationArray[operatorIndex].Operate(dataMap, action, modelObj.Nick, modelObj.PrimaryKey.Name, operationType)
	}
}

// 执行验证器
func executeValidater(modelObj model.Model, action sql_utils.Action, operationType int, dataMap map[string]interface{}) (string, bool) {
	//执行验证器
	for validaterIndex := range modelObj.ValidaterArray {
		message, successful := modelObj.ValidaterArray[validaterIndex].Validate(dataMap, action, modelObj.PrimaryKey.Name, operationType)
		if !successful {
			return message, false
		}
	}
	return "", true
}

// 执行查询
func executeQuery(modelObj model.Model, dataMap map[string]interface{}) ([]map[string]interface{}, string, bool) {
	message := ""
	successful := true
	action := sql_utils.Action{}
	var dataArray []map[string]interface{}
	message, successful = executeValidater(modelObj, action, sql_interface.ModelPermissionsOperationTypeQuery, dataMap)
	if successful {
		executeOperator(modelObj, action, sql_interface.ModelPermissionsOperationTypeQueryForParam, dataMap)
		dataArray = action.QueryModel(modelObj.Nick, dataMap)
		if dataArray != nil {
			for index := range dataArray {
				executeOperator(modelObj, action, sql_interface.ModelPermissionsOperationTypeQuery, dataArray[index])
			}
		}
	}
	return dataArray, message, successful
}
