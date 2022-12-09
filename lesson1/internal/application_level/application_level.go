package application_level

import "myPackage/internal/domain_level"

type AppInterface interface {
	AddNewUser(userName string) error
	ChangeEnvironment(userName string, en string) error
	SearchUserByUserOrEnv(search string) (domain_level.User, error)
	SearchEnvByUserOrEnv(search string) (string, error)
}
type App struct {
	UserProvider domain_level.UserProvider
}

func (a *App) AddNewUser(userName string) error {
	// ...
}

func (a *App) ChangeEnvironment(userName string, en string) error {
	// ...
}

func (a *App) SearchUserByUserOrEnv(search string) (domain_level.User, error) {
	// ...
}

func (a *App) SearchEnvByUserOrEnv(search string) (string, error) {
	// ...
}
