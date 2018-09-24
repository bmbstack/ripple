package helper

const (
	CacheRedisSection = "ripple"
	CacheSeconds      = 3600 //600s
)

const (
	SuccessCode = 200 + iota //200
)

const (
	ErrorCode                     = 500 + iota //500
	ErrorCodeParamsValidateFailed              //501
)

const (
	ErrorMsgParamsValidateFailed string = "参数校验失败"
)

func SuccessJSON(data interface{}) interface{} {
	return map[string]interface{}{
		"code": SuccessCode,
		"msg":  "成功",
		"data": data,
	}
}

func ErrorJSON(msg interface{}, code int) interface{} {
	return map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": []map[string]interface{}{},
	}
}