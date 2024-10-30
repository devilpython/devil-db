package model

import (
	"fmt"
	"strings"

	"github.com/devilpython/devil-db/db/model_operation"
	"github.com/devilpython/devil-db/db/model_validater"
	"github.com/devilpython/devil-db/db/model_xml"
	"github.com/devilpython/devil-db/db/sql_interface"

	"github.com/devilpython/devil-tools/utils"
)

// 加载模型映射表
func LoadModelMap(path string) *ModelConfig {
	config := ModelConfig{}
	config.ModelMap = make(map[string]*Model)
	xmlModelList := model_xml.ModelList{}
	hasModel := utils.LoadXmlObject(path, &xmlModelList)
	if hasModel {
		utils.CopyData(xmlModelList, &config)
		for index := range xmlModelList.ModelArray {
			mModel := Model{}
			loadModel(&mModel, xmlModelList.ModelArray[index])
			config.ModelMap[mModel.Nick] = &mModel
			if mModel.Nick == "account" {
				config.AccountModel = mModel
			} else if mModel.Nick == "admin" {
				config.AdminModel = mModel
			} else if mModel.Nick == "token" {
				config.TokenModel = mModel
			}
		}
		setChildren(&config.ModelMap)
	}
	return &config
}

// 设置子模型
func setChildren(modelMap *map[string]*Model) {
	for _, modelObj := range *modelMap { //为父模型设置子模型
		for index := range modelObj.FieldArray {
			if modelObj.FieldArray[index].TargetModel != "" {
				parentModel, hasParent := (*modelMap)[modelObj.FieldArray[index].TargetModel]
				if hasParent {
					if parentModel.ChildrenField == nil {
						parentModel.ChildrenField = make(map[string]Field)
					}
					parentModel.ChildrenField[modelObj.Nick] = modelObj.FieldArray[index]
				}
			}
		}
	}
}

// 加载模型
func loadModel(modelPointer *Model, xmlModel model_xml.Model) {
	modelPointer.Nick = xmlModel.Nick
	modelPointer.TableName = xmlModel.Nick
	if xmlModel.TableName != "" {
		modelPointer.TableName = xmlModel.TableName
	}
	strType := strings.ToLower(xmlModel.ReadPermissions)
	if strType == "all" {
		modelPointer.ReadPermissions = sql_interface.ModelPermissionsAll
	} else if strType == "user" {
		modelPointer.ReadPermissions = sql_interface.ModelPermissionsUser
	} else if strType == "admin" {
		modelPointer.ReadPermissions = sql_interface.ModelPermissionsAdmin
	} else if strType == "sys" {
		modelPointer.ReadPermissions = sql_interface.ModelPermissionsSys
	}
	loadModelFieldArray(modelPointer, xmlModel)
	loadValidaterArray(modelPointer, xmlModel.ValidaterListObj)
	loadOperatorArray(modelPointer, xmlModel.OperatorListObj)
}

// 加载模型字段数组
func loadModelFieldArray(modelPointer *Model, xmlModel model_xml.Model) {
	for index := range xmlModel.FieldArray {
		field := Field{}
		xmlField := xmlModel.FieldArray[index]
		loadModelField(&field, xmlField)
		if field.IsPrimaryKey {
			modelPointer.PrimaryKey = field
		}
		modelPointer.FieldArray = append(modelPointer.FieldArray, field)
	}
}

// 加载模型字段
func loadModelField(fieldPointer *Field, xmlField model_xml.Field) {
	fieldPointer.Name = xmlField.Name
	fieldPointer.Type = strings.ToLower(xmlField.Type)
	fieldPointer.Create = xmlField.Create
	fieldPointer.IsPrimaryKey = xmlField.IsPrimaryKey
	fieldPointer.TargetModel = xmlField.TargetModel
	fieldPointer.TargetField = xmlField.TargetField
	if xmlField.IsUserId {
		fieldPointer.Flag = 1
	} else if xmlField.IsUserPassword {
		fieldPointer.Flag = 2
	}
}

//==========================================加载验证器数组=====开始==========================================

