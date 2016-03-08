package collection

import (
	"fmt"
)

// Create creates a new record in the database
func (c *Definition) Create(obj interface{}) error {
	fmt.Println("Trying to create ", obj)
	return nil
}
