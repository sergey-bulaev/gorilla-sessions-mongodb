package mongodbstoregorilla

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStore(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// TODO flag
	mongoURI := "mongodb://localhost:27017"
	mongoDB := "test"
	mongoColl := "mongodbstore_sessions_test"
	dropColl := true

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Error connecting to mongoDB: %v", err)
	}

	defer client.Disconnect(ctx)

	coll := client.Database(mongoDB).Collection(mongoColl)

	if dropColl {
		defer coll.Drop(ctx)
	}

	store, err := NewMongoDBStore(coll, []byte("secret"))
	if err != nil {
		t.Fatalf("Error initializing mongodb store: %v", err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	resp := httptest.NewRecorder()
	// Get a session.
	session, err := store.Get(req, "session-key")

	if !session.IsNew {
		t.Fatalf("Error getting session should be new: %v", err)
	}

	if err != nil {
		t.Fatalf("Error getting session: %v", err)
	}
	session.Values["kek"] = "lol"
	flashes := session.Flashes()
	if len(flashes) != 0 {
		t.Errorf("Expected empty flashes; Got %v", flashes)
	}
	// Add some flashes.
	session.AddFlash("foo")
	session.AddFlash("bar")
	// Custom key.
	// session.AddFlash("baz", "custom_key")
	// Save.
	if err = sessions.Save(req, resp); err != nil {
		t.Fatalf("Error saving session: %v", err)
	}
	header := resp.Header()
	cookies, ok := header["Set-Cookie"]
	if !ok || len(cookies) != 1 {
		t.Fatalf("No cookies. Header: %+v", header)
	}
	value, ok := session.Values["kek"]
	if !ok {
		t.Errorf("expected to have kek in session %t", value)
	}

	req, _ = http.NewRequest("GET", "http://localhost:8080/", nil)
	req.Header.Add("Cookie", cookies[0])
	resp = httptest.NewRecorder()
	// Get a session.
	if session, err = store.Get(req, "session-key"); err != nil {
		t.Fatalf("Error getting session: %v", err)
	}

	if session.IsNew {
		t.Fatalf("Error getting session should not be new: %v", err)
	}

	flashes = session.Flashes()
	if len(flashes) == 0 {
		t.Errorf("Expected 2 flashes; Got %v", flashes)
	}
}
