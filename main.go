package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gospackler/caddyshack"
	"github.com/gospackler/caddyshack/model"
	"github.com/gospackler/caddyshack/resource"
	"github.com/gospackler/couchdb"
)

type CouchStore struct {
	model  *model.Definition
	client *couchdb.Client
	DbObj  *couchdb.Database
	DesDoc map[string]*couchdb.DesignDoc

	N   int
	Res *resource.Definition

	// Fields for caddyshack follows
	// This is needed for identifying adpter for caddyshack.
	StoreName string
	// For all queries, receiver for data.
	ObjType  reflect.Type
	DefQuery *CouchQuery
}

// FIXME : Assert Kind of the objModel passed is a pointer.
func NewCouchStore(res *resource.Definition, objModel caddyshack.StoreObject) (c *CouchStore) {
	objType := reflect.ValueOf(objModel).Elem().Type()
	client := couchdb.NewClient(res.Host, res.Port)
	c = &CouchStore{
		client:    &client,
		StoreName: "couchdb",
		ObjType:   objType,
		DesDoc:    make(map[string]*couchdb.DesignDoc),
		Res:       res,
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
		panic("Could not connect with db " + res.Host + string(res.Port) + err.Error())
	}
	c.DbObj = &dbObj
	c.DefQuery = NewObjQuery(objModel, c)
	return
}

// FIXME Reason out and remove this method in future.
func (c *CouchStore) Init(model *model.Definition) (error, caddyshack.Store) {
	c.model = model
	return nil, c
}

