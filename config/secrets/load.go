package secrets

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Init will set the secrets from the untracked .env of current app environment
// These secrets will be loaded into the os environment
func Init() error {
	path := fmt.Sprintf("config/secrets/.env.%s", os.Getenv("configEnvironment"))
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	return err
}
