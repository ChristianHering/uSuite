package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"os/user"
	"strconv"

	"github.com/gorilla/mux"
	//This is where you should import your uSuite extension packages
	huna "github.com/ChristianHering/Huna"
)

var templates *template.Template

var callbackTemplatesHTML string

func init() {
	templates = template.Must(template.ParseGlob("./templates/*.html"))
}

func main() {
	if !configuration.DevMode {
		uid := os.Getuid()

		user, err := user.LookupId(strconv.Itoa(uid))
		if err != nil {
			panic(err)
		}

		if user.Username != "usuite" {
			panic(`for correct permissions, please run usuite as the user "usuite"`)
		}

		if configuration.FSCheck {
			go checkFilePermissions()
		}
	}

	mux := mux.NewRouter()

	//This is where you should call your uSuite extension handlers
	err := callbackTemplate(huna.Huna(mux, configuration.DataDir, callback))
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	mux.HandleFunc("/userAdd", userAdditionHandler)
	mux.HandleFunc("/userRename", userRenameHandler)
	mux.HandleFunc("/userDelete", userDeletionHandler)

	if configuration.HTTPS.Enabled {
		panic(http.ListenAndServeTLS(configuration.ListenAddr, configuration.HTTPS.CertFile, configuration.HTTPS.KeyFile, nil))
	} else {
		panic(http.ListenAndServe(configuration.ListenAddr, mux))
	}
}

func callback(f string) interface{} {
	switch f {
	case "IsUserAuthenticated":
		return IsUserAuthenticated
	case "IsUserAuthenticatedAsAdmin":
		return IsUserAuthenticatedAsAdmin
	}

	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	authenticated, username, err := IsUserAuthenticated(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !authenticated {
		err = templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	adminAuthenticated, _, err := IsUserAuthenticated(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	data := struct {
		Username              string
		CallbackTemplatesHTML template.HTML
		IsAdmin               bool
	}{
		Username:              username,
		CallbackTemplatesHTML: template.HTML(callbackTemplatesHTML),
		IsAdmin:               adminAuthenticated,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func callbackTemplate(t *template.Template) (err error) {
	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
	if err != nil {
		return err
	}

	callbackTemplatesHTML = callbackTemplatesHTML + buf.String()

	return nil
}
