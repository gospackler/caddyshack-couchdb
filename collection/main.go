package collection

import (
	"github.com/georgethomas111/caddyshack/model"
	"github.com/georgethomas111/couchdb"
)

// Definition specifies the adapter implementaiton for a collection
type Definition struct {
	Name   string
	Model  *model.Definition
	Server *couchdb.Client
	DB     *couchdb.Database
}

// GetName returns the name of the collection
func (c *Definition) GetName() string {
	return c.Name
}

// Close releases the connection back to the pool - HANDLED BY BASE ORM
func (c *Definition) Close() error {
	return nil
}

// Open retrieves a connection from the connection pool - HANDLED BY BASE ORM
func (c *Definition) Open() error {
	return nil
}
