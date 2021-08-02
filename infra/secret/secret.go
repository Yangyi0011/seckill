package secret

import (
	"github.com/dgrijalva/jwt-go"
	"seckill/infra/code"
	"strings"
	"time"
)

const (
	// SecKey 服务器密钥
	SecKey = "hello second kill"
	// ExpiresTime JWT 过期时间，单位：秒
	ExpiresTime = 60*60
	// Issuer JWT 签发人
	Issuer      = "second kill"
	// TokenPrefix JWT 生成 token 所添加的前缀
	TokenPrefix = "Bearer "
)

// CustomClaims 自定义 Claims，继承 jwt.StandardClaims 并添加一些自己需要的信息
type CustomClaims struct {
	UserId uint `json:"userId"`
	Username string `json:"username"`
	Kind     int8 `json:"kind"`
	jwt.StandardClaims
}

// JWT 结构
type JWT struct {
	// SigningKey 密钥信息
	SigningKey []byte
}

// NewJWT 创建一个 JWT 实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(SecKey),
	}
}

// CreateToken 创建 JWT
func (j *JWT) CreateToken(claims CustomClaims) (token string, err error) {
	// 通过 HS256 算法生成 tokenClaims ,这就是我们的 HEADER 部分和 PAYLOAD。
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString(j.SigningKey)
	// 返回添加了前缀的 token
	return TokenPrefix + token, err
}

// ParseToken 解析 JWT
func (j *JWT) ParseToken(tokenStr string) (*CustomClaims, error) {
	auth := strings.Fields(tokenStr)
	if len(auth) < 1 {
		return nil, code.TokenInvalidErr
	}
	// 解析 token 时去除前缀的影响
	token, err := jwt.ParseWithClaims(auth[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, code.TokenMalformedErr
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, code.TokenExpiredErr
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, code.TokenNotValidYetErr
			} else {
				return nil, code.TokenInvalidErr
			}
		}
	}
	if err == nil && token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, code.TokenInvalidErr
}

// RefreshToken 刷新 JWT
func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	auth := strings.Fields(tokenStr)
	if len(auth) < 1 {
		return "", code.TokenInvalidErr
	}
	// 解析 token 时去除前缀的影响
	token, err := jwt.ParseWithClaims(auth[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.StandardClaims.ExpiresAt = time.Now().Add(ExpiresTime * time.Second).Unix()
		return j.CreateToken(*claims)
	}
	return "", code.TokenInvalidErr
}

// InvalidToken 让 JWT 失效
// 这里实测，并没有产生应有的效果
func (j *JWT) InvalidToken(tokenStr string) error {
	auth := strings.Fields(tokenStr)
	if len(auth) < 1 {
		return code.TokenInvalidErr
	}
	// 解析 token 时去除前缀的影响
	token, err := jwt.ParseWithClaims(auth[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 使这个 JWT 立即失效
		claims.StandardClaims.ExpiresAt = time.Now().Unix()
		return nil
	}
	return code.TokenInvalidErr
}