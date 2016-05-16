package main

import (
	"fmt"
	plugin "github.com/kcmerrill/packs.go"
	"os"
)

func main() {
	plugin.Init("putthexampleplugindirhere", os.Args)
	my_password := "asd123"
	password := plugin.Filter("filter_password", my_password)
	fmt.Println("Before: ", "asd123")
	fmt.Println("Hashed: ", password)
}
