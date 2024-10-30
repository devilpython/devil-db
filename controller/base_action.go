package controller

import (
	"github.com/devilpython/devil-tools/utils"
	"github.com/gin-gonic/gin"
)

// 显示版本号
func ShowVersion(context *gin.Context) {
	context.JSON(200, gin.H{
		"successful": true,
		"message":    "ok",
		"version":    "flow-db-2020-02-18-v1",
	})
}

func ShowFavicon(context *gin.Context) {
	utils.ShowMessage(context, true, "ICON")
}
