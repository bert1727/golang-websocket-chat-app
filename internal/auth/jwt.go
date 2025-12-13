package auth

import (
	"os"
)

var jwtSecret = os.Getenv("jwtSecret")

func Auth() {}
