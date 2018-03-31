package user

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func Login(r *http.Request, login, logout string) (bool, string, error) {
	var ok bool
	var url string
	var err error

	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, err = user.LoginURL(ctx, login)
		ok = true
	} else {
		url, err = user.LogoutURL(ctx, logout)
		ok = false
	}
	return ok, url, err
}
