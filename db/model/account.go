package model

import (
	"encoding/xml"
	"sync"

	"github.com/devilpython/devil-db/constants"
	"github.com/devilpython/devil-db/db/model_operation"
	"github.com/devilpython/devil-db/db/model_validater"
	"github.com/devilpython/devil-db/db/model_xml"

	"github.com/devilpython/devil-tools/utils"
)

var accountManager *AccountManager
var hasManager bool
var onceManager sync.Once

// 获得模型配置信息
func GetAccountManager() (*AccountManager, bool) {
	onceManager.Do(func() {
		accountManager, hasManager = loadAccountManager(constants.AccountFilePath)
	})
	return accountManager, hasManager
}

// 账号管理器
type AccountManager struct {
	XMLName                     xml.Name     //顶层的账号管理器名称
	IncorrectPermissionsMessage string       //权限不正确错误消息
	AccountTableErrorMessage    string       //账号表配置错误消息
	AccountModel                AccountModel //账号模型
}

// 账号模型
type AccountModel struct {
	Id             string                          //登录账号的唯一标识字段
	Password       string                          //登录账号的密码字段
	ValidaterArray []model_validater.DataValidater //验证器数组
	OperationArray []model_operation.DataOperator  //数据操作器
}

// 加载模型映射表
func loadAccountManager(path string) (*AccountManager, bool) {
	manager := AccountManager{}
	xmlManager := model_xml.AccountManager{}
	hasAccountManager := utils.LoadXmlObject(path, &xmlManager)
	if hasAccountManager {
		utils.CopyData(xmlManager, &manager)
		loadValidaterArrayForAccountManager(&manager.AccountModel, xmlManager.AccountModel.ValidaterListObj)
		loadOperatorArrayForAccountManager(&manager.AccountModel, xmlManager.AccountModel.OperatorListObj)
	}
	return &manager, hasAccountManager
}

// 加载验证器数组
func loadValidaterArrayForAccountManager(accountPointer *AccountModel, validaterList model_xml.ValidaterList) {
	accountPointer.ValidaterArray = loadNilValidaterArray(accountPointer.ValidaterArray, validaterList.NilValidaterArray)
	accountPointer.ValidaterArray = loadLengthValidaterArray(accountPointer.ValidaterArray, validaterList.LengthValidaterArray)
	accountPointer.ValidaterArray = loadRegexValidaterArray(accountPointer.ValidaterArray, validaterList.RegexValidaterArray)
	accountPointer.ValidaterArray = loadExistValidaterArray(accountPointer.ValidaterArray, validaterList.ExistValidaterArray, "account", false)
	accountPointer.ValidaterArray = loadExistValidaterArray(accountPointer.ValidaterArray, validaterList.NotExistValidaterArray, "account", true)
}

// 加载操作器数组
func loadOperatorArrayForAccountManager(accountPointer *AccountModel, operatorList model_xml.OperatorList) {
	accountPointer.OperationArray = loadDataShieldArray(accountPointer.OperationArray, operatorList.DataShieldArray)
	accountPointer.OperationArray = loadDataReviserArray(accountPointer.OperationArray, operatorList.DataReviserArray)
	accountPointer.OperationArray = loadDataPaddingArray(accountPointer.OperationArray, operatorList.DataPaddingArray)
	accountPointer.OperationArray = loadDataExchangerArray(accountPointer.OperationArray, operatorList.DataExchangerArray)
}
