package main

import (
	"ViewLog/back/ready"
)

func main() {
	ready.Config()
	// ready.Db()
	ready.Gin()
}
