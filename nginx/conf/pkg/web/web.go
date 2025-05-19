package web

import (
	"io/fs"
	"net/http"
	"om/frontend"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func staticHandler() http.Handler {
	fsys, err := fs.Sub(frontend.EmbededFiles, "dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}

type jwtCustomClaims struct {
	Uid      uint   `json:"uid"`
	Username string `json:"username"`
	Role     int    `json:"role"`
	jwt.RegisteredClaims
}

func New(jwtKey string) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         3600,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().Header.Get("Accept"), "image/")
		},
		MinLength: 1024,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogHost:      true,
		LogRemoteIP:  true,
		LogMethod:    true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.WithFields(log.Fields{
					"remote_ip":  v.RemoteIP,
					"host":       v.Host,
					"method":     v.Method,
					"uri":        v.URI,
					"user_agent": v.UserAgent,
					"status":     v.Status,
					"latency":    v.Latency,
				}).Info("req")
			} else {
				log.WithFields(log.Fields{
					"remote_ip":  v.RemoteIP,
					"host":       v.Host,
					"method":     v.Method,
					"uri":        v.URI,
					"user_agent": v.UserAgent,
					"status":     v.Status,
					"latency":    v.Latency,
				}).Error(v.Error.Error())
			}
			return nil
		},
	}))

	e.GET("/*", echo.WrapHandler(staticHandler()))

	e.POST("/api/v1/login", login(jwtKey))

	ag := e.Group("/api/v1/admin")

	ag.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(jwtKey),
	}))

	ag.GET("/stats", getStats)
	ag.GET("/rts", getRts)

	ag.GET("/users", getUsers)
	ag.POST("/users", addUser)
	ag.PUT("/users", setUser)
	ag.DELETE("/users", delUsers)

	ag.GET("/sites", getSites)
	ag.POST("/sites", addSite)
	ag.PUT("/sites", setSite)
	ag.DELETE("/sites", delSites)

	ag.GET("/certs", getCerts)
	ag.POST("/certs", addCert)
	ag.PUT("/certs", setCert)
	ag.DELETE("/certs", delCerts)

	ag.GET("/upstreams", getUpstreams)
	ag.POST("/upstreams", addUpstream)
	ag.PUT("/upstreams", setUpstream)
	ag.DELETE("/upstreams", delUpstreams)

	ag.GET("/om_config", getOmConfig)
	ag.PUT("/om_config", setOmConfig)
	ag.GET("/version", getVersion)
	ag.GET("/update", update)

	return e
}
