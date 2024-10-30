package sql_utils

import (
	"fmt"

	"github.com/devilpython/devil-db/db/model"
	"github.com/devilpython/devil-db/db/xsession"
)

// 数据库动作对象
type Action struct {
}

// 保存模型
func (action Action) SaveModel(sqlMap map[string]string) bool {
	if len(sqlMap) > 0 {
		session := xsession.GetDbSession()
		for _, sql := range sqlMap {
			fmt.Println("......sql:", sql)
			_, err := session.Exec(sql)
			if err == nil {
				return true
			} else {
				fmt.Println("............error:", err)
				return false
			}
		}
		return true
	}
	return false
}

// 删除模型
func (action Action) DeleteModel(nick string, dataMap map[string]interface{}) bool {
	sqlMap := make(map[string]string)
	SetRemoveSqlMap(nick, sqlMap, dataMap)
	if len(sqlMap) > 0 {
		session := xsession.GetDbSession()
		for _, sql := range sqlMap {
			fmt.Println("......sql:", sql)
			result, err := session.Exec(sql)
			if err == nil {
				count, _ := result.RowsAffected()
				if count == 0 {
					return false
				}
			} else {
				fmt.Println("............error:", err)
				return false
			}
		}
		return true
	}
	return false
}

// 查询模型
func (action Action) QueryModel(nick string, dataMap map[string]interface{}) []map[string]interface{} {
	sqlMap := make(map[string]string)
	SetQuerySqlMap(nick, sqlMap, dataMap)
	if len(sqlMap) > 0 {
		session := xsession.GetDbSession()
		for _, sql := range sqlMap {
			fmt.Println("......sql:", sql)
			result, err := session.SQL(sql).Query().List()
			if err == nil {
				return result
			} else {
				fmt.Println("............error:", err)
			}
		}
	}
	return nil
}

// 查询模型
func (action Action) GetModelData(nick string, dataMap map[string]interface{}) map[string]interface{} {
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		if len(modelObj.PrimaryKey.Name) > 0 {
			data, hasData := dataMap[modelObj.PrimaryKey.Name]
			if hasData {
				paramMap := make(map[string]interface{})
				paramMap[modelObj.PrimaryKey.Name] = data
				result := action.QueryModel(nick, paramMap)
				if result != nil && len(result) == 1 {
					return result[0]
				}
			}
		}
	}
	return nil
}
