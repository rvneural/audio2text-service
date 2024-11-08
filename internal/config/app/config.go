package app

import "os"

const (
	ADDR = ":8082"
)

var (
	BEARER_KEY = os.Getenv("KEY")
)
