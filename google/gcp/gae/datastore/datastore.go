package datastore

import (
	"net/http"
	"reflect"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"
)

type Datastore struct {
	ctx context.Context
}

func NewDatastore(r *http.Request) Datastore {
	return Datastore{
		ctx: appengine.NewContext(r),
	}
}

type Pass struct {
	Master    string
	Secondary string
	Onetime   string
}

type User struct {
	Name string
	Mail []string
	Pass Pass
	Ver  int

	w http.ResponseWriter
	r *http.Request
}

func (p *Datastore) Get(id int64, kind string, result interface{}) error {
	key := datastore.NewKey(p.ctx, kind, "", id, nil)
	return datastore.Get(p.ctx, key, result)
}

func (p *Datastore) Post(kind string, d interface{}) (int64, error) {
	// interface > struct for datastore put
	dd := reflect.Indirect(reflect.ValueOf(d))
	//dd := reflect.ValueOf(&d)

	log.Debugf(p.ctx, "dd %+v", dd)

	key := datastore.NewIncompleteKey(p.ctx, kind, nil)
	//key, err := datastore.Put(p.ctx, key, e)
	key, err := datastore.Put(p.ctx, key, &dd)
	if err != nil {
		return 0, err
	}
	return key.IntID(), err
}

func (p *Datastore) Put(id int64, kind string, d *interface{}) (int64, error) {
	key := datastore.NewKey(p.ctx, kind, "", id, nil)
	key, err := datastore.Put(p.ctx, key, d)
	return key.IntID(), err
}

func (p *Datastore) Delete(id int64, kind string) error {
	key := datastore.NewKey(p.ctx, kind, "", id, nil)
	return datastore.Delete(p.ctx, key)
}
