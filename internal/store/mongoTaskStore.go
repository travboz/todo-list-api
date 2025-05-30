package store

import (
	"context"
	"errors"

	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStoreTasks struct {
	tasks *mongo.Collection
}

func (t MongoDBStoreTasks) Insert(ctx context.Context, task *data.Task) error {
	task.ID = primitive.NewObjectID()

	_, err := t.tasks.InsertOne(ctx, task)
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

	result := t.tasks.FindOne(ctx, bson.M{"_id": task_id})

	var task data.Task
	if err = result.Decode(&task); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &task, nil

}

// func (t MongoDBStoreTasks) GetOwnerOfTask(context.Context, string) (string, error) {
// 	return "", nil
// }

func (t MongoDBStoreTasks) FetchAllTasks(ctx context.Context) ([]*data.Task, error) {
	filter := bson.M{}

	cursor, err := t.tasks.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var tasks []*data.Task

	for cursor.Next(ctx) {
		var task *data.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil

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

	result := t.tasks.FindOneAndUpdate(
		ctx, filter, update, &opt,
	)

	var updatedTask data.Task

	if err = result.Decode(&updatedTask); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
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

	result, err := t.tasks.DeleteOne(ctx, bson.M{"_id": task_id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrRecordNotFound
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

	result := t.tasks.FindOneAndUpdate(ctx, filter, update, &opt)

	var completedTask data.Task

	if err = result.Decode(&completedTask); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &completedTask, nil
}
