package main

import (
	"fmt"

	"github.com/devilpython/devil-db/constants"
	"github.com/devilpython/devil-db/db/model"
	"github.com/devilpython/devil-db/db/sql_utils"
)

//weak = re.compile(r'')
//level_weak = weak.match(password)
//level_middle = re.match(r'',password)
//level_strong = re.match(r'(\w+|\W+)+',password)
//if level_weak:
//return 0
//else:
//if (level_middle and len(level_middle.group())==len(password)):
//return 1
//else:
//if level_strong and len(level_strong.group())==len(password):
//return 2

func setPointer(data *string) {
	m := "hello pointer"
	data = &m
}

func main() {
	//loginInfo := make(map[string]interface{})
	//loginInfo["id"] = "abcd123"
	//loginInfo["ticket_create_time"] = time.Now().Unix()
	//cache.SetObjectEx("test-ticket", loginInfo, 60)
	//resultInfo := make(map[string]interface{})
	//err := cache.GetObject("test-ticket", &resultInfo)
	//if err != nil {
	//	fmt.Println("..........error:", err)
	//}else  {
	//	fmt.Println("..........info:", resultInfo)
	//}

	testInsertSql()
	testUpdateSql()
	testDeleteSql()
	testQuerySql()

	//message, hasMessage := utils.GetMessage("no-ticket")
	//if hasMessage {
	//	fmt.Println(".......", message)
	//}
	//manager, hasManager := model.GetAccountManager()
	//fmt.Println(".............manager:", manager, hasManager)
}

func append(dataMap map[string]interface{}) {
	dataMap["id"] = "abc"
}

func testInsertSql() {
	tokenModel, hasModel := model.GetModel("token")
	if hasModel {
		dataMap := make(map[string]interface{})
		dataMap["account_id"] = "abcde"
		dataMap["token"] = "mytoken"
		sql := sql_utils.CreateInsertSql(tokenModel, dataMap)
		fmt.Println(".............sql:", sql)
	}
}

func testUpdateSql() {
	tokenModel, hasModel := model.GetModel("token")
	if hasModel {
		dataMap := make(map[string]interface{})
		dataMap["account_id"] = "abcde"
		dataMap["token"] = "mytoken"
		sql := sql_utils.CreateUpdateSql(tokenModel, dataMap)
		fmt.Println(".............sql:", sql)
	}
}

func testDeleteSql() {
	tokenModel, hasModel := model.GetModel("token")
	if hasModel {
		dataMap := make(map[string]interface{})
		//dataMap["account_id"] = "abcde"
		dataMap["token"] = "mytoken"
		sql := sql_utils.CreateDeleteSql(tokenModel, dataMap)
		fmt.Println(".............sql:", sql)
	}
}

func testQuerySql() {
	tokenModel, hasModel := model.GetModel("token")
	if hasModel {
		dataMap := make(map[string]interface{})
		//dataMap["fuzzy"] = true
		dataMap["account_id"] = "abcde"
		dataMap["token"] = "mytoken"
		dataMap[constants.QueryWhereCondition] = "{account_id_fuzzy} and {token}"
		sql := sql_utils.CreateQuerySql(tokenModel, dataMap)
		fmt.Println(".............sql:", sql)
	}
}
