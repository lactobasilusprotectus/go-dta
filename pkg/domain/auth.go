package domain

import (
	"context"
	"go-dts/pkg/auth/common"
)

type AuthUseCase interface {
	Register(user User) (err error)
	Login(ctx context.Context, email, password string) (token common.LoginToken, err error)
	Info(ctx context.Context) (info common.LoginInfo, err error)
}
