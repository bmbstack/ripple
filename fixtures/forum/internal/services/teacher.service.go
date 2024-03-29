// Code generated by ripple g, You can edit it again.
// source: fixtures/forum/internal/dto/teacher.dto.go

package services

import (
	"github.com/labstack/echo/v4"
	"sync"
)

var (
	teacherService     *TeacherService
	teacherServiceOnce sync.Once
)

// GetTeacherService 单例
func GetTeacherService(ctx echo.Context) *TeacherService {
	teacherServiceOnce.Do(func() {
		teacherService = &TeacherService{}
	})
	teacherService.ctx = ctx
	return teacherService
}

type TeacherService struct {
	ctx echo.Context
}
		