// 加载验证器数组
func loadValidaterArray(modelPointer *Model, validaterList model_xml.ValidaterList) {
	modelPointer.ValidaterArray = loadNilValidaterArray(modelPointer.ValidaterArray, validaterList.NilValidaterArray)
	modelPointer.ValidaterArray = loadLengthValidaterArray(modelPointer.ValidaterArray, validaterList.LengthValidaterArray)
	modelPointer.ValidaterArray = loadRegexValidaterArray(modelPointer.ValidaterArray, validaterList.RegexValidaterArray)
	modelPointer.ValidaterArray = loadExistValidaterArray(modelPointer.ValidaterArray, validaterList.ExistValidaterArray, modelPointer.Nick, false)
	modelPointer.ValidaterArray = loadExistValidaterArray(modelPointer.ValidaterArray, validaterList.NotExistValidaterArray, modelPointer.Nick, true)
}

// 加载空验证器数组
func loadNilValidaterArray(validaterArray []model_validater.DataValidater, xmlValidaterArray []model_xml.NilValidater) []model_validater.DataValidater {
	for index := range xmlValidaterArray {
		nilValidater := model_validater.NilValidater{}
		utils.CopyData(xmlValidaterArray[index], &nilValidater)
		nilValidater.OperateType = getActionOperationType(xmlValidaterArray[index].ForAction)
		validater := model_validater.DataValidater{}
		validater.Validater = nilValidater
		validaterArray = append(validaterArray, validater)
	}
	return validaterArray
}

// 加载长度验证器数组
func loadLengthValidaterArray(validaterArray []model_validater.DataValidater, xmlLengthValidaterArray []model_xml.LengthValidater) []model_validater.DataValidater {
	for index := range xmlLengthValidaterArray {
		lengthValidater := model_validater.LengthValidater{}
		utils.CopyData(xmlLengthValidaterArray[index], &lengthValidater)
		lengthValidater.OperateType = getActionOperationType(xmlLengthValidaterArray[index].ForAction)
		validater := model_validater.DataValidater{}
		validater.Validater = lengthValidater
		validaterArray = append(validaterArray, validater)
	}
	return validaterArray
}

// 加载正则验证器数组
func loadRegexValidaterArray(validaterArray []model_validater.DataValidater, xmlRegexValidaterArray []model_xml.RegexValidater) []model_validater.DataValidater {
	for index := range xmlRegexValidaterArray {
		regexValidater := model_validater.RegexValidater{}
		utils.CopyData(xmlRegexValidaterArray[index], &regexValidater)
		regexValidater.OperateType = getActionOperationType(xmlRegexValidaterArray[index].ForAction)
		regexValidater.Regex = xmlRegexValidaterArray[index].Regex
		validater := model_validater.DataValidater{}
		validater.Validater = regexValidater
		validaterArray = append(validaterArray, validater)
	}
	return validaterArray
}

// 加载存在验证器数组
func loadExistValidaterArray(validaterArray []model_validater.DataValidater, xmlExistValidaterArray []model_xml.IsNotExistValidater, modelNick string, isNot bool) []model_validater.DataValidater {
	//fmt.Println("..................loadExistValidaterArray:", isNot)
	for index := range xmlExistValidaterArray {
		existValidater := model_validater.ExistValidater{}
		existValidater.IsNot = isNot
		utils.CopyData(xmlExistValidaterArray[index], &existValidater)
		existValidater.OperateType = getActionOperationType(xmlExistValidaterArray[index].ForAction)
		existValidater.TargetModel = modelNick
		if xmlExistValidaterArray[index].TargetModel != "" {
			existValidater.TargetModel = xmlExistValidaterArray[index].TargetModel
		}
		validater := model_validater.DataValidater{}
		validater.Validater = existValidater
		validaterArray = append(validaterArray, validater)
	}
	return validaterArray
}

//==========================================加载验证器数组=====结束==========================================

//==========================================加载操作器数组=====开始==========================================

