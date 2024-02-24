package spotistats

import (
	"github.com/erfgypO/spotistats/lib/server"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func runApi() {
	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/sign-up", server.HandleSignUp)
	e.POST("/sign-in", server.HandleSignIn)

	e.GET("/auth-redirect", server.HandleAuthRedirect)

	authGroup := e.Group("/auth")
	authGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	authGroup.GET("/url", server.HandleGetAuthUrl)
	authGroup.GET("/user/me", server.HandleGetMe)
	authGroup.PUT("/user/update-password", server.HandleUpdatePassword)

	authGroup.GET("/stats", server.HandleGetStats)
	authGroup.GET("/stats/hourly", server.HandleGetStatsGroupedByHour)

	e.Logger.Fatal(e.Start(os.Getenv("API_URL")))
}
