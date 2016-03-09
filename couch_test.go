package adapter

import (
	"github.com/georgethomas111/caddyshack"
	"github.com/georgethomas111/caddyshack/model"
	"github.com/georgethomas111/caddyshack/resource"
	"testing"
)

// Create a compatable storeObject
type TestObj struct {
	Name  string
	Value string
	Id    string
}

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
	couchStore := NewCouchStore(res)
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

	testObj := &TestObj{
		Name:  "abcd",
		Value: "1234",
	}
	caddyName := model.Name + model.Adapter
	err, caddy := cs.GetCaddy(caddyName)
	if err != nil {
		t.Error("Error while retreiving caddy ", err)
	}

	err = caddy.StoreIns.Create(testObj)
	if err != nil {

		t.Error("Error creating object in the test Store")
	}

	t.Log("Created Object with key", testObj.GetKey())

	err, obj := caddy.StoreIns.ReadOne(testObj.GetKey())
	if err != nil {
		t.Error("Error while retreiving object")
	}

	t.Log("Retreived Object", obj)

	if obj.GetKey() != testObj.GetKey() {
		t.Error("Retreived wrong object")
	}
}
