package main

import (
	"fmt"
	plugin "github.com/kcmerrill/packs.go"
)

func main() {
	plugin.Init("https://raw.githubusercontent.com/kcmerrill/packs.go/master/example/plugins/available", "./plugins/")
	/* Download from a specific page */
	//plugin.Download("https://raw.githubusercontent.com/kcmerrill/packs.go/master/example/plugins/available/append.py")
	/* Download from our official repository */
	//plugin.Download("password.py")

	/* Get going ... */
	my_password := "asd123"
	password := plugin.Filter("filter_password", my_password)
	fmt.Println("Before: ", "asd123")
	fmt.Println("Hashed: ", password)

	/* Displaying helper functions */
	fmt.Println("Enabled: ", plugin.IsEnabled("append"))
	fmt.Println("Enabled: ", plugin.IsEnabled("doesnotexist"))

	/* A quick message about downloading and installing plugins */
	fmt.Println("Now ... run this program again. This time, use --download-plugin append")
	fmt.Println("Now ... run this program again. This time, use --download-plugin password")
	fmt.Println("Now ... run this program again. This time, use --download-plugin append2")
}
