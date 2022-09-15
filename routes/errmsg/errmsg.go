package errmsg

//状态码
const ()

//定义map存放映射关系
var codeMsg map[int]string

//输出信息
func GetErrMsg(code int) string {
	return codeMsg[code]
}
