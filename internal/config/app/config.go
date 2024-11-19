package app

import "os"

const (
	ADDR = ":8082"
)

var (
	DB_URL = os.Getenv("DB_URL")
)
