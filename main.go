package adapter

import (
	"encoding/json"

	"reflect"

	"github.com/bushwood/caddyshack"
	"github.com/bushwood/caddyshack/model"
	"github.com/bushwood/caddyshack/resource"
	"github.com/bushwood/couchdb"
)

type CouchStore struct {
	model  *model.Definition
	client *couchdb.Client
	DbObj  *couchdb.Database
	DesDoc map[string]*couchdb.DesignDoc

	// Fields for caddyshack follows
	// This is needed for identifying adpter for caddyshack.
	StoreName string
	// For all queries, receiver for data.
	ObjType reflect.Type
}

// FIXME : Assert Kind of the objModel passed is a pointer.
func NewCouchStore(res *resource.Definition, objModel caddyshack.StoreObject) (c *CouchStore) {

	client := couchdb.NewClient(res.Host, res.Port)
	c = &CouchStore{
		client:    &client,
		StoreName: "couchdb",
		ObjType:   reflect.ValueOf(objModel).Elem().Type(),
		DesDoc:    make(map[string]*couchdb.DesignDoc),
	}

	dbObj := c.client.DB(res.Name)
	status, err := dbObj.Exists()
	if err == nil {
		if status == false {
			err = dbObj.Create()
			if err != nil {
				panic("Could not create a database " + err.Error())
			}
		}
	} else {
		panic("Could not connect with db " + err.Error())
	}
	c.DbObj = &dbObj
	return
}

// FIXME Reason out and remove this method in future.
func (c *CouchStore) Init(model *model.Definition) (error, caddyshack.Store) {
	c.model = model
	return nil, c
}

func (c *CouchStore) GetDesignDoc(docName string) *couchdb.DesignDoc {

	_, exists := c.DesDoc[docName] //Checking if the view exists.
	if exists == true {
		return c.DesDoc[docName]
	} else {
		c.DesDoc[docName] = couchdb.NewDesignDoc(docName, c.DbObj)
		return c.DesDoc[docName]
	}
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

func (c *CouchStore) GetStoreObj(jsonObj []byte) (error, caddyshack.StoreObject) {

	dynmaicObj := reflect.New(c.ObjType).Interface()
	err := json.Unmarshal(jsonObj, dynmaicObj)
	if err != nil {
		return err, nil
	}

	obj := dynmaicObj.(caddyshack.StoreObject)
	return nil, obj
}

func (c *CouchStore) ReadOne(key string) (error, caddyshack.StoreObject) {

	doc := couchdb.NewDocument(key, "", c.DbObj)
	jsonObj, err := doc.GetDocument()
	if err != nil {
		return err, nil
	}

	err, obj := c.GetStoreObj(jsonObj)
	obj.SetKey(doc.Id)
	return err, obj
}

func (c *CouchStore) Read(query caddyshack.Query) (error, []caddyshack.StoreObject) {

	err, objects := query.Execute()
	// Use the rawJson to check for the view.

	return err, objects
}

// The object passed should have CouchWrapperUpdate as an anonymous field containing the details.
func (c *CouchStore) UpdateOne(obj caddyshack.StoreObject) (err error) {

	// FIXME Actually a hack which works because of the implementation.
	err = c.Create(obj)
	return
}

func (c *CouchStore) DestroyOne(key string) error {
	// Destroy not yet implemented need to implement it in the lower level. Missed it!
	return nil
}
