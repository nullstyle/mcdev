package dotenv

import (
	"flag"

	"github.com/joho/godotenv"
)

var env = flag.Bool("env", true, "load environment using .env files")

// Load loads the current directory's .env file into the current process
// provided the `env` flag is true.
func Load() error {
	if !*env {
		return nil
	}

	return godotenv.Load()
}
