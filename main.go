package main

import (
	"auditor.z9fr.xyz/server/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = bootstrap.RootApp.Execute()
}
