package dto

// ReqStudentLearn
// @RippleApi
// @Uri /student/learn
// @Method POST
type ReqStudentLearn struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespStudentLearn struct {
	Name string `json:"name"`
}

// ReqStudentHello
// @RippleApi
// @Uri /student/hello
// @Method POST
type ReqStudentHello struct {
	ID uint64 `form:"id" json:"id" binding:"required"`
}

type RespStudentHello struct {
	ID    uint64 `json:"id"`
	Extra string `json:"extra"`
}
