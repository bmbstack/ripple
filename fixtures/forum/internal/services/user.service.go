// Code generated by ripple g, You can edit it again.
// source: fixtures/forum/internal/dto/user.dto.go

package services

import (
	"github.com/labstack/echo/v4"
	"sync"
)

var (
	userService     *UserService
	userServiceOnce sync.Once
)

// GetUserService 单例
func GetUserService(ctx echo.Context) *UserService {
	userServiceOnce.Do(func() {
		userService = &UserService{}
	})
	userService.ctx = ctx
	return userService
}

type UserService struct {
	ctx echo.Context
}
		