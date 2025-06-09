package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/travboz/backend-projects/todo-list-api/internal/store"
	"github.com/travboz/backend-projects/todo-list-api/internal/store/mocks"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Storage: &store.Storage{
			Users: &mocks.UsersStoreMock{},
			Tasks: &mocks.TasksStoreMock{},
		},
	}
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

// Implement a get() method on our custom testServer type. This makes a GET
// request to a given url path using the test server client, and returns the
// response status code, headers and body.
func (ts *testServer) makeGetRequest(t *testing.T, urlPath, token string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) makeGetRequestWithToken(t *testing.T, urlPath, token string) (int, http.Header, string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

// // InsertTestData inserts a dummy user, task, and token for testing
// func InsertTestData(t *testing.T, ctx context.Context, usersColl, tasksColl, tokensColl *mongo.Collection) (user data.User, task data.Task, tokenValue string) {
// 	t.Helper()

// 	// Insert user
// 	user = data.User{
// 		ID:       primitive.NewObjectID(),
// 		Name:     "Test User",
// 		Email:    "test@example.com",
// 		Password: "hashed-password", // hash not necessary for testing
// 	}
// 	_, err := usersColl.InsertOne(ctx, user)
// 	if err != nil {
// 		t.Fatalf("failed to insert user: %v", err)
// 	}

// 	// Insert task
// 	task = data.Task{
// 		ID:          primitive.NewObjectID(),
// 		Owner:       user.ID.Hex(), // store as string
// 		Title:       "Green Mocks & Ham",
// 		Description: "Test task",
// 		Completed:   false,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// 	_, err = tasksColl.InsertOne(ctx, task)
// 	if err != nil {
// 		t.Fatalf("failed to insert task: %v", err)
// 	}

// 	// Insert token
// 	tokenValue = "testtoken"
// 	token := data.Token{
// 		ID:        primitive.NewObjectID(),
// 		UserID:    user.ID,
// 		Token:     tokenValue,
// 		CreatedAt: time.Now(),
// 	}
// 	_, err = tokensColl.InsertOne(ctx, token)
// 	if err != nil {
// 		t.Fatalf("failed to insert token: %v", err)
// 	}

// 	return user, task, tokenValue
// }

// // CleanupTestData removes all test documents related to the given user ID
// func CleanupTestData(t *testing.T, ctx context.Context, userID string, usersColl, tasksColl, tokensColl *mongo.Collection) {
// 	t.Helper()

// 	// Remove the user
// 	_, err := usersColl.DeleteOne(ctx, bson.M{"_id": userID})
// 	if err != nil {
// 		t.Fatalf("failed to delete test user: %v", err)
// 	}

// 	// Remove any tasks owned by the user
// 	_, err = tasksColl.DeleteMany(ctx, bson.M{"owner": userID})
// 	if err != nil {
// 		t.Fatalf("failed to delete test tasks: %v", err)
// 	}

// 	// Remove tokens linked to the user
// 	objectID, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		t.Fatalf("invalid userID: %v", err)
// 	}

// 	_, err = tokensColl.DeleteMany(ctx, bson.M{"user_id": objectID})
// 	if err != nil {
// 		t.Fatalf("failed to delete test tokens: %v", err)
// 	}
// }
