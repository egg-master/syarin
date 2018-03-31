package datastore

import (
	"net/http"

	"github.com/mjibson/goon"
	//"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func Put(src interface{}, r *http.Request) (*datastore.Key, error) {
	n := goon.NewGoon(r)
	key, err := n.Put(src)

	ctx := appengine.NewContext(r)
	log.Debugf(ctx, "%+v %+v", key, err)
	return key, err
}

func Get(src interface{}, r *http.Request) {
	n := goon.NewGoon(r)
	n.Get(src)
}

func Delete(key *datastore.Key, r *http.Request) {
	n := goon.NewGoon(r)
	n.Delete(key)
}

func Key(src interface{}, r *http.Request) *datastore.Key {
	n := goon.NewGoon(r)
	return n.Key(src)
}

func GetMulti(res interface{}, r *http.Request) error {
	n := goon.NewGoon(r)
	return n.GetMulti(res)

}