func (c *CouchStore) GetDesignDoc(docName string) *couchdb.DesignDoc {
	_, exists := c.DesDoc[docName] //Checking if the view exists.

	// FIXME check in the db as well to make sure the document does not exist there.

	if exists == true {

		return c.DesDoc[docName]
	} else {

		err, doc := couchdb.RetreiveDocFromDb(docName, c.DbObj)
		log.Debug("Checking if document with name ", docName, " is present.")
		if err != nil {
			c.DesDoc[docName] = couchdb.NewDesignDoc(docName, c.DbObj)
		} else {
			c.DesDoc[docName] = doc
		}
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

// {"id":"bb0c3212953bdc7fad2ade66160c244d","key":"bb0c3212953bdc7fad2ade66160c244d","value":{"name":"abcd","surprise":"1234","field1":"field1","age":0,"id":""}}
// An example of json that could come up.
// Decodes a key : value type to a the registered object and returns it.
func (c *CouchStore) GetStoreObj(jsonObj []byte) (error, caddyshack.StoreObject) {
	jsonStream := bytes.NewBuffer(jsonObj)
	jsonDecoder := json.NewDecoder(jsonStream)

	type Object struct {
		Id    string          `json:"id"`
		Key   interface{}     `json:"key"`
		Value json.RawMessage `json:"value"`
	}

	tempObj := new(Object)
	err := jsonDecoder.Decode(tempObj)
	if err != nil {
		return err, nil
	}

	dynmaicObj := reflect.New(c.ObjType).Interface()
	err = json.Unmarshal(tempObj.Value, dynmaicObj)
	if err != nil {
		return err, nil
	}

	obj := dynmaicObj.(caddyshack.StoreObject)
	obj.SetKey(tempObj.Id)
	return nil, obj
}

func (c *CouchStore) MarshalStoreObjects(data []byte) (result []caddyshack.StoreObject, err error) {
	jsonStream := strings.NewReader(string(data))
	jsonDecoder := json.NewDecoder(jsonStream)

	type ObjInfo struct {
		NumRows int               `json:"total_rows"`
		Offset  int               `json:"offset"`
		Array   []json.RawMessage `json:"rows"`
	}

	objInfo := new(ObjInfo)

	err = jsonDecoder.Decode(objInfo)

	for _, row := range objInfo.Array {
		// Does the reflection part
		err, storeObj := c.GetStoreObj(row)
		if err != nil {
			err = errors.New("Marshal Object" + err.Error())
			break
		}
		result = append(result, storeObj)
	}

	return
}

func (c *CouchStore) ReadFromObj(obj caddyshack.StoreObject) ([]caddyshack.StoreObject, error) {
	key := obj.GetKey()
	return c.ReadFromView(c.Res.DesDoc, c.DefQuery.ViewName, key)
}

func (c *CouchStore) ReadByKey(key string) ([]caddyshack.StoreObject, error) {
	return c.ReadFromView(c.Res.DesDoc, c.DefQuery.ViewName, key)
}

// ReadOne : This works only with the default id that couchdb generates.
func (c *CouchStore) ReadOne(key string) (error, caddyshack.StoreObject) {
	doc := couchdb.NewDocument(key, "", c.DbObj)
	jsonObj, err := doc.GetDocument()
	if err != nil {
		return err, nil
	}

	dynmaicObj := reflect.New(c.ObjType).Interface()
	err = json.Unmarshal(jsonObj, dynmaicObj)
	if err != nil {
		return err, nil
	}

	obj := dynmaicObj.(caddyshack.StoreObject)
	obj.SetKey(doc.Id)
	return err, obj
}

func (c *CouchStore) DeleteOne(obj caddyshack.StoreObject) error {
	doc := couchdb.NewDocument(obj.GetKey(), "", c.DbObj)
	_, err := doc.GetDocument()
	if err != nil {
		return err
	}
	return doc.Delete()
}

func (c *CouchStore) DestroyOne(key string) error {
	doc := couchdb.NewDocument(key, "", c.DbObj)
	_, err := doc.GetDocument()
	if err != nil {
		return err
	}
	return doc.Delete()
}

func (c *CouchStore) ReadFromView(desDocName string, viewName string, key string) ([]caddyshack.StoreObject, error) {
	if !strings.Contains(desDocName, "/") {
		desDocName = "_design/" + desDocName
	}
	data, err := c.DbObj.GetView(desDocName, viewName, "key=\""+key+"\"")

	if err != nil {
		newErr := fmt.Errorf("Error retreiving : Key = %s ViewName = %s desDoc = %s :  %s", key, viewName, desDocName, err.Error())
		return nil, newErr
	} else {
		result, err := c.MarshalStoreObjects(data)
		if err != nil {
			return nil, errors.New("Could not Marshal json" + err.Error())
		}
		return result, nil
	}
}

func (c *CouchStore) ReadOneFromView(desDocName string, viewName string, key string) (caddyshack.StoreObject, error) {
	results, err := c.ReadFromView(desDocName, viewName, key)
	if err != nil {
		return nil, err
	}

	return results[0], nil
}

func (c *CouchStore) Read(query caddyshack.Query) (error, []caddyshack.StoreObject) {
	err, objects := query.Execute()
	// Use the rawJson to check for the view.

	return err, objects
}

// Read Default uses the default query object to make the request.
func (c *CouchStore) ReadDef() (err error, objects []caddyshack.StoreObject) {
	err, objects = c.Read(c.DefQuery)
	if err != nil {
		err = errors.New("Read Def :" + err.Error())
	}
	return
}

func (c *CouchStore) ReadN(query *CouchQuery) (objs []caddyshack.StoreObject, err error) {
	if query.BufferSize == 0 {
		return nil, errors.New("BufferSize is 0 in query for readN which is invalid")
	}

	if query.Skip == 0 && query.Limit == 0 {
		query.Skip = 0
		query.Limit = query.BufferSize
	} else {
		query.Skip = query.Skip + query.BufferSize
	}

	err, objs = c.Read(query)

	if len(objs) == 0 {
		// FIXME : Possible lose of error information.
		err = io.EOF
		query.Skip = 0
		query.Limit = 0
	}
	log.Debug("Query =", query)
	log.Debug("Objects count read =", len(objs))
	log.Debug("Error = ", err)
	return objs, err
}

// The object passed should have CouchWrapperUpdate as an anonymous field containing the details.
func (c *CouchStore) UpdateOne(obj caddyshack.StoreObject) (err error) {
	byteObj, err := json.Marshal(obj)
	doc := couchdb.NewDocument(obj.GetKey(), "", c.DbObj)
	err = doc.Update(byteObj)
	return
}
