package utils

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/stevenleeg/gobb/config"
	"github.com/stevenleeg/gobb/models"
	"net/http"
)

var Store *sessions.CookieStore

func GetCookieStore(r *http.Request) *sessions.CookieStore {
	if Store == nil {
		cookie_key, _ := config.Config.GetString("gobb", "cookie_key")
		Store = sessions.NewCookieStore([]byte(cookie_key))
	}

	return Store
}

func GetCurrentUser(r *http.Request) *models.User {
	cached := context.Get(r, "user")
	if cached != nil {
		return cached.(*models.User)
	}

	session, _ := GetCookieStore(r).Get(r, "sirsid")

	username := session.Values["username"]
	if username == nil || username == "" {
		return nil
	}
	password := session.Values["password"]
	if password == nil || password == "" {
		return nil
	}
	err, current_user := models.AuthenticateUser(username.(string), password.(string))

	if err != nil {
		return nil
	}

	context.Set(r, "user", current_user)
	return current_user
}
