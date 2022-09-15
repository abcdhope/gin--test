package errmsg

//状态码
const (
	SUCCSE = 200
	ERROR  = 500

	//用户错误1000...
	ERROR_USERNAME_USED = 1001
	//文章错误2000...
	//文章类型错误3000...
)

//定义map存放映射关系
var codeMsg = map[int]string{
	SUCCSE:              "OK",
	ERROR:               "FAIL",
	ERROR_USERNAME_USED: "用户名已被使用",
}

//输出信息
func GetErrMsg(code int) string {
	return codeMsg[code]
}
