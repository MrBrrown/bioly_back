package main

import "github.com/alexedwards/argon2id"

func main() {
	passwords := []string{
		"password123",
		"admin",
		"rootuser",
		"123456",
		"pass",
	}

	for _, pwd := range passwords {
		hashed, err := argon2id.CreateHash(pwd, argon2id.DefaultParams)
		if err != nil {
			panic(err)
		}
		println("Password:", pwd, "Hash:", hashed)
	}
}
