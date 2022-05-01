package main

import (
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"

	goessentials "github.com/ChristianHering/GoEssentials"
)

//checkPassword checks if the given password
//matches the given user's current password
func checkPassword(username string, password string) (match bool, err error) {
	if goessentials.FileNotExist(filepath.Join(configuration.DataDir, username, "password.bcrypt")) {
		return false, nil
	}

	salt, err := os.ReadFile(filepath.Join(configuration.DataDir, username, "password.salt"))
	if err != nil {
		return false, err
	}

	hashedPassword, err := os.ReadFile(filepath.Join(configuration.DataDir, username, "password.bcrypt"))
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, append(salt, []byte(password)...))
	if err != nil {
		return false, err
	}

	return true, nil
}

//setPassword sets a new password for the given user
func setPassword(username string, password string) (err error) {
	if password == "" {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configuration.DataDir, username, "password.bcrypt"), hashedPassword, 0760)
	if err != nil {
		return err
	}

	return nil
}
