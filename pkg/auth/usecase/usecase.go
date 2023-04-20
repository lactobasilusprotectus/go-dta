package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-dts/pkg/auth/common"
	"go-dts/pkg/common/password"
	commonTime "go-dts/pkg/common/time"
	"go-dts/pkg/domain"
	"go-dts/pkg/util/cache"
	"go-dts/pkg/util/jwt"
	"time"
)

type AuthUseCase struct {
	userRepo  domain.UserRepository
	jwtModule jwt.JwtInterface
	cache     cache.KVCacheInterface
	time      commonTime.TimeInterface
}

func NewAuthUseCase(userRepo domain.UserRepository, jwtModule jwt.JwtInterface, cache cache.KVCacheInterface,
	time commonTime.TimeInterface) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo,
		jwtModule: jwtModule,
		cache:     cache,
		time:      time,
	}
}

func (a *AuthUseCase) Register(user domain.User) (err error) {
	//hash password
	hashedPassword, err := password.HashPassword(user.Password)

	if err != nil {
		return fmt.Errorf("something wrong: %w", err)
	}

	user.Password = hashedPassword

	//save to database
	return a.userRepo.InsertUser(user)
}

func (a *AuthUseCase) Login(ctx context.Context, email, pass string) (token common.LoginToken, err error) {
	//get user from database
	user, err := a.userRepo.FindUserByEmail(email)

	if err != nil {
		return common.LoginToken{}, common.ErrEmailNotFound
	}

	//check password
	if len(user.Email) > 0 {
		if !password.CheckPasswordHash(pass, user.Password) {
			err = common.ErrWrongPassword
			return
		}

		token, err = a.generateLoginToken(ctx, user.ID)
		if err != nil {
			return
		}
	} else {
		err = common.ErrUserNotFound
		return
	}

	return
}

func (a *AuthUseCase) Info(ctx context.Context) (info common.LoginInfo, err error) {
	//TODO implement me
	panic("implement me")
}

// generate new login token with new session ID
func (a *AuthUseCase) generateLoginToken(ctx context.Context, userID int64) (token common.LoginToken, err error) {
	sessionID := uuid.New().String()

	at, err := a.generateToken(ctx, sessionID, userID, common.AccessTokenType, common.AccessTokenLifetime)
	if err != nil {
		return
	}

	rt, err := a.generateToken(ctx, sessionID, userID, common.RefreshTokenType, common.RefreshTokenLifetime)
	if err != nil {
		return
	}

	token.AccessToken = at
	token.RefreshToken = rt
	return
}

// generate new token
func (a *AuthUseCase) generateToken(ctx context.Context, sessionID string, userID int64, tokenType string,
	lifeTime time.Duration) (token string, err error) {
	return a.jwtModule.GenerateToken(ctx, jwt.JwtData{
		SessionID:  sessionID,
		IdentityID: userID,
		Type:       tokenType,
		Lifetime:   lifeTime,
	})
}
