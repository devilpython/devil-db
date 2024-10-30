package model_xml

import "encoding/xml"

//账号管理器
type AccountManager struct {
	XMLName                     xml.Name     `xml:"account-manager"`            //顶层的账号管理器名称
	AccountModel                AccountModel `xml:"account-model"`              //账号模型
}

//账号模型
type AccountModel struct {
	Id               string        `xml:"id-field,attr"`       //登录标识字段
	Password         string        `xml:"password-field,attr"` //登录密码字段
	ValidaterListObj ValidaterList `xml:"data-validate"`       //验证器列表对象
	OperatorListObj  OperatorList  `xml:"data-operation"`      //数据操作器列表对象
}
