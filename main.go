package adapters

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	couchColl "github.com/bushwood/caddyshack-couchdb/collection"
	"github.com/bushwood/caddyshack/collection"
	"github.com/bushwood/caddyshack/model"
	"github.com/bushwood/caddyshack/resource"
	"github.com/bushwood/couchdb"
)

// Adapter exports the struct instance of the adapter
var Adapter = &Definition{
	Name: "couchdb",
}

// Definition defines the implementation of the adapter interface
type Definition struct {
	Name   string
	Config resource.Definition
}

// GetConfig returns the config resource of the adapter
func (adp *Definition) GetConfig() (rsc resource.Definition) {
	return adp.Config
}

// SetConfig sets the config resource of the adapter
func (adp *Definition) SetConfig(rsc resource.Definition) error {
	adp.Config = rsc
	if adp.Config.Host == "" {
		return errors.New("No host found for adapter [" + adp.Name + "]")
	}
	if adp.Config.Port == 0 {
		adp.Config.Port = 5984 // default to couch port
	}
	return nil
}

// GetName returns the name of the adapter
func (adp *Definition) GetName() string {
	return adp.Name
}

// SetName sets the name resource of the adapter
func (adp *Definition) SetName(name string) error {
	adp.Name = name
	return nil
}

// BuildCollection builds a collection for the adapter for teh provided model
func (adp *Definition) BuildCollection(m model.Definition) (collection.Definition, error) {
	var fallback collection.Definition
	config := adp.GetConfig()

	couch := couchdb.NewClient(config.Host, config.Port)
	couch.SetAuth(config.Username, couch.Password)
	couch.SetTimeout(config.Timeout)
	db := couch.DB(m.Name)
	exists, eErr := db.Exists()
	if eErr != nil {
		log.Info("fuck")
		return fallback, eErr
	}
	if !exists {
		cErr := db.Create()
		log.Info(cErr)
		if cErr != nil {
			return fallback, cErr
		}
	}

	return &couchColl.Definition{}, nil
	return &couchColl.Definition{m.Name, adp, m, &couch, &db}, nil
}
