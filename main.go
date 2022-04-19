package main

import "heyui/server"

// @title Users API
// @version 1.0
// @description This is a series of APIs for managing users. Use JWT as access token for each API, except sign in and sign up.
// @host http://localhost:8080 https://localhost:8443
// @BasePath /
func main() {
	server.Start()
}
