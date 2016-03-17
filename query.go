package adapter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/bushwood/caddyshack"
	"github.com/bushwood/couchdb"
)

// View Object is placed over here as a query in couch is possible only with a view
// Helps in forming the javascript that can be used for working on stuff.
type ViewObj struct {
	Name     string
	ViewType reflect.Type
}

// FIXME : Assert Kind of the viewObj passed is a pointer.
func NewViewObj(name string, viewObj caddyshack.StoreObject) *ViewObj {
	// Check if View Exists in the DB.
	// Create View Using the tags if thats the case.
	return &ViewObj{
		Name:     name,
		ViewType: reflect.ValueOf(viewObj).Elem().Type(),
	}
}

// Create the View if it does not exist.
func (v *ViewObj) GetCondition() string {
	// Use the tags to form the condition.
	return "Not implemented"
}

// Initial Version
type CouchQuery struct {
	Condition string // Code for the view RawJson
	ViewName  string
	desDoc    *couchdb.DesignDoc
}

func NewQuery(line string, viewName string, desDoc string, db *CouchStore) (couchQuery *CouchQuery) {

	// Assuming a design doc is already created.
	desDocObj := db.GetDesignDoc(desDoc)

	couchQuery = &CouchQuery{
		desDoc:    desDocObj,
		Condition: line,
		ViewName:  viewName,
	}
	return
}

// Use reflection to create the query from the tags.
func NewObjQuery(obj caddyshack.StoreObject, viewName string, desDoc string, db *CouchStore) (query *CouchQuery) {

	viewObj := NewViewObj(viewName, obj)
	desDocObj := db.GetDesignDoc(desDoc)

	query = &CouchQuery{
		desDoc:    desDocObj,
		Condition: viewObj.GetCondition(),
	}
	return
}

func (q *CouchQuery) SetCondition(cond string) {
	q.Condition = cond
}

func (q *CouchQuery) GetCondition() string {
	return q.Condition
}

func (q *CouchQuery) Execute() (err error, result []caddyshack.StoreObject) {
	// Currently O(n) w.r.t to views
	status := q.desDoc.CheckExists(q.ViewName)

	if status == true {
		err, data := q.desDoc.Db.GetView(q.desDoc.Id, q.ViewName)
		if err != nil {
			return errors.New("Error retreiving view : " + err.Error()), nil
		} else {
			// Print for now create store Object later.
			// FIXME Handle unmarshalling over here.
			fmt.Println(string(data))
		}
	} else {
		// The intutive version would be creating an object and then adding methods to it.
		newView := &couchdb.View{Name: q.ViewName}
		newView.RawStatus = true
		newView.RawJson = q.Condition
		q.desDoc.AddView(newView)
		err = q.desDoc.SaveDoc()
		// possible infinite recursion. Should be fun :D
		err, result = q.Execute()
	}

	return
	// If it exists get the view back.
	// Otherwise Get Retrieve the Data and Marshal the store Object from the json..
}
