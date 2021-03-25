package mock

import "github.com/akwanmaroso/ddd-drugs/infrastructure/auth"

// AuthInterface is a mock user auth interface
type AuthInterface struct {
	CreateAuthFn    func(uint64, *auth.TokenDetails) error
	FetchAuthFn     func(string) (uint64, error)
	DeleteRefreshFn func(string) error
	DeleteTokensFn  func(*auth.AccessDetails) error
}

func (auth *AuthInterface) DeleteRefresh(refreshUuid string) error {
	return auth.DeleteRefreshFn(refreshUuid)
}

func (auth *AuthInterface) DeleteTokens(accessDetails auth.AccessDetails) error {
	return auth.DeleteTokens(accessDetails)
}

func (auth *AuthInterface) FetchAuth(uuid string) (uint64, error) {
	return auth.FetchAuthFn(uuid)
}

func (auth *AuthInterface) CreateAuth(userId uint64, tokenDetails *auth.TokenDetails) error {
	return auth.CreateAuthFn(userId, tokenDetails)
}
