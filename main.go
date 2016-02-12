package adapters

import (
	couchColl "github.com/bushwood/caddyshack-couchdb/collection"
	"github.com/bushwood/caddyshack/collection"
	"github.com/bushwood/caddyshack/model"
	"github.com/bushwood/caddyshack/resource"
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
	return &couchColl.Definition{m.Name, adp, m}, nil
}
