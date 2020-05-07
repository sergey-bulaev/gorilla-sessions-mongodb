# gorilla-sessions-mongodb

[Gorilla's Session](http://www.gorillatoolkit.org/pkg/sessions) store implementation for mongoDB using [official Go driver](https://pkg.go.dev/go.mongodb.org/mongo-driver?tab=overview)


## Installation
    go get github.com/2-72/gorilla-sessions-mongodb

### Example
```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	mongodbstore "github.com/2-72/gorilla-sessions-mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf("error connecting to mongodb %v \n", err)
		return
	}
	coll := client.Database("app_db").Collection("sessions")
	store, err := mongodbstore.NewMongoDBStore(coll, []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		log.Printf("error initializing mongodb store %v \n", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, _ := store.Get(r, "session-name")
		// Set some session values.
		session.Values["foo"] = "bar"
		session.Values[42] = 43
		// Save it before we write to the response/return from the handler.
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

```