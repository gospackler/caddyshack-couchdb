package adapters

import (
	"github.com/bushwood/caddyshack/resource"
)

// DBNAME isa contant used for interal database intantiation
const DBNAME string = "_CADDYSHACK"

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

// Close releases the connection back to the pool
func (adp *Definition) Close() error {
	return nil
}

// Open retrieves a connection from the connection pool
func (adp *Definition) Open() error {
	return nil
}