// 加载操作器数组
func loadOperatorArray(modelPointer *Model, operatorList model_xml.OperatorList) {
	modelPointer.OperationArray = loadDataShieldArray(modelPointer.OperationArray, operatorList.DataShieldArray)
	modelPointer.OperationArray = loadDataReviserArray(modelPointer.OperationArray, operatorList.DataReviserArray)
	modelPointer.OperationArray = loadDataPaddingArray(modelPointer.OperationArray, operatorList.DataPaddingArray)
	modelPointer.OperationArray = loadDataExchangerArray(modelPointer.OperationArray, operatorList.DataExchangerArray)
}

// 加载数据屏蔽器数组
func loadDataShieldArray(operatorArray []model_operation.DataOperator, xmlDataShieldArray []model_xml.DataShield) []model_operation.DataOperator {
	for index := range xmlDataShieldArray {
		dataShield := model_operation.DataShield{}
		dataShield.Name = xmlDataShieldArray[index].Name
		dataShield.OperateType = getActionOperationType(xmlDataShieldArray[index].ForAction)
		operator := model_operation.DataOperator{}
		operator.Operator = dataShield
		operatorArray = append(operatorArray, operator)
	}
	return operatorArray
}

// 加载数据修改器数组
func loadDataReviserArray(operatorArray []model_operation.DataOperator, xmlDataReviserArray []model_xml.DataReviser) []model_operation.DataOperator {
	for index := range xmlDataReviserArray {
		dataReviser := model_operation.DataReviser{}
		utils.CopyData(xmlDataReviserArray[index], &dataReviser)
		dataReviser.OperateType = getActionOperationType(xmlDataReviserArray[index].ForAction)
		operator := model_operation.DataOperator{}
		operator.Operator = dataReviser
		operatorArray = append(operatorArray, operator)
	}
	return operatorArray
}

// 加载数据填充器数组
func loadDataPaddingArray(operatorArray []model_operation.DataOperator, xmlDataPaddingArray []model_xml.DataPadding) []model_operation.DataOperator {
	for index := range xmlDataPaddingArray {
		dataPadding := model_operation.DataPadding{}
		utils.CopyData(xmlDataPaddingArray[index], &dataPadding)
		dataPadding.OperateType = getActionOperationType(xmlDataPaddingArray[index].ForAction)
		operator := model_operation.DataOperator{}
		operator.Operator = dataPadding
		operatorArray = append(operatorArray, operator)
	}
	return operatorArray
}

// 加载数据转换器数组
func loadDataExchangerArray(operatorArray []model_operation.DataOperator, xmlDataExchangerArray []model_xml.DataExchanger) []model_operation.DataOperator {
	for index := range xmlDataExchangerArray {
		dataExchanger := model_operation.DataExchanger{}
		utils.CopyData(xmlDataExchangerArray[index], &dataExchanger)
		strData := strings.TrimSpace(xmlDataExchangerArray[index].ExchangeData)
		dataMap := make(map[string]interface{})
		err := utils.ConvertJsonToData(strData, &dataMap)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			dataExchanger.ExchangeData = dataMap
			dataExchanger.OperateType = getActionOperationType(xmlDataExchangerArray[index].ForAction)
			operator := model_operation.DataOperator{}
			operator.Operator = dataExchanger
			operatorArray = append(operatorArray, operator)
		}
	}
	return operatorArray
}

//==========================================加载操作器数组=====结束==========================================

// 获得动作的操作类型
func getActionOperationType(actionInfo string) int {
	result := 0
	actionInfo = strings.ToLower(actionInfo)
	if strings.Contains(actionInfo, "save") {
		result |= sql_interface.ModelPermissionsOperationTypeSave
	}
	if strings.Contains(actionInfo, "query") {
		result |= sql_interface.ModelPermissionsOperationTypeQuery
	}
	if strings.Contains(actionInfo, "query-param") {
		result |= sql_interface.ModelPermissionsOperationTypeQueryForParam
	}
	if strings.Contains(actionInfo, "remove") {
		result |= sql_interface.ModelPermissionsOperationTypeRemove
	}
	return result
}
