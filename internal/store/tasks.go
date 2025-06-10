package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	appErrors "github.com/travboz/backend-projects/todo-list-api/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TasksStore struct {
	db    *mongo.Collection
	cache *redis.Client
}

// NewTasksStore creates a new TasksStore that implements TasksModel interface
func NewTasksStore(collection *mongo.Collection, cache *redis.Client) TasksModel {
	return &TasksStore{
		db:    collection,
		cache: cache,
	}
}

func (ts *TasksStore) invalidateCache(ctx context.Context, id string) {
	cacheKey := fmt.Sprintf("%s:%s", TasksCacheKeyBase, id)
	ts.cache.Del(ctx, cacheKey)
}

func (ts *TasksStore) Insert(ctx context.Context, task *data.Task) error {
	task.ID = primitive.NewObjectID()

	_, err := ts.db.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TasksStore) GetTaskById(ctx context.Context, id string) (*data.Task, error) {
	cacheKey := fmt.Sprintf("%s:%s", TasksCacheKeyBase, id)

	// 1. Try cache first
	if task, ok := ts.getTaskFromCache(ctx, cacheKey); ok {
		return task, nil
	}

	// 2. Cache miss or error - try database
	task, err := ts.getTaskFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Found in database - cache it for next time
	ts.cacheTask(ctx, cacheKey, task)

	return task, nil
}

func (ts *TasksStore) getTaskFromCache(ctx context.Context, key string) (*data.Task, bool) {
	val, err := ts.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, false // cache miss or error
	}

	var task data.Task
	if err := json.Unmarshal([]byte(val), &task); err != nil {
		return nil, false // failed to parse cached value
	}

	return &task, true
}

func (ts *TasksStore) getTaskFromDB(ctx context.Context, id string) (*data.Task, error) {
	taskID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := ts.db.FindOne(ctx, bson.M{"_id": taskID})

	var task data.Task
	if err = result.Decode(&task); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrRecordNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (ts *TasksStore) cacheTask(ctx context.Context, key string, task *data.Task) {
	if data, err := json.Marshal(task); err == nil {
		ts.cache.Set(ctx, key, data, CacheExpiryTime)
	}
}

func (ts *TasksStore) FetchAllTasks(ctx context.Context, p data.Filters, search string) ([]*data.Task, data.Metadata, error) {
	filter := bson.D{}
	if search != "" {
		filter = bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: search}}}}
	}

	if strings.Contains(search, "page_size=") {
		return nil, data.Metadata{}, appErrors.ErrInvalidQueryString
	}

	limit := int64(p.Limit())
	skip := int64(p.Offset())

	findOptions := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}

	// Count total matching documents
	totalRecords, err := ts.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, data.Metadata{}, err
	}

	cursor, err := ts.db.Find(ctx, filter, &findOptions)
	if err != nil {
		return nil, data.Metadata{}, err
	}

	defer cursor.Close(ctx)

	var tasks []*data.Task

	for cursor.Next(ctx) {
		var task *data.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, data.Metadata{}, err
		}

		tasks = append(tasks, task)
	}

	metadata := data.CalculateMetadata(int(totalRecords), p.Page, p.Limit())

	return tasks, metadata, nil

}
func (ts *TasksStore) UpdateTask(ctx context.Context, id string, task *data.Task) (*data.Task, error) {
	task_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // invalid id
	}

	filter := bson.M{"_id": task_id}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: task.Title},
			{Key: "description", Value: task.Description},
			{Key: "completed", Value: task.Completed},
		}},
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := ts.db.FindOneAndUpdate(
		ctx, filter, update, &opt,
	)

	var updatedTask data.Task

	if err = result.Decode(&updatedTask); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, appErrors.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	ts.invalidateCache(ctx, id)

	return &updatedTask, nil

}
func (ts *TasksStore) DeleteTask(ctx context.Context, id string) error {
	task_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := ts.db.DeleteOne(ctx, bson.M{"_id": task_id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return appErrors.ErrRecordNotFound
	}

	ts.invalidateCache(ctx, id)

	return nil
}

func (ts *TasksStore) CompleteTask(ctx context.Context, id string) (*data.Task, error) {
	taskID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // invalid id
	}

	filter := bson.M{"_id": taskID}
	update := bson.M{
		"$set": bson.M{
			"completed": true,
		},
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := ts.db.FindOneAndUpdate(ctx, filter, update, &opt)

	var completedTask data.Task

	if err = result.Decode(&completedTask); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrRecordNotFound
		}
		return nil, err
	}

	ts.invalidateCache(ctx, id)

	return &completedTask, nil
}

func (ts *TasksStore) GetTaskOwnerId(ctx context.Context, id string) (string, error) {
	task_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	result := ts.db.FindOne(ctx, bson.M{"_id": task_id})

	var task data.Task
	if err = result.Decode(&task); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return "", appErrors.ErrRecordNotFound
		default:
			return "", err
		}
	}

	return task.Owner, nil
}
