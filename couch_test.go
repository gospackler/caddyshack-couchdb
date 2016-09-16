package adapter

import (
	"github.com/gospackler/caddyshack"
	"github.com/gospackler/caddyshack/model"
	"github.com/gospackler/caddyshack/resource"
	"github.com/gospackler/couchdb"
	"testing"
)

// Create a compatable storeObject
type TestObj struct {
	Name   string `json:"name"`
	Value  string `json:"surprise"`
	Field1 string `json:"field1"`
	Age    int    `json:"age"`
	Id     string `json:"id"`
}

type TestObjCond struct {
	Name   string `json:"name" condition:"name == 'Updated'"`
	Value  string `json:"surprise" condition:"surprise == -1"`
	Field1 string `json:"field1"`
	Age    int    `json:"age"`
	Id     string `json:"id"`
}

func (t *TestObjCond) GetKey() string {
	return t.Id
}

func (t *TestObjCond) SetKey(id string) {
	t.Id = id
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

func (t *TestObj) String() string {
	return "\n Name " + t.Name +
		"\n Value " + t.Value +
		"\n Field1 " + t.Field1
}

var CouchStoreIns *CouchStore

func getCouchStore(t *testing.T) *CouchStore {
	// Add model definition in future to it.
	cs := caddyshack.New()

	res := &resource.Definition{
		Host:   "127.0.0.1",
		Port:   5984,
		Name:   "adaptertest",
		DesDoc: "queries",
	}

	// From storedemo.go
	couchStore := NewCouchStore(res, &TestObj{})
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
	return couchStore
}

func getStoreCond(t *testing.T) *CouchStore {
	newTestObj := new(TestObjCond)
	res := &resource.Definition{
		Host:   "127.0.0.1",
		Port:   5984,
		Name:   "adaptertest",
		DesDoc: "queries",
	}

	store := NewCouchStore(res, newTestObj)
	return store
}
func TestInit(t *testing.T) {
	getCouchStore(t)
}

func TestCreate(t *testing.T) {
	testObj := &TestObj{
		Name:   "abcd",
		Value:  "1234",
		Field1: "field1",
		Age:    10,
	}

	err := Caddy.StoreIns.Create(testObj)
	if err != nil {

		t.Error("Error creating object in the test Store")
	}

	Key = testObj.GetKey()
	if Key == "" {
		t.Error("Error while creating, Key not obtained.")
	}
	t.Log("Object Created Key :=", testObj.GetKey())
}

func TestReadOneAndUpdate(t *testing.T) {
	err, obj := Caddy.StoreIns.ReadOne(Key)
	if err != nil {
		t.Error("Error while retreiving object")
	}

	if obj.GetKey() != Key {
		t.Error("Retreived wrong object")
	}

	actualObj := obj.(*TestObj)
	t.Log("Got the actual object back.", actualObj)

	actualObj.Name = "Updated"
	actualObj.Value = "-1"
	err = Caddy.StoreIns.UpdateOne(actualObj)
	if err != nil {
		t.Log("Error while updating object, ", err)
	}
	t.Log("Check the updated object in the DB if delete is disabled with key", Key)
}

func TestRead(t *testing.T) {
	// Every Query is the request to a view.
	query := NewQuery("function(doc) {emit(doc.field1);}", "new_view", "new_design", getCouchStore(t))
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

func TestReadCustomMR(t *testing.T) {
	mapR := "\"map\": \"function(doc) {emit(doc.field1);}\", \"reduce\" : \"function(keys, values) {console.log('reduce');}\""
	query := NewMRQuery(mapR, "new_view_reduce", "new_design", getCouchStore(t))
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
	store := getCouchStore(t)
	query := NewObjQuery(new(TestObj), store)
	err, objs := store.Read(query)

	if err != nil {
		t.Error("Obj Query failed", err)
	} else {
		for _, obj := range objs {
			testObj := obj.(*TestObj)
			t.Log(testObj)
		}
	}
	// Obj query type.
}

func TestObjQueryCondRead(t *testing.T) {
	store := getStoreCond(t)
	query := NewObjQuery(new(TestObjCond), store)
	err, objs := store.Read(query)

	if err != nil {
		t.Error("Obj Query Condition failed", err)
	} else {
		for _, obj := range objs {
			testObj := obj.(*TestObjCond)
			t.Log(testObj)
		}
	}
	// Obj query type.
}

func TestObjDoubleQuery(t *testing.T) {
	store := getStoreCond(t)
	query := NewObjQuery(new(TestObjCond), store)
	err, objs := store.Read(query)

	if err != nil {
		t.Error("Obj Query Condition failed", err)
	} else {
		for _, obj := range objs {
			testObj := obj.(*TestObjCond)
			t.Log(testObj)
		}
	}
	NewObjQuery(new(TestObjCond), store)
}

// Test for the new default method
func TestDef(t *testing.T) {
	store := getStoreCond(t)
	err, objs := store.ReadDef()

	if err != nil {
		t.Error("Read With Default query failed ", err)
	} else {
		for _, obj := range objs {
			testObj := obj.(*TestObjCond)
			t.Log(testObj)
		}
	}
}

func TestGetKeyFromView(t *testing.T) {
	store := getStoreCond(t)
	newTestObj := new(TestObjCond)
	query := NewObjQuery(newTestObj, store)
	viewName := query.GetViewName(newTestObj)
	desDocName := "queries"
	obj, err := store.ReadOneFromView(desDocName, viewName, Key)
	if err != nil {
		t.Error("Error reading one from view", err)
	}
	t.Log(obj)
}

func TestDelete(t *testing.T) {
	newTestObj := new(TestObjCond)
	store := getStoreCond(t)
	query := NewObjQuery(newTestObj, store)
	viewName := query.GetViewName(newTestObj)
	desDocName := "queries"
	obj, err := store.ReadOneFromView(desDocName, viewName, Key)
	if err != nil {
		t.Error("Error reading one from view", err)
	}
	err = store.DeleteOne(obj)
	if err != nil {
		t.Error("Error reading one from view", err)
	}
	t.Log("Successfully deleted object with name ", Key)
}

// FIXME : Add tests with two object Queries.

// FIXME : Add a way to retrieve an Object if it already exists.
// FIXME : Add update option
