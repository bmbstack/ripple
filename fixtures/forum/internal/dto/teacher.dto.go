package dto

// ReqTeacherTeach
// @RippleApi
// @Uri /teacher/teach
// @Method GET
type ReqTeacherTeach struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespTeacherTeach struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// ReqTeacherHello
// @RippleApi
// @Uri /teacher/hello
// @Method POST
type ReqTeacherHello struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespTeacherHello struct {
	ID    uint64 `json:"id"`
	Extra string `json:"extra"`
}

