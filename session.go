package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

//UserSession represents a user's session
type UserSession struct {
	Session    string
	Expiration time.Time
}

//generateSession generates a new, random session
func generateSession() (session string, err error) {
	b := make([]byte, 256)

	_, err = rand.Read(b)
	if err != nil {
		return
	}

	session = base64.URLEncoding.EncodeToString(b)

	return session, nil
}

//checkSession checks to see if the given
//session is valid for the given user
func checkSession(username string, session string) (authenticated bool, err error) {
	err = refreshSessions(username)
	if err != nil {
		return false, err
	}

	b, err := os.ReadFile(filepath.Join(configuration.DataDir, username, "sessions.json"))
	if err != nil {
		return false, err
	}

	var sessions []UserSession

	err = json.Unmarshal(b, &sessions)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(sessions); i++ {
		if sessions[i].Session == session {
			return true, nil
		}
	}

	return false, nil
}

//addSession adds a session to the given user's session file
func addSession(username string, session UserSession) (err error) {
	err = refreshSessions(username)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(filepath.Join(configuration.DataDir, username, "sessions.json"))
	if err != nil {
		return err
	}

	var sessions []UserSession

	err = json.Unmarshal(b, &sessions)
	if err != nil {
		return err
	}

	sessions = append(sessions, session)

	b, err = json.MarshalIndent(sessions, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configuration.DataDir, username, "sessions.json"), b, 0660)
	if err != nil {
		return err
	}

	return nil
}

//refreshSessions deletes expired sessions
//from the given user's session file
func refreshSessions(username string) (err error) {
	b, err := os.ReadFile(filepath.Join(configuration.DataDir, username, "sessions.json"))
	if err != nil {
		return err
	}

	var sessions []UserSession

	err = json.Unmarshal(b, &sessions)
	if err != nil {
		return err
	}

	var validSessions []UserSession

	for i := 0; i < len(sessions); i++ {
		if sessions[i].Expiration.After(time.Now()) {
			validSessions = append(validSessions, sessions[i])
		}
	}

	b, err = json.MarshalIndent(validSessions, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configuration.DataDir, username, "sessions.json"), b, 0660)
	if err != nil {
		return err
	}

	return nil
}

//removeSession removes the given session
//from the given user's session file
func removeSession(username string, session string) (err error) {
	err = refreshSessions(username)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(filepath.Join(configuration.DataDir, username, "sessions.json"))
	if err != nil {
		return err
	}

	var sessions []UserSession

	err = json.Unmarshal(b, &sessions)
	if err != nil {
		return err
	}

	for i := 0; i < len(sessions); i++ {
		if sessions[i].Session == session {
			if i == len(sessions)-1 {
				sessions = sessions[:i]
			} else {
				sessions = append(sessions[:i], sessions[i+1:]...)
			}
		}
	}

	b, err = json.MarshalIndent(sessions, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(configuration.DataDir, username, "sessions.json"), b, 0660)
	if err != nil {
		return err
	}

	return nil
}
