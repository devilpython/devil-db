package xsession

import (
	"github.com/devilpython/devil-db/global_keys"

	devil "github.com/devilpython/devil-tools/utils"
	"github.com/xormplus/xorm"
)

// 获得数据库会话
func GetDbSession() *xorm.Session {
	session, _ := devil.GetGlobalData(global_keys.KeyDbSession)
	return session.(*xorm.Session)
}
