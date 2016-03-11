## Caddyshack-couchdb

This is the couchdb adapter for caddyshack. The skeleton is in place based on the caddyshack right now. 

### How it works?
It is basically a type which implements the Store interface of caddyshack. Within each of the functions the couchdb adapter is used to get the work done.

For Reads the reflection 

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

The adapter requires a type to be registered with it while initializing. This type will be created while retrieving the object. *(json.Unmarshal)* uses it and dynamically created using reflection.

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
actualObj := obj.(*RetTestObj)  // Explicit assertion is needed to retrieve the object.
fmt.Println("Got the actual object back.", actualObj.TestObj)
```
Explicit assertion is needed because go being a statically compiled language requires the type to be bound to it. The interface object is received back which needs to be converted to the object type we need.

For more details, have a look at *couch_test.go* which contains the case of *registering, creating an object and retrieving* the object from a couch db.
