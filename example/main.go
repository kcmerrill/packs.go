package main

import (
	"fmt"
	plugin "github.com/kcmerrill/packs.go"
)

func main() {
	plugin.Init("kcmerrill/packs.go/", "./plugins/")
	/* Download from a specific page */
	plugin.Download("https://raw.githubusercontent.com/kcmerrill/packs.go/master/example/plugins/available/append.py")
	/* Download from our official repository */
	plugin.Download("example/plugins/available/password.py")

	/* Get going ... */
	my_password := "asd123"
	password := plugin.Filter("filter_password", my_password)
	fmt.Println("Before: ", "asd123")
	fmt.Println("Hashed: ", password)

	/* Displaying helper functions */
	fmt.Println("Enabled: ", plugin.IsEnabled("append.py"))
	fmt.Println("Enabled: ", plugin.IsEnabled("doesnotexist.py"))
}
