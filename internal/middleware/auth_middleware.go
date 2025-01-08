package middleware

import (
	"net/http"

	jwt "github.com/SkandarEverest/refresh-golang/pkg/jwt"
	"github.com/spf13/viper"

	"github.com/SkandarEverest/refresh-golang/internal/exception"
	"github.com/labstack/echo/v4"
)

func Auth(config *viper.Viper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			jwtToken, err := extractJWTTokenFromHeader(ctx.Request())
			if err != nil {
				ex, ok := err.(*exception.CustomError)
				if ok {
					return ctx.JSON(ex.StatusCode, ex)
				}
				panic(err)
			}

			claim, err := jwt.ClaimToken(jwtToken, []byte(config.GetString("JWT_SECRET")))
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, exception.Unauthorized("Invalid token"))
			}

			ctx.Set("user", claim)

			// default user passing middleware if token is valid
			return next(ctx)
		}
	}
}

func extractJWTTokenFromHeader(r *http.Request) (string, error) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return "", exception.BadRequest("missing auth token")
	}

	return authToken[len("Bearer "):], nil
}
