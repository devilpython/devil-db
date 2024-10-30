package sql_utils

import (
	"fmt"

	"github.com/devilpython/devil-db/db/model"
	"github.com/devilpython/devil-db/db/xsession"
)

// 设置插入SQL映射表
func SetInsertSqlMap(nick string, sqlMap map[string]string, dataMap map[string]interface{}) {
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		sql := CreateInsertSql(modelObj, dataMap)
		if sql != "" {
			sqlMap[modelObj.Nick] = sql
		}
	}
}

// 设置更新SQL映射表
func SetUpdateSqlMap(nick string, sqlMap map[string]string, dataMap map[string]interface{}) {
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		sql := CreateUpdateSql(modelObj, dataMap)
		if sql != "" {
			sqlMap[modelObj.Nick] = sql
		}
	}
}

// 设置保存SQL映射表
func SetSaveSqlMap(nick string, sqlMap map[string]string, dataMap map[string]interface{}) {
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		if modelObj.PrimaryKey.IsPrimaryKey { //检查有没有主键
			sql := ""
			_, hasPrimaryKey := dataMap[modelObj.PrimaryKey.Name]  //检查数据是否有主键
			if hasPrimaryKey && modelObj.PrimaryKey.Name == "id" { //有主键则更新
				sql = CreateUpdateSql(modelObj, dataMap)
			} else { //没主键则插入
				sql = CreateInsertSql(modelObj, dataMap)
			}
			if sql != "" {
				sqlMap[modelObj.Nick] = sql
			}
		}
	}
}

// 设置查询SQL映射表
func SetQuerySqlMap(nick string, sqlMap map[string]string, dataMap map[string]interface{}) {
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		sql := CreateQuerySql(modelObj, dataMap)
		if sql != "" {
			sqlMap[modelObj.Nick] = sql
		}
	}
}

// 设置删除SQL映射表
func SetRemoveSqlMap(nick string, sqlMap map[string]string, dataMap map[string]interface{}) {
	modelObj, hasModel := model.GetModel(nick)
	if hasModel {
		sql := CreateDeleteSql(modelObj, dataMap)
		if sql != "" {
			sqlMap[modelObj.Nick] = sql
			if modelObj.PrimaryKey.IsPrimaryKey {
				primaryKey, hasPrimaryKey := dataMap[modelObj.PrimaryKey.Name]
				if hasPrimaryKey {
					session := xsession.GetDbSession()
					for childNick, childField := range modelObj.ChildrenField {
						childDataMap := make(map[string]interface{})
						childDataMap[childField.Name] = primaryKey
						queryMap := make(map[string]string)
						SetQuerySqlMap(childNick, queryMap, childDataMap)
						for _, sql := range queryMap {
							dbData, err := session.SQL(sql).Query().List()
							if err == nil {
								for index := range dbData {
									SetRemoveSqlMap(childNick, sqlMap, dbData[index])
								}
							} else {
								fmt.Println("............error:", err)
							}
						}
					}
				}
			}
		}
	}
}
