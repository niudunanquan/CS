package web

import (
	"net/http"
	"om/pkg/db"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func isAdmin(c echo.Context) bool {
	return c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims).Role == 0
}

func getCurrentUid(c echo.Context) uint {
	return c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims).Uid
}

func login(jwtKey string) func(echo.Context) error {
	return func(c echo.Context) error {
		var req struct {
			db.User
			Otp string `json:"otp"`
		}

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		password := req.Password
		otp := req.Otp

		err = req.GetByUsername()
		if err != nil {
			return c.JSON(http.StatusOK, echo.Map{
				"error": "Error username",
			})
		}

		if req.EnableOtp && !req.VerifyOtp(otp) {
			return c.JSON(http.StatusOK, echo.Map{
				"error": "Error otp",
			})
		}

		if !req.VerifyPwd(password) {
			return c.JSON(http.StatusOK, echo.Map{
				"error": "Error password",
			})
		}

		claims := &jwtCustomClaims{
			req.ID,
			req.Username,
			req.Role,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tk, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			return c.JSON(http.StatusOK, echo.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": tk,
		})
	}
}

func getUsers(c echo.Context) error {
	var user db.User
	var users []db.User
	var err error

	if isAdmin(c) {
		users, err = user.GetAll()
	} else {
		err = user.Get(getCurrentUid(c))
		if err == nil {
			users = []db.User{user}
		}
	}

	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, users)
}

func setUser(c echo.Context) error {
	var req struct {
		db.User
		Otp string `json:"otp"`
	}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	if !isAdmin(c) && getCurrentUid(c) != req.ID {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "Not allowed",
		})
	}

	err = req.HashPwd(req.Password)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	columns := map[string]interface{}{"email": req.Email, "password": req.Password}
	if isAdmin(c) {
		columns["role"] = req.Role
	}
	if req.EnableOtp && req.Otp != "" {
		if !req.VerifyOtp(req.Otp) {
			return c.JSON(http.StatusOK, echo.Map{
				"error": "error otp",
			})
		}
		columns["enable_otp"] = true
	}
	err = req.Updates(columns)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func addUser(c echo.Context) error {
	if !isAdmin(c) {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "Not allowed",
		})
	}

	var user db.User
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	err = user.Insert()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}

func delUsers(c echo.Context) error {
	if !isAdmin(c) {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "Not allowed",
		})
	}

	var req struct {
		Keys []uint `json:"keys"`
	}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	if slices.Contains(req.Keys, getCurrentUid(c)) {
		return c.JSON(http.StatusOK, echo.Map{
			"error": "Can't delete yourself",
		})
	}

	user := db.User{}
	err = user.Delete(req.Keys)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "OK")
}
