package adapters

import (
    // "time"

    // "github.com/bushwood/caddyshack/adapter"
    "github.com/bushwood/caddyshack/resource"
    // "github.com/satori/go.uuid"
    // couch "github.com/rhinoman/couchdb-go"
)

const DBNAME string = "_CADDYSHACK"

var Adapter = &Definition{
    Name: "couchdb",
}

type Definition struct {
    Name string
    Config resource.Definition
}

// GetConfig returns the config resource of the adapter
func (adp *Definition) GetConfig() (rsc resource.Definition) {
    return adp.Config
}

// SetConfig sets the config resource of the adapter
func (adp *Definition) SetConfig(rsc resource.Definition) (error) {
    adp.Config = rsc
    return nil
}

// GetConfig returns the name of the adapter
func (adp *Definition) GetName() (string) {
    return adp.Name
}

// SetConfig sets the name resource of the adapter
func (adp *Definition) SetName(name string) (error) {
    adp.Name = name
    return nil
}

// Connect connects the adapter to the database
func (adp *Definition) Open() (error) {
    return nil
}

// Close closes the connection to teh database
func (adp *Definition) Close() (error) {
    return nil
}

//
// var CouchDB caddyshack.Adapter = caddyshack.Adapter{
//     Connect: func(config caddyshack.Resource) (caddyshack.ORM, error) {
//         h := config.Host
//         p := atoi(config.Port)
//         t := time.Duration(config.Timeout)
//         uname := config.Username
//         pass := config.Password
//         conn, err := couch.NewConnection(h, p, t)
//         auth := couch.BasicAuth{Username: uname, Password: pass}
//         return NewCouchORM(conn.SelectDB(DBNAME, &auth))
//     },
// }
//
// func NewCouchORM(db couch.Database) (caddyshack.ORM, error) {
//     return caddyshack.ORM{
//         Create: func(doc interface{}) (interface{}, error) {
//             id := uuid.NewV4().String()
//             db.Save(doc, id, "")
//             return interface{}, nil
//         },
//         Find: func(query interface{}) (interface{}, error) {
//             db.Read()
//             return interface{}, nil
//         },
//         FindOne: func(query interface{}) (interface{}, error) {
//             db.Read()
//             return interface{}, nil
//         },
//         Update: func(query interface{}, doc interface{}) (interface{}, error){
//             db.Save()
//             return interface{}, nil
//         },
//         Destroy: func(query interface{}) (interface{}, error) {
//             db.Delete()
//             return interface{}, nil
//         },
//         Connect func() (error) {
//             return nil
//         }
//         Close: func() (error) {
//             return nil
//         }
//     }
// }
