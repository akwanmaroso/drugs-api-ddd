package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Token struct{}

func NewToken() *Token {
	return &Token{}
}

type TokenInterface interface {
	CreateToken(userId uint64) (*TokenDetails, error)
	ExtractTokenMetadata(r *http.Request) (*AccessDetails, error)
}

var _ TokenInterface = &Token{}

func (token *Token) CreateToken(userId uint64) (*TokenDetails, error) {
	tokenDetails := &TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.TokenUuid = uuid.NewV4().String()

	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetails.RefreshUuid = tokenDetails.TokenUuid + "++" + strconv.Itoa(int(userId))

	var err error
	// create access token
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = tokenDetails.TokenUuid
	accessTokenClaims["user_id"] = userId
	accessTokenClaims["exp"] = tokenDetails.AtExpires
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokenDetails.AccessToken, err = accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// create refresh token
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	refreshTokenClaims["user_id"] = userId
	refreshTokenClaims["exp"] = tokenDetails.RtExpires
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenDetails.RefreshUuid, err = refreshToken.SignedString([]byte(os.Getenv("REFRESH_TOKEN")))
	if err != nil {
		return nil, err
	}
	return tokenDetails, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check if the token method is "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func (token *Token) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	fmt.Println("WE ENTERED METADATA")
	tokenJWT, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := tokenJWT.Claims.(jwt.MapClaims)
	if ok && tokenJWT.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{TokenUuid: accessUuid, UserId: userId}, nil
	}
	return nil, err
}
