package sql_utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/devilpython/devil-db/constants"
	"github.com/devilpython/devil-db/db/model"
)

// 创建插入SQL语句
func CreateInsertSql(modelObj model.Model, dataMap map[string]interface{}) string {
	var strSql bytes.Buffer
	strSql.WriteString("INSERT INTO ")
	strSql.WriteString(modelObj.TableName)
	strSql.WriteString("(")
	var strValues bytes.Buffer
	fieldIndex := 0
	for index := range modelObj.FieldArray {
		fieldName := modelObj.FieldArray[index].Name
		fieldValue, hasField := dataMap[fieldName]
		fieldValue, hasField = createFieldValue(fieldValue, hasField, modelObj.FieldArray[index], dataMap)
		if hasField && modelObj.FieldArray[index].Type != "timestamp" {
			if fieldIndex > 0 {
				strSql.WriteString(",")
				strValues.WriteString(",")
			}
			strSql.WriteString(fieldName)
			strValues.WriteString(toSqlValue(fieldValue))
			fieldIndex += 1
		}
	}
	if fieldIndex > 0 {
		strSql.WriteString(") VALUES (")
		strSql.WriteString(strValues.String())
		strSql.WriteString(")")
		return strSql.String()
	} else {
		return ""
	}
}

// 创建更新SQL语句
func CreateUpdateSql(modelObj model.Model, dataMap map[string]interface{}) string {
	var strSql bytes.Buffer
	strSql.WriteString("UPDATE ")
	strSql.WriteString(modelObj.TableName)
	strSql.WriteString(" SET ")
	var strValues bytes.Buffer
	fieldIndex := 0
	primaryKey := ""
	primaryValue := ""
	accountId := ""
	for index := range modelObj.FieldArray {
		fieldName := modelObj.FieldArray[index].Name
		fieldValue, hasField := dataMap[fieldName]
		fieldValue, hasField = createFieldValue(fieldValue, hasField, modelObj.FieldArray[index], dataMap)
		if hasField && modelObj.FieldArray[index].Type != "timestamp" {
			if modelObj.FieldArray[index].IsPrimaryKey {
				primaryKey = fieldName
				primaryValue = toSqlValue(fieldValue)
			} else {
				if fieldName == "account_id" {
					accountId = toSqlValue(fieldValue)
				} else {
					if fieldIndex > 0 {
						strValues.WriteString(",")
					}
					strValues.WriteString(fieldName)
					strValues.WriteString(" = ")
					strValues.WriteString(toSqlValue(fieldValue))
					fieldIndex += 1
				}
			}
		}
	}
	if fieldIndex > 0 && primaryKey != "" {
		strSql.WriteString(strValues.String())
		strSql.WriteString(" WHERE ")
		strSql.WriteString(primaryKey)
		strSql.WriteString(" = ")
		strSql.WriteString(primaryValue)
		if accountId != "" {
			strSql.WriteString(" AND account_id = ")
			strSql.WriteString(accountId)
		}
		return strSql.String()
	} else {
		return ""
	}
}

// 创建插入SQL语句
func CreateDeleteSql(modelObj model.Model, dataMap map[string]interface{}) string {
	var strSql bytes.Buffer
	strSql.WriteString("DELETE FROM ")
	strSql.WriteString(modelObj.TableName)
	strSql.WriteString(" WHERE ")
	initLength := strSql.Len()
	hasPrimaryKey := false
	if modelObj.PrimaryKey.IsPrimaryKey {
		var fieldValue interface{}
		fieldValue, hasPrimaryKey = dataMap[modelObj.PrimaryKey.Name]
		if hasPrimaryKey {
			strSql.WriteString(modelObj.PrimaryKey.Name)
			strSql.WriteString(" = ")
			strSql.WriteString(toSqlValue(fieldValue))
			for index := range modelObj.FieldArray {
				fieldName := modelObj.FieldArray[index].Name
				fieldValue, hasField := dataMap[fieldName]
				if fieldName == "account_id" && hasField && modelObj.PrimaryKey.Name != "account_id" {
					strSql.WriteString(" and account_id = ")
					strSql.WriteString(toSqlValue(fieldValue))
					break
				}
			}
		}
	}
	if !hasPrimaryKey {
		for index := range modelObj.FieldArray {
			fieldName := modelObj.FieldArray[index].Name
			fieldValue, hasField := dataMap[fieldName]
			if hasField {
				if strSql.Len() > initLength {
					strSql.WriteString(" and ")
				}
				strSql.WriteString(fieldName)
				strSql.WriteString(" = ")
				strSql.WriteString(toSqlValue(fieldValue))
			}
		}
	}
	if strSql.Len() > initLength {
		return strSql.String()
	} else {
		return ""
	}
}

