package dotenv

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

var env = flag.Bool("env", true, "load environment using .env files")

// Load loads the current directory's .env file into the current process
// provided the `env` flag is true.
func Load() {
	if !*env {
		return
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("warn: %v", err)
	}
}
