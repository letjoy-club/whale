package main

import (
	"flag"
	"fmt"

	"github.com/letjoy-club/mida-tool/authenticator"
)

func main() {
	uid := flag.String("uid", "1000", "user id")
	secret := flag.String("secret", "youyue", "secret")
	flag.Parse()

	auth := authenticator.Authenticator{Key: []byte(*secret)}

	token, _ := auth.SignID(*uid)
	fmt.Println()
	fmt.Println(token)
}