// 创建查询SQL语句
func CreateQuerySql(modelObj model.Model, dataMap map[string]interface{}) string {
	var strSql bytes.Buffer
	strSql.WriteString("SELECT * FROM ")
	strSql.WriteString(modelObj.TableName)
	var strQuery bytes.Buffer
	createQueryString(&strQuery, modelObj, dataMap)
	if strQuery.Len() > 0 {
		strSql.WriteString(" WHERE ")
		strSql.WriteString(strQuery.String())
	}
	return strSql.String()
}

// 创建查询字符串
func createQueryString(strQuery *bytes.Buffer, modelObj model.Model, dataMap map[string]interface{}) {
	paramMap := make(map[string]map[string]string)
	for index := range modelObj.FieldArray {
		value, exist := dataMap[modelObj.FieldArray[index].Name]
		if exist {
			setQueryParamMap(paramMap, modelObj.FieldArray[index], value)
		}
	}
	condition, hasWhereCondition := dataMap[constants.QueryWhereCondition]
	strCondition, conditionOk := condition.(string)
	if hasWhereCondition && conditionOk {
		setQueryParamStringUseCondition(strQuery, paramMap, strCondition)
	} else {
		fuzzy := getFuzzy(dataMap)
		setQueryParamString(strQuery, paramMap, fuzzy)
	}
}

// 获得模糊匹配参数
func getFuzzy(dataMap map[string]interface{}) bool {
	fuzzy := false
	fuzzyData, exist := dataMap["fuzzy"]
	if exist {
		fuzzyData = fmt.Sprintf("%v", fuzzyData)
		fuzzyStr := fuzzyData.(string)
		strings.ToLower(fuzzyStr)
		fuzzyBool, err := strconv.ParseBool(fuzzyStr)
		if err == nil {
			fuzzy = fuzzyBool
		}
	}
	return fuzzy
}

// 设置查询参数字符串
func setQueryParamStringUseCondition(strQuery *bytes.Buffer, paramMap map[string]map[string]string, condition string) {
	for fieldName, value := range paramMap {
		var queryVar bytes.Buffer
		queryVar.WriteString("{")
		queryVar.WriteString(fieldName)
		queryVar.WriteString("}")
		condition = strings.ReplaceAll(condition, queryVar.String(), value["base"]) //替换基础参数
		queryVar.Reset()
		queryVar.WriteString("{")
		queryVar.WriteString(fieldName)
		queryVar.WriteString("_fuzzy")
		queryVar.WriteString("}")
		condition = strings.ReplaceAll(condition, queryVar.String(), value["like"]) //替换模糊参数
	}
	strQuery.WriteString(condition)
}

// 设置查询参数字符串
func setQueryParamString(strQuery *bytes.Buffer, paramMap map[string]map[string]string, fuzzy bool) {
	for _, value := range paramMap {
		if strQuery.Len() > 0 {
			strQuery.WriteString(" and ")
		}
		if fuzzy {
			strQuery.WriteString(value["like"])
		} else {
			strQuery.WriteString(value["base"])
		}
	}
}

// 设置查询参数映射表
func setQueryParamMap(paramMap map[string]map[string]string, field model.Field, value interface{}) {
	var likeQuery bytes.Buffer
	var baseQuery bytes.Buffer
	value = fmt.Sprintf("%v", value)
	valueStr := value.(string)
	valueStr = strings.TrimSpace(valueStr)
	if len(valueStr) > 0 {
		if valueStr == "<nil>" {
			baseQuery.WriteString(field.Name)
			baseQuery.WriteString(" = NULL")
		} else if field.Type == "string" {
			likeQuery.WriteString(field.Name)
			likeQuery.WriteString(" like '%")
			likeQuery.WriteString(valueStr)
			likeQuery.WriteString("%'")

			baseQuery.WriteString(field.Name)
			baseQuery.WriteString(" = '")
			baseQuery.WriteString(valueStr)
			baseQuery.WriteString("'")
		} else {
			baseQuery.WriteString(field.Name)
			baseQuery.WriteString(" = ")
			baseQuery.WriteString(valueStr)
		}
		paramValue := make(map[string]string)
		paramValue["base"] = baseQuery.String()
		if likeQuery.Len() > 0 {
			paramValue["like"] = likeQuery.String()
		} else {
			paramValue["like"] = baseQuery.String()
		}
		paramMap[field.Name] = paramValue
	}
}

// 转换到SQL的数据值
func toSqlValue(fieldValue interface{}) string {
	strValue, isString := fieldValue.(string)
	if isString {
		var buffer bytes.Buffer
		buffer.WriteString("'")
		buffer.WriteString(strValue)
		buffer.WriteString("'")
		return buffer.String()
	} else {
		return fmt.Sprintf("%v", fieldValue)
	}
}
