package echomiddleware

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type CustomClaims struct {
	UserID int `json:"id"`
	jwt.RegisteredClaims
}

const Secret = "secret"
const CookieName = "user"

func InitJWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContinueOnIgnoredError: true,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(CustomClaims)
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)

			claims := user.Claims.(*CustomClaims)

			c.Set("userId", claims.UserID)

		},
		ErrorHandler: func(c echo.Context, err error) error {
			userID := rand.Int()

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
				},
				UserID: userID,
			})

			tokenString, errT := token.SignedString([]byte(Secret))
			if errT != nil {
				return errT
			}

			cookie := http.Cookie{
				Name:    "user",
				Value:   tokenString,
				Expires: time.Now().Add(3 * time.Hour),
			}

			c.SetCookie(&cookie)
			fmt.Println(cookie)
			c.Set("userId", userID)
			return nil
		},
		TokenLookup: "cookie:user",
		SigningKey:  []byte(Secret),
	})
}

func GetUser(c echo.Context) (int, error) {
	a := c.Get("userId").(int)
	if a == 0 {
		return 0, errors.New("отсутствует userId")
	}
	return a, nil
}
