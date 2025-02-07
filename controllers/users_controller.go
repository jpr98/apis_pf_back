package controllers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jpr98/apis_pf_back/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Users represents a users controller
type Users struct {
	userStore models.UserStore
}

// NewUsersController creates a new users controller
func NewUsersController(us models.UserStore) Users {
	return Users{userStore: us}
}

// ValidateToken checks if a token is valid
func (u *Users) ValidateToken(c echo.Context) error {
	tokenString := c.Param("token")

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Valid token",
	})
}

// GetByID returns a user by a given id
func (u *Users) GetByID(c echo.Context) error {
	id := c.Param("id")
	user, err := u.userStore.GetByID(id)
	if err != nil {
		c.Logger().Errorf("Can't find user", err)
		return c.String(http.StatusNotFound, "Can't find user")
	}
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

// Create creates a user
func (u *Users) Create(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		c.Logger().Error("Can't bind body to JSON")
		return c.String(http.StatusBadRequest, "Can't bind body to json")
	}
	if !u.userStore.ValidEmail(user.Email) {
		return c.String(http.StatusConflict, "Email taken")
	}
	createdUser, err := u.userStore.Create(*user)
	if err != nil {
		c.Logger().Errorf("Can't create user", err)
		return c.String(http.StatusInternalServerError, "Can't create user")
	}

	createdUser.Password = ""
	return c.JSON(http.StatusCreated, createdUser)
}

// Update updates a user's info
func (u *Users) Update(c echo.Context) error {
	eu := new(models.EditUser)
	if err := c.Bind(eu); err != nil {
		return c.String(http.StatusBadRequest, "Can't bind body to json")
	}

	id := c.Param("id")
	if id != getTokenStringClaimByKey(c, "id") {
		return c.String(http.StatusForbidden, "You can only update your own info")
	}

	err := u.userStore.Update(id, *eu)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusAccepted, "Updated")
}

// AuthBody is the content for auth requests
type AuthBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login gives a user a token and logs the user in
func (u *Users) Login(c echo.Context) error {
	var auth AuthBody
	if err := c.Bind(&auth); err != nil {
		return c.String(http.StatusBadRequest, "Can't bind request body")
	}

	user, err := u.userStore.GetByEmail(auth.Email)
	if err != nil {
		return echo.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(72 * time.Hour).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
		"id":    user.ID.Hex(),
	})
}
