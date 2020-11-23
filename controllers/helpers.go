package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func getTokenStringClaimByKey(c echo.Context, key string) string {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return ""
	}
	claims := user.Claims.(jwt.MapClaims)
	value, ok := claims[key].(string)
	if !ok {
		return ""
	}
	return value
}
