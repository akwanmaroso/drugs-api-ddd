package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type AuthInterface interface {
	CreateAuth(uint64, *TokenDetails) error
	FetchAuth(string) (uint64, error)
	DeleteRefresh(string) error
	DeleteToken(*AccessDetails) error
}

type ClientData struct {
	client *redis.Client
}

var _ AuthInterface = &ClientData{}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{client: client}
}

type AccessDetails struct {
	TokenUuid string
	UserId    uint64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func (clientData *ClientData) CreateAuth(userId uint64, tokenDetails *TokenDetails) error {
	atExpires := time.Unix(tokenDetails.AtExpires, 0)
	rtExpress := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	atCreated, err := clientData.client.Set(context.Background(), tokenDetails.TokenUuid, strconv.Itoa(int(userId)), atExpires.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := clientData.client.Set(context.Background(), tokenDetails.RefreshUuid, strconv.Itoa(int(userId)), rtExpress.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "" {
		return errors.New("no record inserted")
	}
	return nil

}

func (clientData *ClientData) FetchAuth(tokenUuid string) (uint64, error) {
	userId, err := clientData.client.Get(context.Background(), tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userId, 10, 64)
	return userID, nil
}

func (clientData *ClientData) DeleteToken(accessDetails *AccessDetails) error {
	// get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", accessDetails.TokenUuid, accessDetails.UserId)
	// delete access token
	deletedAt, err := clientData.client.Del(context.Background(), accessDetails.TokenUuid).Result()
	if err != nil {
		return err
	}
	// delete refresh token
	deleteRt, err := clientData.client.Del(context.Background(), refreshUuid).Result()
	if err != nil {
		return err
	}

	// when the record is deleted, the return value is 1
	if deletedAt != 1 || deleteRt != 1 {
		return errors.New("something went error")
	}
	return nil
}

func (clientData *ClientData) DeleteRefresh(refreshUuid string) error {
	//delete refresh token
	deleted, err := clientData.client.Del(context.Background(), refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
