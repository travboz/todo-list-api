package mongo

import (
	"context"
	"errors"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	appErrors "github.com/travboz/backend-projects/todo-list-api/internal/errors"
	"github.com/travboz/backend-projects/todo-list-api/internal/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStoreTasks struct {
	Tasks *mongo.Collection
}

func (ms *MongoStorage) NewMongoTasksModel() store.TasksModel {
	return MongoDBStoreTasks{ms.DB.Collection("tasks")}
}

func (t MongoDBStoreTasks) Insert(ctx context.Context, task *data.Task) error {
	task.ID = primitive.NewObjectID()

	_, err := t.Tasks.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	return nil
}
func (t MongoDBStoreTasks) GetTaskById(ctx context.Context, id string) (*data.Task, error) {
	task_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := t.Tasks.FindOne(ctx, bson.M{"_id": task_id})

	var task data.Task
	if err = result.Decode(&task); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, appErrors.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &task, nil

}

func (t MongoDBStoreTasks) FetchAllTasks(ctx context.Context, p data.Filters, search string) ([]*data.Task, data.Metadata, error) {
	filter := bson.D{{"$text", bson.D{{"$search", search}}}}

	limit := int64(p.Limit())
	skip := int64(p.Offset())

	findOptions := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}

	cursor, err := t.Tasks.Find(ctx, filter, &findOptions)
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

	metadata := data.CalculateMetadata(len(tasks), p.Page, p.Limit())

	return tasks, metadata, nil

}
func (t MongoDBStoreTasks) UpdateTask(ctx context.Context, id string, task *data.Task) (*data.Task, error) {
	task_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // invalid id
	}

	filter := bson.M{"_id": task_id}
	update := bson.D{
		{"$set", bson.D{
			{"title", task.Title},
			{"description", task.Description},
			{"completed", task.Completed},
		}},
	}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := t.Tasks.FindOneAndUpdate(
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

	return &updatedTask, nil

}
func (t MongoDBStoreTasks) DeleteTask(ctx context.Context, id string) error {
	task_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := t.Tasks.DeleteOne(ctx, bson.M{"_id": task_id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return appErrors.ErrRecordNotFound
	}

	return nil
}

func (t MongoDBStoreTasks) CompleteTask(ctx context.Context, id string) (*data.Task, error) {
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

	result := t.Tasks.FindOneAndUpdate(ctx, filter, update, &opt)

	var completedTask data.Task

	if err = result.Decode(&completedTask); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrRecordNotFound
		}
		return nil, err
	}

	return &completedTask, nil
}

func (t MongoDBStoreTasks) GetTaskOwnerId(ctx context.Context, id string) (string, error) {
	task_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	result := t.Tasks.FindOne(ctx, bson.M{"_id": task_id})

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
