package dto

// ReqUserInfo
// @RippleApi
// @Uri /user/info
// @Method GET
// @Version v1
type ReqUserInfo struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespUserInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// ReqUserExtra
// @RippleApi
// @Uri /user/extra
// @Method GET
// @Version v1
type ReqUserExtra struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespUserExtra struct {
	ID    uint64 `json:"id"`
	Extra string `json:"extra"`
}

// ReqUserSay
// @RippleApi
// @Uri /user/say
// @Method POST
// @Version v1
type ReqUserSay struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespUserSay struct {
	ID    uint64 `json:"id"`
	Extra string `json:"extra"`
}
