package interfaces

import (
	"fmt"
	"github.com/akwanmaroso/ddd-drugs/application"
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/infrastructure/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

type Authenticate struct {
	userInterface    application.UserAppInterface
	refreshInterface auth.AuthInterface
	tokenInterface   auth.TokenInterface
}

func NewAuthenticate(userInterface application.UserAppInterface, refreshInterface auth.AuthInterface, tokenInterface auth.TokenInterface) *Authenticate {
	return &Authenticate{userInterface: userInterface, refreshInterface: refreshInterface, tokenInterface: tokenInterface}
}

func (authenticate *Authenticate) Login(c *gin.Context) {
	var user *entity.User
	var responTokenErr = map[string]string{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json provided")
		return
	}

	validateUser := user.Validate("login")
	if len(validateUser) > 2 {
		c.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}

	u, userErr := authenticate.userInterface.GetUserByEmailAndPassword(user)
	if userErr != nil {
		c.JSON(http.StatusInternalServerError, userErr)
		return
	}

	token, tokenError := authenticate.tokenInterface.CreateToken(u.ID)
	if tokenError != nil {
		responTokenErr["token_error"] = tokenError.Error()
		c.JSON(http.StatusUnprocessableEntity, tokenError.Error())
		return
	}

	saveError := authenticate.refreshInterface.CreateAuth(u.ID, token)
	if saveError != nil {
		c.JSON(http.StatusInternalServerError, saveError.Error())
		return
	}

	userData := make(map[string]interface{})
	userData["access_token"] = token.AccessToken
	userData["refresh_token"] = token.RefreshToken
	userData["id"] = u.ID
	userData["fullname"] = u.FullName

	c.JSON(http.StatusOK, user)
}

func (authenticate *Authenticate) Logout(c *gin.Context) {
	metadata, err := authenticate.tokenInterface.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}
	deleteErr := authenticate.refreshInterface.DeleteToken(metadata)
	if deleteErr != nil {
		c.JSON(http.StatusUnauthorized, deleteErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func (authenticate *Authenticate) Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	refreshToken := mapToken["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	// check any error
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	// check token is valid
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// get uuid
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "cannot get uuid")
			return
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "error occurred")
			return
		}

		deleteError := authenticate.refreshInterface.DeleteRefresh(refreshUuid)
		if deleteError != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		token, createErr := authenticate.tokenInterface.CreateToken(userId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}

		saveError := authenticate.refreshInterface.CreateAuth(userId, token)
		if saveError != nil {
			c.JSON(http.StatusForbidden, saveError.Error())
			return
		}

		tokens := map[string]string{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh token expired")
	}
}
