## Caddyshack-couchdb

This is the couchdb adapter for caddyshack. The caddyshack code and idea can be found here. 

https://github.com/bushwood/caddyshack

### How it works?

It is basically a type which implements the Store interface of caddyshack. Within each of the functions the couchdb adapter is used to get the work done.

Reading elements works. Only the tags need to be played with. 
```go
type TestObj struct {
        Name  string
        Value string
        Id    string
}
```

```go
type RetTestObj struct {
        couchdb.CouchWrapperUpdate
        TestObj
}
```
The retrieve and save structs need to be separate as fields get added like _id and _rev while receiving.

The adapter requires a type to be registered with it while initializing. This type will be created while retrieving the object. *(json.Unmarshal)* uses it and is dynamically created using reflection.

```go
couchStore := NewCouchStore(res, &RetTestObj{})
```

Some working functions
For creating an object of type testObj in couch. tags could be used to get the right names.

```go
err = caddy.StoreIns.Create(testObj)
```

To get an object when a key is passed,
```go
err, obj := caddy.StoreIns.ReadOne(key)
actualObj := obj.(*RetTestObj)  // Assertion needed as golang is compiled
fmt.Println("Got the actual object back.", actualObj.TestObj)
```
Explicit assertion is needed because go being a statically compiled language requires the type to be bound to it. The interface object is received back which needs to be converted to the object type we need.

For more details, have a look at *couch_test.go* which contains the case of *registering, creating an object and retrieving* the object from a couch db.

## Queries

Couch Wrapper supports queries the way caddyshack expects it to. There are two types of Queries possible. 
### 1. Basic Couch Query
In this type of query raw javascript can be passed to populate the fields that the json unmarshal can use. 

``` go
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
````
### 2. Object Query

In this type of query the Objects tags are read to create the query dynamically. ie it creates permanent views in the database based on the tags.
    
    ``` go
        type TestObj struct {
            Name   string `json:"name"`
            Value  string `json:"surprise"`
            Field1 string `json:"field1"`
            Age    int    `json:"age"`
            Id     string `json:"id"`
        }
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
                        testObj := obj.(*TestObj)
                        t.Log(testObj)
                }
        }
    ```

The view is created dynamically based on the tags and the registered object in this case *TestObj* is used to create the view. The name of the view will be *testobj* and would be saved in the design document named *_design/queries*

### 3. Object Query With Conditions

``` go
type TestObjCond struct {
        Name   string `json:"name"`
        Value  string `json:"surprise"`
        Field1 string `json:"field1"`
        Age    int    `json:"age" condition:"age < 20"`
        Id     string `json:"id"`
}
```

This is similar to Object Query by allows conditions before emitting fields. All the conditions will have an *AND* operation between them to facilitate the easy retreival of objects.
