package controller

import (
	"github.com/devilpython/devil-db/db/model_action"
	"github.com/devilpython/devil-db/db/sql_interface"
	"github.com/devilpython/devil-db/global_keys"

	devil "github.com/devilpython/devil-tools/utils"
	"github.com/gin-gonic/gin"
)

// 删除动作
func Remove(context *gin.Context) {
	nick, dataMap, accountId, _, validSuccess, message := validateAccount(context, sql_interface.ModelPermissionsOperationTypeSave)
	if validSuccess {
		dataMap["account_id"] = accountId
		message, validSuccess = model_action.RemoveData(nick, dataMap)
	}
	if validSuccess {
		id, hasId := dataMap["id"]
		if hasId {
			devil.ShowIdMessage(context, true, message, id)
		} else {
			devil.SetGlobalData(global_keys.DbMustCommit, true)
			devil.ShowMessage(context, true, "ok")
		}
	} else {
		devil.ShowMessage(context, false, message)
	}
}
