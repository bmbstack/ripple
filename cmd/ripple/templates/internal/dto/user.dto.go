package dto

// ReqUserInfo
// @RippleApi
// @Uri /v1/user/info
// @Method GET
type ReqUserInfo struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespUserInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// ReqUserExtra
// @RippleApi
// @Uri /v1/user/extra
// @Method GET
type ReqUserExtra struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespUserExtra struct {
	ID    uint64 `json:"id"`
	Extra string `json:"extra"`
}

// ReqUserSay
// @RippleApi
// @Uri /v1/user/say
// @Method POST
type ReqUserSay struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespUserSay struct {
	ID    uint64 `json:"id"`
	Extra string `json:"extra"`
}
