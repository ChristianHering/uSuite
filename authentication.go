package main

import (
	"net/http"
	"path/filepath"

	goessentials "github.com/ChristianHering/GoEssentials"
	"gopkg.in/alessio/shellescape.v1"
)

//IsUserAuthenticated returns weather the user is
//authenticated and is an administrative user or not
func IsUserAuthenticatedAsAdmin(r *http.Request) (authenticated bool, username string, err error) {
	cookie, err := r.Cookie("username")
	if err != nil {
		return false, "", nil
	}

	username = shellescape.Quote(cookie.Value)

	if goessentials.FileNotExist(filepath.Join(configuration.DataDir, username, "isAdministrator")) {
		return false, username, nil
	}

	return IsUserAuthenticated(r)
}

//IsUserAuthenticated returns weather the user is authenticated
func IsUserAuthenticated(r *http.Request) (authenticated bool, username string, err error) {
	cookie, err := r.Cookie("username")
	if err != nil {
		return false, "", nil
	}

	username = shellescape.Quote(cookie.Value)

	if goessentials.FolderNotExist(filepath.Join(configuration.DataDir, username)) {
		println("poke")
		return false, username, nil
	}

	if goessentials.FileNotExist(filepath.Join(configuration.DataDir, username, "password.bcrypt")) {
		return true, username, nil
	}

	cookie, err = r.Cookie("session")
	if err != nil {
		return false, "", nil
	}

	session := shellescape.Quote(cookie.Value)

	authenticated, err = checkSession(username, session)
	if err != nil {
		return false, username, err
	}

	return authenticated, username, nil
}
