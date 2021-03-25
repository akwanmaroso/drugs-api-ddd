package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akwanmaroso/ddd-drugs/domain/entity"
	"github.com/akwanmaroso/ddd-drugs/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Login_Success(t *testing.T) {
	userApp.GetUserByEmailAndPasswordFn = func(user *entity.User) (*entity.User, map[string]string) {
		return &entity.User{
			ID:       1,
			FullName: "johndoe",
		}, nil
	}
	fakeToken.CreateTokenFn = func(userId uint64) (*auth.TokenDetails, error) {
		return &auth.TokenDetails{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			TokenUuid:    "8b2c2993-5eb6-4e82-b65e-a3c2ec4ddc61",
			RefreshUuid:  "7fb37b26-984e-4f8e-a7ab-fbf0f92303e6",
			AtExpires:    12345,
			RtExpires:    123455,
		}, nil
	}
	fakeAuth.CreateAuthFn = func(u uint64, details *auth.TokenDetails) error {
		return nil
	}

	fakeLoginInput := `{"email":"johndoe@example.com", "password":"12345678"}`
	r := gin.Default()
	r.POST("/login", authMock.Login)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(fakeLoginInput))
	if err != nil {
		t.Errorf("Got some error: %v\n", err.Error())
	}

	resRecord := httptest.NewRecorder()
	r.ServeHTTP(resRecord, req)

	fmt.Println("The response: ", string(resRecord.Body.Bytes()))

	response := make(map[string]interface{})

	err = json.Unmarshal(resRecord.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Got some error when unmarshalling: %v\n", err.Error())
	}

	assert.Equal(t, resRecord.Code, http.StatusOK)
	assert.EqualValues(t, response["access_token"], "access-token")
	assert.EqualValues(t, response["refresh_token"], "refresh-token")
	assert.EqualValues(t, response["fullname"], "johndoe")
}
