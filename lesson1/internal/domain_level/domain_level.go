package domain_level

type User struct {
	Name         string
	Environments []string
}

type UserProvider interface {
	AddNewUser(userName string) error
	ChangeEnvironment(userName string, en string) error
	SearchUserByUserOrEnv(search string) (User, error)
	SearchEnvByUserOrEnv(search string) (string, error)
}

func AddNewUser(userName string) error {
	// ...
}

func ChangeEnvironment(userName string, en string) error {
	// ...
}

func SearchUserByUserOrEnv(search string) (User, error) {
	// ...
}

func SearchEnvByUserOrEnv(search string) (string, error) {
	// ...
}
