package dto

// ReqUserInfo ripple gen
// @Path /v1/user/info
// @Method GET
type ReqUserInfo struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespUserInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}
