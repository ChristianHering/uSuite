package main

import (
	"io/fs"
	"log"
	"os"
	"os/user"
	"strconv"
)

//checkFilePermissions makes sure the data directory
//folder has the correct ownership and permissions
//
//This doesn't mess with user directories/files in case
//you wanted to use custom groups for each user folder
//to allow for per-user access to user data folders
func checkFilePermissions() {
	err := os.MkdirAll(configuration.DataDir, 0764)
	if err != nil {
		log.Println("failed to make configuration directory with error: ", err)
	}

	err = changeOwnership(configuration.DataDir) //Change file ownership
	if err != nil {
		log.Println("failed to change ownership for ", configuration.DataDir)
	}

	err = changePermissions(configuration.DataDir, 0764) //Change file permissions
	if err != nil {
		log.Println("failed to change permissions for ", configuration.DataDir)
	}

	err = os.Mkdir(configuration.DataDir+"/administrator", 0760)
	if err != nil && !os.IsExist(err) {
		log.Println("failed to make administrator user data folder")
	}
}

//changeOwnership changes the ownership of the file
//or folder at the given path to usuite:usuite
func changeOwnership(path string) error {
	group, err := user.LookupGroup("usuite")
	if err != nil {
		log.Println("failed to lookup usuite group id with error: ", err)
		return err
	}

	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		log.Println("failed to convert gid to int")
		return err
	}

	err = os.Chown(path, os.Getuid(), gid)
	if err != nil {
		log.Println("failed to set ownership of data directory. This may or may not cause issues")
		return err
	}

	return nil
}

//changePermissions changes the permissions for the file
//or folder located at 'path' to the passed permissions
func changePermissions(path string, permissions fs.FileMode) error {
	pathInfo, err := os.Stat(path)
	if err != nil {
		log.Println("failed to get info for data directory: ", err)
		return err
	}

	if pathInfo.Mode() != permissions {
		err := os.Chmod(path, permissions)
		if err != nil {
			log.Println("chmod failed with error: ", err)
			return err
		}
	}

	return nil
}
