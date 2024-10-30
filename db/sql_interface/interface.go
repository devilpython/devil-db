package sql_interface

//模型权限
type ModelPermissions int

const (
	ModelPermissionsAll   ModelPermissions = 1
	ModelPermissionsUser  ModelPermissions = 2
	ModelPermissionsAdmin ModelPermissions = 3
	ModelPermissionsSys   ModelPermissions = 4
)

var ModelPermissionsOperationTypeQuery = 1         //查询操作，操作器针对查询结果
var ModelPermissionsOperationTypeSave = 2          //保存操作
var ModelPermissionsOperationTypeRemove = 4        //删除操作
var ModelPermissionsOperationTypeQueryForParam = 8 //查询操作，操作器针对查询参数

//SQL动作接口
type SqlActionInterface interface {
	//保存模型
	SaveModel(sqlMap map[string]string) bool

	//删除模型
	DeleteModel(nick string, dataMap map[string]interface{}) bool

	//查询模型
	QueryModel(nick string, dataMap map[string]interface{}) []map[string]interface{}

	//通过ID获得模型数据
	GetModelData(nick string, dataMap map[string]interface{}) map[string]interface{}
}
