package config

import "fmt"

const (
	Host    = "localhost"
	Port    = 8080
	HashLen = 6
)

var (
	Addr = fmt.Sprintf("%s:%d", Host, Port)
)
