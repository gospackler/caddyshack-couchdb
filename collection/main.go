package collection

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bushwood/caddyshack/adapter"
	"github.com/bushwood/caddyshack/model"
	"github.com/patrickjuchli/couch"
)

// Definition specifies the adapter implementaiton for a collection
type Definition struct {
	Name    string
	Adapter adapter.Definition
	Model   model.Definition
	URL     string
	Creds   *couch.Credentials
	Server  *couch.Server
	DB      *couch.Database
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
	log.Info(c.Adapter.GetConfig())
	return nil
}
