package collection

import (
	"github.com/bushwood/caddyshack/adapter"
	"github.com/bushwood/caddyshack/model"
)

// Definition specifies the adapter implementaiton for a collection
type Definition struct {
	Name    string
	Adapter adapter.Definition
	Model   model.Definition
}

// GetName returns the name of the collection
func (c *Definition) GetName() string {
	return c.Name
}

// Close releases the connection back to the pool
func (c *Definition) Close() error {
	return nil
}

// Open retrieves a connection from the connection pool
func (c *Definition) Open() error {
	return nil
}
