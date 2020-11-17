package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func getTokenStringClaimByKey(c echo.Context, key string) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}
