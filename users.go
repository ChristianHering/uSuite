package main

import (
	"errors"
	"os"
	"path/filepath"

	goessentials "github.com/ChristianHering/GoEssentials"
)

var errUsernameTaken = errors.New("failed to change username - username taken")

//addUser adds a user data folder in the root of our
//main data directory with the correct permissions
func addUser(userName string, userPass string) error {
	err := os.Mkdir(filepath.Join(configuration.DataDir, userName), 0760)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configuration.DataDir, userName, "sessions.json"), []byte{0x5b, 0x5d}, 0760)
	if err != nil {
		return err
	}

	//_ = changeOwnership(filepath.Join(configuration.DataDir, userName))

	err = setPassword(userName, userPass)
	if err != nil {
		return err
	}

	return nil
}

//renameUser moves a user's existing data folder
//to a folder with a new name, or returns an error
//if a folder or file already exists with that name
func renameUser(oldUserName string, newUserName string) error {
	if goessentials.FolderExists(newUserName) {
		return errUsernameTaken
	}

	return os.Rename(oldUserName, newUserName)
}

//deleteUser deletes all the data in the given user's folder
func deleteUser(userName string) error {
	return os.RemoveAll(filepath.Join(configuration.DataDir, userName))
}
