package echomiddleware

import (
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

const SECRET = "secret"
const COOKIE_NAME = "user"

func InitJWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContinueOnIgnoredError: true,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(CustomClaims)
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)

			claims := user.Claims.(*CustomClaims)

			c.Set("userId", claims.Id)

		},
		ErrorHandler: func(c echo.Context, err error) error {
			userID := rand.Int()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
				},
				Id: userID,
			})

			tokenString, err := token.SignedString([]byte(SECRET))
			if err != nil {
				return err
			}

			cookie := http.Cookie{
				Name:    "user",
				Value:   tokenString,
				Secure:  true,
				Expires: time.Now().Add(3 * time.Hour),
			}

			c.SetCookie(&cookie)
			c.Set("userId", userID)
			return nil
		},
		TokenLookup: "cookie:user",
		SigningKey:  []byte(SECRET),
	})
}

func GetUser(c echo.Context) (int, error) {
	a := c.Get("userId").(int)
	if a == 0 {
		return 0, errors.New("Отсутствует userId")
	}
	return a, nil
}
