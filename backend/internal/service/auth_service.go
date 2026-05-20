package service

import (
	"context"
	"errors"
	"time"

	"github.com/linximanager/backend/internal/config"
	"github.com/linximanager/backend/internal/model"
	jwtpkg "github.com/linximanager/backend/internal/pkg/jwt"
	"github.com/linximanager/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepo
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepo, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *model.User `json:"user"`
}

// LoginByPhone 手机号+密码登录（验证码逻辑在生产中替换）
func (s *AuthService) LoginByPhone(ctx context.Context, phone, code string) (*LoginResult, error) {
	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 自动注册
			user = &model.User{
				Phone:  phone,
				Role:   "consumer",
				Status: "active",
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	if user.Status == "frozen" {
		return nil, errors.New("account frozen")
	}
	return s.buildTokens(user)
}

// LoginByWechat 微信授权登录
func (s *AuthService) LoginByWechat(ctx context.Context, openid, unionid, nickname, avatarURL string, role string) (*LoginResult, error) {
	user, err := s.userRepo.FindByOpenid(ctx, openid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &model.User{
				Openid:    openid,
				Unionid:   unionid,
				Nickname:  nickname,
				AvatarURL: avatarURL,
				Role:      role,
				Status:    "active",
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	if user.Status == "frozen" {
		return nil, errors.New("account frozen")
	}
	return s.buildTokens(user)
}

// RefreshToken 刷新 Access Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error) {
	uid, err := jwtpkg.ParseRefreshToken(refreshToken, s.cfg.JWT.RefreshSecret)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user.Status == "frozen" {
		return nil, errors.New("account frozen")
	}
	return s.buildTokens(user)
}

func (s *AuthService) buildTokens(user *model.User) (*LoginResult, error) {
	var mid int64
	accessToken, err := jwtpkg.GenerateAccessToken(user.ID, user.Role, mid, s.cfg.JWT.AccessSecret, s.cfg.JWT.AccessExpire)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwtpkg.GenerateRefreshToken(user.ID, s.cfg.JWT.RefreshSecret, s.cfg.JWT.RefreshExpire)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// ensure compile-time usage of time package
var _ = time.Now
