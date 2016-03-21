package adapter

import (
	"github.com/bushwood/caddyshack"
	"github.com/bushwood/caddyshack/model"
	"github.com/bushwood/caddyshack/resource"
	"github.com/bushwood/couchdb"
	"testing"
)

// Create a compatable storeObject
type TestObj struct {
	Name   string `json:"name" query:""`
	Value  string `json:"surprise"`
	Field1 string `json:"field1"`
	Age    int    `json:"age" query:"age < 20"`
	Id     string `json:"id"`
}

type RetTestObj struct {
	couchdb.CouchWrapperUpdate
	TestObj
}

var Caddy *caddyshack.Caddies
var Key string

func (t *TestObj) GetKey() string {
	return t.Id
}

func (t *TestObj) SetKey(id string) {
	t.Id = id
}

var CouchStoreIns *CouchStore

func TestInit(t *testing.T) {

	// Add model definition in future to it.
	cs := caddyshack.New()

	res := &resource.Definition{
		Host: "127.0.0.1",
		Port: 5984,
		Name: "adaptertest",
	}

	// From storedemo.go
	couchStore := NewCouchStore(res, &RetTestObj{})
	err := cs.LoadStore(couchStore)

	if err != nil {
		t.Error("Error while loading a store.")
	}

	model := &model.Definition{
		Adapter: "couchdb",
		Name:    "testModel",
	}
	err = cs.AddModel(model)
	if err != nil {
		t.Error("Error while Adding Model to caddyshack", err)
	}

	caddyName := model.Name + model.Adapter
	err, Caddy = cs.GetCaddy(caddyName)
	if err != nil {
		t.Error("Error while retreiving caddy ", err)
	}
	CouchStoreIns = couchStore
}

func TestCreate(t *testing.T) {

	testObj := &TestObj{
		Name:   "abcd",
		Value:  "1234",
		Field1: "field1",
	}

	err := Caddy.StoreIns.Create(testObj)
	if err != nil {

		t.Error("Error creating object in the test Store")
	}

	Key = testObj.GetKey()
	t.Log("Created Object with key", testObj.GetKey())
}

var RetrObject *RetTestObj

func TestReadOne(t *testing.T) {

	err, obj := Caddy.StoreIns.ReadOne(Key)
	if err != nil {
		t.Error("Error while retreiving object")
	}

	if obj.GetKey() != Key {
		t.Error("Retreived wrong object")
	}

	actualObj := obj.(*RetTestObj)
	t.Log("Got the actual object back.", actualObj.TestObj)
	RetrObject = actualObj
}

func TestUpdate(t *testing.T) {

	RetrObject.TestObj.Name = "Updated"
	RetrObject.TestObj.Value = "-1"
	err := Caddy.StoreIns.UpdateOne(RetrObject)
	if err != nil {
		t.Log("Error while updating object, ", err)
	}
	t.Log("Check the updated object in the DB if delete is disabled with key", Key)
}

func TestRead(t *testing.T) {

	// Every Query is the request to a view.

	//	NewQuery("function(doc) {emit(doc.field1);}", "new_view", "new_design", CouchStoreIns)
	query := NewQuery("function(doc) {emit(doc.field1);}", "new_view", "new_design", CouchStoreIns)
	err, objects := Caddy.StoreIns.Read(query)
	if err != nil {
		t.Error("Error while reading query ", query, " ", err)
	} else {
		t.Log("Read", objects)
	}
	for _, obj := range objects {
		t.Log(obj.GetKey())
	}
}

func TestObjQueryRead(t *testing.T) {

	newTestObj := new(TestObj)

	res := &resource.Definition{
		Host:   "127.0.0.1",
		Port:   5984,
		Name:   "adaptertest",
		DesDoc: "queries",
	}

	store := NewCouchStore(res, newTestObj)
	query := NewObjQuery(newTestObj, store, res)
	err, objs := store.Read(query)

	if err != nil {
		t.Error("Obj Query failed", err)
	} else {
		for _, obj := range objs {
			t.Log(obj.GetKey())
		}
	}
	// Obj query type.
}
