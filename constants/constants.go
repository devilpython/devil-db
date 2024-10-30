package constants

const ConfigFilePath = "config/config.json"  //配置文件路径
const ModelFilePath = "config/model.xml"     //模型文件路径
const AccountFilePath = "config/account.xml" //账号文件路径
//const MessageFilePath = "config/message.xml" //配置文件路径

const FlowUpdateTime = "flow-update-%s"      //流程更新时间
const AnnVecPrefix = "ann-vec-%s"            //ANN创建索引的词向量前缀
const AnnIdPrefix = "ann-id-%s"              //ANN创建索引的ID前缀
const AnnIndexPrefix = "ann-index-%s"        //ANN的索引前缀
const AnnTimePrefix = "ann-time-%s"          //ANN的时间前缀
const ModelMap = "model-map"                 //模型映射表
const TicketKey = "ticket-%s"                //账号票据键值
const QueryWhereCondition = "query_where"    //查询的WHERE条件参数，格式:"{field1} and ({field2_fuzzy} or {field3})"
