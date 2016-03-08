package collection

import (
	log "github.com/Sirupsen/logrus"
	"github.com/georgethomas111/caddyshack/query"
)

// Read retrieves a collection from the database
func (c *Definition) Read(q query.Definition) (interface{}, error) {

	for _, value := range q.Where {
		log.Info(value)
	}
	// c.Open()
	// result, err := c.DB.Query("views", "by_email", make(map[string]interface{}))
	// log.Info(result.Rows)
	// defer c.Close()
	// return err
	return nil, nil
}

// ReadOne retrieves a set of collections from the database
func (c *Definition) ReadOne(id string) (interface{}, error) {
	return nil, nil
}
