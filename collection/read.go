package collection

// Read retrieves a collection from the database
func (c *Definition) Read() error {
	c.Open()
	defer c.Close()
	return nil
}

// ReadOne retrieves a set of collections from the database
func (c *Definition) ReadOne(id string) error {
	return nil
}
