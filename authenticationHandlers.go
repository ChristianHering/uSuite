package main

import (
	"net/http"
	"path/filepath"
	"time"

	goessentials "github.com/ChristianHering/GoEssentials"
	"gopkg.in/alessio/shellescape.v1"
)

//loginHandler logs in the user and redirects them to /
func loginHandler(w http.ResponseWriter, r *http.Request) {
	authenticated, _, err := IsUserAuthenticated(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if len(r.Form["username"]) == 0 || len(r.Form["password"]) == 0 {
		err := templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	username := shellescape.Quote(r.Form["username"][0])
	password := shellescape.Quote(r.Form["password"][0])

	if goessentials.FolderNotExist(filepath.Join(configuration.DataDir, username)) {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	if goessentials.FileExists(filepath.Join(configuration.DataDir, username, "password.bcrypt")) {
		match, err := checkPassword(username, password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if !match {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}
	}

	sessionString, err := generateSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	sessionExpiration := time.Now().Add(configuration.SessionTTL)

	session := UserSession{
		Session:    sessionString,
		Expiration: sessionExpiration,
	}

	err = addSession(username, session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   username,
		Expires: sessionExpiration,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   sessionString,
		Expires: sessionExpiration,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//logoutHandler logs out the user and redirects them to /
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	authenticated, username, err := IsUserAuthenticated(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if !authenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	sessionCookie, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = removeSession(username, sessionCookie.Value)
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
