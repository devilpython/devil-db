package model_account

import (
	"fmt"

	"github.com/devilpython/devil-db/db/model"
	"github.com/devilpython/devil-db/db/model_action"
)

// 通过Token获得账号ID
func GetAccountIdForToken(token string) (string, bool, string) {
	accountId := ""
	isAdmin := false
	message := ""
	tokenModel, hasTokenModel := model.GetModel("token")
	if hasTokenModel { //检查是否配置了token表
		tokenMap := make(map[string]interface{})
		tokenMap["token"] = token
		dataArray, _message, successful := model_action.QueryData(tokenModel.Nick, tokenMap)
		if successful && dataArray != nil && len(dataArray) > 0 { //存在token
			accountId = fmt.Sprintf("%v", dataArray[0][tokenModel.PrimaryKey.Name])
			adminModel, hasAdminModel := model.GetModel("admin")
			if hasAdminModel { //检查是否存在管理员表
				adminMap := make(map[string]interface{})
				var adminData map[string]interface{}
				adminData, message, successful = model_action.GetData(adminModel.Nick, adminMap)
				if successful && adminData != nil {
					isAdmin = true
				}
			}
		} else {
			message = _message
		}
	}
	return accountId, isAdmin, message
}
