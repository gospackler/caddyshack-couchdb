package adapter

import (
	"github.com/georgethomas111/caddyshack"
	"github.com/georgethomas111/caddyshack/model"
	"github.com/georgethomas111/caddyshack/resource"
	"github.com/georgethomas111/couchdb"
	"testing"
)

// Create a compatable storeObject
type TestObj struct {
	Name  string
	Value string
	Id    string
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
}

func TestCreate(t *testing.T) {

	testObj := &TestObj{
		Name:  "abcd",
		Value: "1234",
	}

	err := Caddy.StoreIns.Create(testObj)
	if err != nil {

		t.Error("Error creating object in the test Store")
	}

	Key = testObj.GetKey()
	t.Log("Created Object with key", testObj.GetKey())
}

var RetrObject *RetTestObj

func TestRead(t *testing.T) {

	err, obj := Caddy.StoreIns.ReadOne(Key)
	if err != nil {
		t.Error("Error while retreiving object")
	}

	if obj.GetKey() != Key {
		t.Error("Retreived wrong object")
	}

	actualObj := obj.(*RetTestObj)
	t.Log("Got the actula object back.", actualObj.TestObj)
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
