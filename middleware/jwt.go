package middleware

import (
	"errors"
	"ginblogtest/routes/errmsg"
	"ginblogtest/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWT struct {
	JwtKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(utils.JwtKey),
	}
}

//验证信息
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 定义错误
var (
	TokenExpired     = errors.New("token已过期,请重新登录")
	TokenNotValidYet = errors.New("token无效,请重新登录")
	TokenMalformed   = errors.New("token不正确,请重新登录")
	TokenInvalid     = errors.New("这不是一个token,请重新登录")
)

//生成token
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey) //生成完整令牌
}

//解析token
func (j *JWT) ParserToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})

	//根据错误返回特定的错误
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}

	return nil, TokenInvalid
}

//jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//读取头部的token信息
		/*
					{
			  headers: {
			    'Authorization': 'Bearer ' + token
			  }
		*/
		var code int
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(
				http.StatusOK, gin.H{
					"status":  code,
					"message": errmsg.GetErrMsg(code),
				})
			c.Abort()
			return
		}
		//解析tokenHeader
		//'Bearer ' + token
		checkToken := strings.Split(tokenHeader, " ")
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(
				http.StatusOK, gin.H{
					"status":  code,
					"message": errmsg.GetErrMsg(code),
				})
		}
		//解析token
		j := NewJWT() //获得密钥以及解析方法
		claims, err := j.ParserToken(checkToken[1])
		if err != nil {
			//token过期
			if err == TokenExpired {
				c.JSON(http.StatusOK, gin.H{
					"status":  errmsg.ERROR,
					"message": "token授权已过期,请重新登录",
					"data":    nil,
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusOK, gin.H{
				"status":  errmsg.ERROR,
				"message": err.Error(),
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Set("username", claims)
		c.Next()
	}
}
