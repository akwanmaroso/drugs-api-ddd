package interfaces

import (
	"github.com/akwanmaroso/ddd-drugs/utils/mock"
)

var (
	userApp   mock.UserAppInterface
	drugApp   mock.DrugAppInterface
	fakeAuth  mock.AuthInterface
	fakeToken mock.TokenInterface

	userMock = NewUsers(&userApp, &fakeToken, &fakeAuth)
	drugMock = NewDrug(&drugApp, &userApp, &fakeToken, &fakeAuth)
	authMock = NewAuthenticate(&userApp, &fakeAuth, &fakeToken)
)
