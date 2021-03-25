package mock

import (
	"github.com/akwanmaroso/ddd-drugs/infrastructure/auth"
	"net/http"
)

// TokenInterface is a mock user token interface
type TokenInterface struct {
	CreateTokenFn          func(userId uint64) (*auth.TokenDetails, error)
	ExtractTokenMetadataFn func(*http.Request) (*auth.AccessDetails, error)
}

func (t *TokenInterface) CreateToken(userId uint64) (*auth.TokenDetails, error) {
	return t.CreateTokenFn(userId)
}

func (t *TokenInterface) ExtractTokenMetadata(r *http.Request) (*auth.AccessDetails, error) {
	return t.ExtractTokenMetadataFn(r)
}
