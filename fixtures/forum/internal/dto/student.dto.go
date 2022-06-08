package dto

// ReqStudentSay
// @RippleApi
// @Uri /student/say
// @Method POST
type ReqStudentSay struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespStudentSay struct {
	ID    uint64 `json:"id"`
	Extra string `json:"extra"`
}
