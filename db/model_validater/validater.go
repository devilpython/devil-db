package model_validater

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/devilpython/devil-db/db/sql_interface"
)

// 数据验证器对象，里面封装了具体的验证器
type DataValidater struct {
	Validater interface{}
}

// 空验证器
type NilValidater struct {
	FieldName   string //字段名
	Message     string //错误消息
	OperateType int    //支持的数据库操作类型
}

// 长度验证器
type LengthValidater struct {
	FieldName   string //字段名
	MinLength   int    //最小长度
	MaxLength   int    //最大长度
	Message     string //错误消息
	OperateType int    //支持的数据库操作类型
}

// 正则表达式验证器
type RegexValidater struct {
	FieldName   string //字段名
	Regex       string //正则表达式
	Message     string //错误消息
	OperateType int    //支持的数据库操作类型
}

// 验证字段
type ValidateField struct {
	Name        string //字段名
	TargetField string //目标字段名
}

// 存在验证器
type ExistValidater struct {
	TargetModel    string          //目标模型
	ConditionField string          //条件字段
	ConditionValue string          //条件值
	FieldArray     []ValidateField //验证的字段数组
	Message        string          //错误消息
	OperateType    int             //支持的数据库操作类型
	IsNot          bool            //是否为不存在验证
}

// 验证器接口
type ValidaterInterface interface {
	//验证方法
	Validate(dataMap map[string]interface{}, ai sql_interface.SqlActionInterface, primaryKey string, operationType int) (string, bool)
}

// 验证器的验证方法
func (validater DataValidater) Validate(dataMap map[string]interface{}, ai sql_interface.SqlActionInterface, primaryKey string, operationType int) (string, bool) {
	switch validaterObj := validater.Validater.(type) {
	case NilValidater:
		if validaterObj.OperateType&operationType > 0 {
			return validateNil(dataMap, validaterObj, primaryKey)
		}
	case LengthValidater:
		if validaterObj.OperateType&operationType > 0 {
			return validateLength(dataMap, validaterObj)
		}
	case RegexValidater:
		if validaterObj.OperateType&operationType > 0 {
			return validateRegex(dataMap, validaterObj)
		}
	case ExistValidater:
		if validaterObj.OperateType&operationType > 0 {
			if validaterObj.IsNot {
				return validateNotExist(ai, dataMap, validaterObj)
			} else {
				return validateExist(ai, dataMap, validaterObj)
			}
		}
	}
	return "", true
}

// 验证是否空值
func validateNil(dataMap map[string]interface{}, validater NilValidater, primaryKey string) (string, bool) {
	_, hasPrimaryKey := dataMap[primaryKey]
	value, hasValue := dataMap[validater.FieldName]
	strValue := ""
	ok := false
	if hasValue {
		strValue, ok = value.(string)
	}
	if hasPrimaryKey && hasValue && ok && strings.TrimSpace(strValue) == "" {
		return validater.Message, false
	} else if !hasPrimaryKey && (!hasValue || strings.TrimSpace(strValue) == "") {
		return validater.Message, false
	} else {
		return "", true
	}
}

// 验证数据长度
func validateLength(dataMap map[string]interface{}, validater LengthValidater) (string, bool) {
	value, hasValue := dataMap[validater.FieldName]
	strValue := ""
	var arrayValue []interface{}
	isString := false
	isArray := false
	if hasValue {
		strValue, isString = value.(string)
		if !isString {
			arrayValue, isArray = value.([]interface{})
		}
	}
	if validater.MaxLength > 0 && ((isString && utf8.RuneCountInString(strValue) >= validater.MaxLength) || (isArray && len(arrayValue) >= validater.MaxLength)) {
		return getMessage(validater.Message, "{max-length}", validater.MaxLength), false
	} else if validater.MinLength > 0 && ((isString && utf8.RuneCountInString(strValue) <= validater.MinLength) || (isArray && len(arrayValue) <= validater.MinLength)) {
		return getMessage(validater.Message, "{min-length}", validater.MaxLength), false
	} else {
		return "", true
	}
}

// 验证正则表达式
func validateRegex(dataMap map[string]interface{}, validater RegexValidater) (string, bool) {
	value, hasValue := dataMap[validater.FieldName]
	strValue := ""
	ok := false
	if hasValue {
		strValue, ok = value.(string)
	}
	if ok {
		regexValidate, _ := regexp.MatchString(validater.Regex, strValue)
		if !regexValidate {
			return validater.Message, false
		}
	}
	return "", true
}

// 验证数据是否存在
func validateExist(ai sql_interface.SqlActionInterface, dataMap map[string]interface{}, validater ExistValidater) (string, bool) {
	mustValidate := checkValidateCondition(dataMap, validater)
	if mustValidate {
		targetDataMap := createdTargetDataMap(dataMap, validater)
		if len(targetDataMap) > 0 && isExist(ai, targetDataMap, validater) {
			return validater.Message, false
		}
	}
	return "", true
}

// 验证数据是否不存在
func validateNotExist(ai sql_interface.SqlActionInterface, dataMap map[string]interface{}, validater ExistValidater) (string, bool) {
	mustValidate := checkValidateCondition(dataMap, validater)
	if mustValidate {
		targetDataMap := createdTargetDataMap(dataMap, validater)
		if len(targetDataMap) > 0 && !isExist(ai, targetDataMap, validater) {
			return validater.Message, false
		}
	}
	return "", true
}

// 获得消息
func getMessage(message string, name string, value int) string {
	if len(name) > 0 {
		return strings.ReplaceAll(message, name, strconv.Itoa(value))
	} else {
		return message
	}
}

// 检查验证条件
func checkValidateCondition(dataMap map[string]interface{}, validater ExistValidater) bool {
	mustValidate := true
	if validater.ConditionField != "" {
		value, hasValue := dataMap[validater.ConditionField]
		if hasValue {
			strValue := fmt.Sprintf("%v", value)
			if strValue != validater.ConditionValue {
				mustValidate = false
			}
		}
	}
	return mustValidate
}

// 创建目标数据映射表
func createdTargetDataMap(dataMap map[string]interface{}, validater ExistValidater) map[string]interface{} {
	targetDataMap := make(map[string]interface{})
	for index := range validater.FieldArray {
		value, hasValue := dataMap[validater.FieldArray[index].Name]
		if hasValue {
			if len(validater.FieldArray[index].TargetField) > 0 {
				targetDataMap[validater.FieldArray[index].TargetField] = value
			} else {
				targetDataMap[validater.FieldArray[index].Name] = value
			}
		}
	}
	return targetDataMap
}

// 验证是否存在
func isExist(ai sql_interface.SqlActionInterface, targetDataMap map[string]interface{}, validater ExistValidater) bool {
	if len(targetDataMap) > 0 {
		resultMap := ai.QueryModel(validater.TargetModel, targetDataMap)
		if resultMap != nil && len(resultMap) > 0 {
			return true
		}
	}
	return false
}
