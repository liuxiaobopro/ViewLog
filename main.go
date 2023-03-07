package main

import (
	"ViewLog/back/ready"
)

func main() {
	Init()
}

func Init() {
	ready.Config()
	ready.Db()
	ready.Gin()
}
