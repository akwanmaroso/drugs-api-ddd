package interfaces

import (
	"github.com/akwanmaroso/ddd-drugs/application"
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Users struct {
	userInterface    application.UserAppInterface
	tokenInterface   auth.TokenInterface
	refreshInterface auth.AuthInterface
}

func NewUsers(userInterface application.UserAppInterface, tokenInterface auth.TokenInterface, refreshInterface auth.AuthInterface) *Users {
	return &Users{userInterface: userInterface, tokenInterface: tokenInterface, refreshInterface: refreshInterface}
}

func (u *Users) SaveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	// validate request
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}

	newUser, err := u.userInterface.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (u *Users) GetUsers(c *gin.Context) {
	users := entity.Users{}
	var err error
	users, err = u.userInterface.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (u *Users) GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := u.userInterface.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}
