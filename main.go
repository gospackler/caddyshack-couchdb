package adapters

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	couchColl "github.com/bushwood/caddyshack-couchdb/collection"
	"github.com/bushwood/caddyshack/collection"
	"github.com/bushwood/caddyshack/model"
	"github.com/bushwood/caddyshack/resource"
	"github.com/patrickjuchli/couch"
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
	if adp.Config.Port == "" {
		adp.Config.Port = "5984" // default to couch port
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
	config := adp.GetConfig()

	var c *couch.Credentials
	if config.Username != "" && config.Password != "" {
		c = couch.NewCredentials(config.Username, config.Password)
	} else {
		c = nil
	}

	if config.Host == "" {
		return &couchColl.Definition{}, errors.New("No host found for adapter [" + adp.Name + "]")
	}

	u := "http"
	if config.Secure == true {
		u += "s"
	}
	u += "://" + config.Host + ":" + config.Port
	s := couch.NewServer(u, c)
	db := s.Database(m.Name)
	if !db.Exists() {
		log.Info(c)
		err := db.Create()
		if err != nil {
			return &couchColl.Definition{}, err
		}
	}
	return &couchColl.Definition{m.Name, adp, m, u, c, s, db}, nil
}
