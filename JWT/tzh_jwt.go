package tzh_utils

//  自定义JWT过期时间跟Secret

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("CXVBNMYHUJKI")

type MyClaims struct {
	jwt.StandardClaims        // jwt包自带的jwt.StandardClaims只包含了官方字段
	UserName           string `json:"username"` // 自定义要存储的信息,可有多个
}

// 生成token
func GenToken(username string) (string, error) {
	c := MyClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "tzh666",                                   // 签发人
		},
		username,
	}
	/*
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
		token.SignedString(MySecret)
	*/
	// 使用指定的签名方法创建签名对象,此处指定用S256,再调用SignedString函数使用指定的secret签名并获得完整的编码后的字符串token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(MySecret)
}

// token解析
func ParseToken(Token string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(Token, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	// 校验token
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
