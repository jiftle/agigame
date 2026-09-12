package jwtutil

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

// SignMethod 签名算法
var SignMethod = jwt.SigningMethodHS256

// Claims 自定义载荷
type Claims struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// 已撤销 token 内存黑名单：jti -> 过期时间(Unix)
var (
	revokedMu sync.RWMutex
	revoked   = make(map[string]int64)
)

func secret(ctx context.Context) []byte {
	if s := os.Getenv("ADMINBASE_JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte(g.Cfg().MustGet(ctx, "jwt.secret", "adminbase-default-secret").String())
}

// Generate 生成 token
func Generate(ctx context.Context, userId int, username, tokenType string, expire time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:   userId,
		Username: username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        grand.S(32),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			Issuer:    "adminbase",
		},
	}
	return jwt.NewWithClaims(SignMethod, claims).SignedString(secret(ctx))
}

// Parse 解析 token
func Parse(ctx context.Context, token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return secret(ctx), nil
	}, jwt.WithValidMethods([]string{SignMethod.Alg()}))
	if err != nil {
		return nil, err
	}
	if isRevoked(claims.ID) {
		return nil, errors.New("令牌已被撤销")
	}
	return claims, nil
}

// Revoke 撤销 token（退出登录 / 令牌轮换）
func Revoke(ctx context.Context, token string) {
	if token == "" {
		return
	}
	claims, err := Parse(ctx, token)
	if err != nil || claims.ID == "" {
		return
	}
	exp := time.Now().Add(7 * 24 * time.Hour).Unix()
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Unix()
	}
	revokedMu.Lock()
	revoked[claims.ID] = exp
	revokedMu.Unlock()
	cleanupRevoked()
}

func isRevoked(id string) bool {
	if id == "" {
		return false
	}
	revokedMu.RLock()
	exp, ok := revoked[id]
	revokedMu.RUnlock()
	return ok && exp > time.Now().Unix()
}

func cleanupRevoked() {
	now := time.Now().Unix()
	revokedMu.Lock()
	for id, exp := range revoked {
		if exp <= now {
			delete(revoked, id)
		}
	}
	revokedMu.Unlock()
}
