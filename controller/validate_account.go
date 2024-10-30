package controller

import (
	"github.com/devilpython/devil-db/db/model"
	"github.com/devilpython/devil-db/db/sql_interface"
	"github.com/devilpython/devil-db/global_keys"

	"github.com/devilpython/devil-tools/utils"

	"github.com/devilpython/devil-tools/middle_ware"
	"github.com/gin-gonic/gin"
)

// 验证账号
func validateAccount(context *gin.Context, operationType int) (string, map[string]interface{}, string, sql_interface.ModelPermissions, bool, string) {
	dataMap, message, has := getDataMapParam()
	if has {
		nick := context.Param("nick")
		accountId, _ := utils.GetGlobalData(global_keys.KeyAccountId)
		accountIdStr, _ := accountId.(string)
		dataLevel := model.GetLevel(nick, operationType)
		currentLevel := sql_interface.ModelPermissionsAll
		if dataLevel > sql_interface.ModelPermissionsAll {
			currentLevel = model.GetCurrentLevel()
		}
		if currentLevel < dataLevel {
			message, _ = utils.GetMessage("incorrect-permissions")
		} else {
			return nick, dataMap, accountIdStr, currentLevel, true, ""
		}
	}
	return "", dataMap, "", sql_interface.ModelPermissionsAll, false, message
}

// 获得数据映射表参数
func getDataMapParam() (map[string]interface{}, string, bool) {
	message := model.GetPostDataErrorMessage()
	paramObj, has := utils.GetGlobalData(middle_ware.KeyPostData)
	var dataMap map[string]interface{}
	if has {
		dataMap, _ = paramObj.(map[string]interface{})
		message = ""
	}
	return dataMap, message, has
}
