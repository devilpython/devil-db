package model

import (
	"sync"

	"github.com/devilpython/devil-db/constants"
	"github.com/devilpython/devil-db/db/model_operation"
	"github.com/devilpython/devil-db/db/model_validater"
	"github.com/devilpython/devil-db/db/sql_interface"
	"github.com/devilpython/devil-db/global_keys"

	devil "github.com/devilpython/devil-tools/utils"
)

var config *ModelConfig
var onceModel sync.Once

// 获得模型配置信息
func GetModel(nick string) (Model, bool) {
	doLoadModel()
	model, has := config.ModelMap[nick]
	return *model, has
}

// 获得缺少主键消息
func GetMissingPrimaryKeyMessage() string {
	doLoadModel()
	return config.MissingPrimaryKeyMessage
}

// 获得请求的数据错误消息
func GetPostDataErrorMessage() string {
	doLoadModel()
	return config.PostDataErrorMessage
}

// 获得查询参数错误消息
func GetQueryParameterErrorMessage() string {
	doLoadModel()
	return config.QueryParameterErrorMessage
}

// 加载模型
func doLoadModel() {
	onceModel.Do(func() {
		config = LoadModelMap(constants.ModelFilePath)
	})
}

// 字段
type Field struct {
	Name         string //字段名
	Type         string //字段类型
	IsPrimaryKey bool   //是否是主键
	Create       string //创建类型
	Flag         int    //user-id: 1, user-password: 2
	TargetModel  string //关联的目标模型
	TargetField  string //关联的目标字段
}

// 模型
type Model struct {
	Nick             string                          //模型昵称
	TableName        string                          //表名
	ReadPermissions  sql_interface.ModelPermissions  //模型读权限
	WritePermissions sql_interface.ModelPermissions  //模型写权限
	FieldArray       []Field                         //字段数组
	ValidaterArray   []model_validater.DataValidater //验证器数组
	OperationArray   []model_operation.DataOperator  //数据操作器
	PrimaryKey       Field                           //主键字段
	ChildrenField    map[string]Field                //子模型对应的字段
}

// 模型配置结构
type ModelConfig struct {
	ModelMap                   map[string]*Model //模型映射表
	MissingPrimaryKeyMessage   string            //缺少主键消息
	PostDataErrorMessage       string            //请求的数据错误消息
	QueryParameterErrorMessage string            //查询参数错误消息
	AccountModel               Model             //账号模型
	AdminModel                 Model             //管理员模型
	TokenModel                 Model             //token模型
}

// 获得数据等级
func GetLevel(nick string, operationType int) sql_interface.ModelPermissions {
	model, hasModel := GetModel(nick)
	if hasModel {
		switch operationType {
		case sql_interface.ModelPermissionsOperationTypeQuery:
			return model.ReadPermissions
		case sql_interface.ModelPermissionsOperationTypeSave:
			return model.WritePermissions
		case sql_interface.ModelPermissionsOperationTypeRemove:
			return model.WritePermissions
		}
	}
	return sql_interface.ModelPermissionsUser
}

// 获得当前权限等级
func GetCurrentLevel() sql_interface.ModelPermissions {
	adminStatus, has := devil.GetGlobalData(global_keys.KeyIsAdmin)
	if has {
		isAdmin, ok := adminStatus.(bool)
		if ok && isAdmin {
			return sql_interface.ModelPermissionsAdmin
		}
	}
	return sql_interface.ModelPermissionsUser
}
