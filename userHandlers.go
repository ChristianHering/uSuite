package main

import (
	"net/http"
	"path/filepath"

	"gopkg.in/alessio/shellescape.v1"
)

//userAdditionHandler adds a new user to the system
func userAdditionHandler(w http.ResponseWriter, r *http.Request) {
	authenticated, _, err := IsUserAuthenticatedAsAdmin(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if len(r.Form["username"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	username := shellescape.Quote(r.Form["username"][0])
	password := shellescape.Quote(r.Form["password"][0])

	err = addUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//userRenamingHandler renames a user's folder,
//effectively changing their username if available
func userRenameHandler(w http.ResponseWriter, r *http.Request) {
	authenticated, username, err := IsUserAuthenticated(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if len(r.Form["username"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	newUsername := shellescape.Quote(r.Form["username"][0])

	err = renameUser(filepath.Join(configuration.DataDir, username), filepath.Join(configuration.DataDir, newUsername))
	if err == errUsernameTaken {
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: newUsername,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//userDeletionHandler deletes a user and all of their data
func userDeletionHandler(w http.ResponseWriter, r *http.Request) {
	authenticated, username, err := IsUserAuthenticated(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	err = deleteUser(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
