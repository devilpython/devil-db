package model_xml

import (
	"encoding/xml"
)

type ModelList struct {
	XMLName                    xml.Name `xml:"model-list"`                 //顶层的模型列表名称
	MissingPrimaryKeyMessage   string   `xml:"missing-primary-key,attr"`   //缺少主键消息
	PostDataErrorMessage       string   `xml:"post-data-error,attr"`       //请求的数据错误消息
	QueryParameterErrorMessage string   `xml:"query-parameter-error,attr"` //查询参数错误消息
	ModelArray []Model `xml:"model"` //模型数组
}

//模型
type Model struct {
	Nick             string        `xml:"nick,attr"`        //模型昵称
	TableName        string        `xml:"table-name,attr"`  //表名
	ReadPermissions  string        `xml:"read,attr"`        //模型读权限
	WritePermissions string        `xml:"write,attr"`       //模型写权限
	FieldArray       []Field       `xml:"field-list>field"` //字段数组
	ValidaterListObj ValidaterList `xml:"data-validate"`    //验证器列表对象
	OperatorListObj  OperatorList  `xml:"data-operation"`   //数据操作器列表对象
}

//验证器列表
type ValidaterList struct {
	NilValidaterArray      []NilValidater        `xml:"nil-validate"`       //空验证器数组
	LengthValidaterArray   []LengthValidater     `xml:"length-validate"`    //长度验证器数组
	RegexValidaterArray    []RegexValidater      `xml:"regex-validate"`     //正则验证器数组
	ExistValidaterArray    []IsNotExistValidater `xml:"exist-validate"`     //存在验证器数组
	NotExistValidaterArray []IsNotExistValidater `xml:"not-exist-validate"` //不存在验证器数组
}

//操作器列表
type OperatorList struct {
	DataShieldArray    []DataShield    `xml:"data-shield"`   //数据屏蔽器数组
	DataReviserArray   []DataReviser   `xml:"data-revise"`   //数据修改器数组
	DataPaddingArray   []DataPadding   `xml:"data-padding"`  //数据填充器数组
	DataExchangerArray []DataExchanger `xml:"data-exchange"` //数据转换器数组
}

//字段
type Field struct {
	Name           string `xml:"name,attr"`          //字段名
	Type           string `xml:"type,attr"`          //字段类型
	TargetModel    string `xml:"target-model,attr"`  //字段名
	TargetField    string `xml:"target-field,attr"`  //字段类型
	IsPrimaryKey   bool   `xml:"primary-key,attr"`   //是否主键
	Create         string `xml:"create,attr"`        //创建函数
	IsUserId       bool   `xml:"user-id,attr"`       //是否是用户标识
	IsUserPassword bool   `xml:"user-password,attr"` //是否用户密码
}

//===========================验证器=============================

//空验证器
type NilValidater struct {
	FieldName string `xml:"name,attr"`    //字段名
	Message   string `xml:"message,attr"` //错误消息
	ForAction string `xml:"for,attr"`     //为了什么而验证，比如为了保存或为了删除验证
}

//长度验证器
type LengthValidater struct {
	FieldName string `xml:"name,attr"`       //字段名
	MinLength int    `xml:"min-length,attr"` //最小长度
	MaxLength int    `xml:"max-length,attr"` //最大长度
	Message   string `xml:"message,attr"`    //错误消息
	ForAction string `xml:"for,attr"`        //为了什么而验证，比如为了保存或为了删除验证
}

//正则表达式验证器
type RegexValidater struct {
	FieldName string `xml:"name,attr"`    //字段名
	Regex     string `xml:"regex,attr"`   //正则表达式
	Message   string `xml:"message,attr"` //错误消息
	ForAction string `xml:"for,attr"`     //为了什么而验证，比如为了保存或为了删除验证
}

//验证字段
type ValidateField struct {
	Name        string `xml:"name,attr"`         //字段名
	TargetField string `xml:"target-field,attr"` //目标字段
}

//是否存在验证器
type IsNotExistValidater struct {
	TargetModel    string          `xml:"target-model,attr"`    //目标模型
	ConditionField string          `xml:"condition-field,attr"` //条件字段
	ConditionValue string          `xml:"condition-value,attr"` //条件值
	Message        string          `xml:"message,attr"`         //错误消息
	FieldArray     []ValidateField `xml:"validate-field"`       //验证的字段数组
	ForAction      string          `xml:"for,attr"`             //为了什么而验证，比如为了保存或为了删除验证
}

//===========================操作器=============================

//数据屏蔽器
type DataShield struct {
	Name      string `xml:"name,attr"` //字段名
	ForAction string `xml:"for,attr"`  //为了什么而屏蔽，比如为了查询或为了保存而屏蔽数据
}

//数据修改器
type DataReviser struct {
	Name      string `xml:"name,attr"`   //字段名
	Method    string `xml:"method,attr"` //修改方法
	ForAction string `xml:"for,attr"`    //为了什么而修改，比如为了查询或为了保存而修改数据
}

//数据填充器
type DataPadding struct {
	Name      string `xml:"name,attr"`   //新增字段名
	Method    string `xml:"method,attr"` //填充方法
	Param     string `xml:"param,attr"`  //填充方法用到的参数
	ForAction string `xml:"for,attr"`    //为了什么而填充，比如为了查询或为了保存而填充数据
}

//数据转换器
type DataExchanger struct {
	Name           string `xml:"name,attr"`                  //字段名
	ConditionField string `xml:"condition-field,attr"`       //字段名
	ConditionValue string `xml:"condition-value,attr"`       //字段名
	ExchangeData   string `xml:",chardata" json:"omitempty"` //转换数据
	ForAction      string `xml:"for,attr"`                   //为了什么而转换，比如为了查询或为了保存而转换数据
}
