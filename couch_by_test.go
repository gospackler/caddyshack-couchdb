package adapter

import (
	"testing"

	"github.com/gospackler/caddyshack/resource"
)

type TestObjBy struct {
	Name   string `json:"name"`
	Value  string `json:"surprise"`
	Field1 string `json:"field1" by:"field1"`
	Age    int    `json:"age"`
	Id     string `json:"id" by:"id"`
}

func (t *TestObjBy) GetKey() string {
	return t.Field1 + "," + t.Id
}

func (t *TestObjBy) SetKey(id string) {
	t.Id = id
}

func getByCouchStore(t *testing.T) *CouchStore {
	res := &resource.Definition{
		Host:   "127.0.0.1",
		Port:   5984,
		Name:   "adaptertest",
		DesDoc: "queries",
	}

	couchStore := NewCouchStore(res, &TestObjBy{})
	return couchStore
}

func TestByCreate(t *testing.T) {
	testObj := &TestObjBy{
		Name:   "Updated",
		Value:  "-1",
		Field1: "field1",
		Age:    11,
		Id:     "nnnn",
	}

	store := getByCouchStore(t)
	err := store.Create(testObj)
	if err != nil {
		t.Error("Error creating object in the test Store")
	}
}

func TestByRead(t *testing.T) {
	testObjDum := &TestObjBy{
		Field1: "field1",
		Id:     "nnnn",
	}

	store := getByCouchStore(t)
	obj, err := store.ReadOneFromObj(testObjDum)

	if err != nil {
		t.Error(err)
	}

	actualObj := obj.(*TestObjBy)
	if actualObj.Age != 11 {
		t.Error("Could not retreive the object from couch")
	}
}
