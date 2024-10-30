package controller

import (
	"github.com/devilpython/devil-db/db/model_action"
	"github.com/devilpython/devil-db/db/sql_interface"

	devil "github.com/devilpython/devil-tools/utils"
	"github.com/gin-gonic/gin"
)

// 获得动作
func Get(context *gin.Context) {
	nick, dataMap, accountId, _, validSuccess, message := validateAccount(context, sql_interface.ModelPermissionsOperationTypeSave)
	var resultData map[string]interface{}
	if validSuccess {
		dataMap["account_id"] = accountId
		resultData, message, validSuccess = model_action.GetData(nick, dataMap)
	}
	if validSuccess && resultData != nil {
		devil.ShowDataMessage(context, true, message, resultData)
	} else {
		devil.ShowMessage(context, false, message)
	}
}
