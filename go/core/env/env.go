package env

import (
	"os"
)

var (
	Domain     string
	AppEnv     string
	JwtPublic  string
	JwtPrivate string
)

func init() {
	if v := os.Getenv("DOMAIN"); v != "" {
		Domain = v
	}
	if v := os.Getenv("APP_ENV"); v != "" {
		AppEnv = v
	}
	if v := os.Getenv("JWT_PUBLIC_KEY"); v != "" {
		JwtPublic = v
	}
	if v := os.Getenv("JWT_PRIVATE_KEY"); v != "" {
		JwtPrivate = v
	}
}
