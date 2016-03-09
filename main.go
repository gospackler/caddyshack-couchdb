package adapter

import (
	//	"errors"
	"encoding/json"

	"fmt"

	"github.com/georgethomas111/caddyshack"
	"github.com/georgethomas111/caddyshack/model"
	"github.com/georgethomas111/caddyshack/resource"
	"github.com/georgethomas111/couchdb"
)

type CouchStore struct {
	model        *model.Definition
	client       *couchdb.Client
	StoreName    string
	DatabaseName string
	DbObj        *couchdb.Database
}

func NewCouchStore(res *resource.Definition) (couchStore *CouchStore) {

	client := couchdb.NewClient(res.Host, res.Port)
	couchStore = &CouchStore{
		client:       &client,
		StoreName:    "couchdb",
		DatabaseName: res.Name,
	}
	return
}

func (c *CouchStore) Init(model *model.Definition) (error, caddyshack.Store) {
	dbObj := c.client.DB(c.DatabaseName)
	c.model = model

	status, err := dbObj.Exists()
	if err == nil {
		if status == false {
			err = dbObj.Create()
		}
	}
	c.DbObj = &dbObj
	return err, c
}

func (c *CouchStore) GetName() string {

	return c.StoreName
}

func (c *CouchStore) SetName(name string) error {

	c.StoreName = name
	return nil
}

// TODO : This method could be part of the interface in general which can be overridden
// Does it work that way ??
func (c *CouchStore) verify(obj caddyshack.StoreObject) {

}

func (c *CouchStore) Create(obj caddyshack.StoreObject) (err error) {

	strObj, err := json.Marshal(obj)

	doc := couchdb.NewDocument("", "", c.DbObj)
	err = doc.Create(strObj)
	obj.SetKey(doc.Id)
	return
}

func (c *CouchStore) ReadOne(key string) (error, caddyshack.StoreObject) {

	doc := couchdb.NewDocument(key, "", c.DbObj)
	jsonObj, err := doc.GetDocument()

	type UpdateObj struct {
		couchdb.CouchWrapperUpdate
		TestObj
	}

	obj := &UpdateObj{}
	json.Unmarshal(jsonObj, obj)

	fmt.Println("***Json :  Object", obj)
	obj.SetKey(doc.Id)
	return err, obj
}
