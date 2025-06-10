package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/travboz/backend-projects/todo-list-api/internal/data"
	"github.com/travboz/backend-projects/todo-list-api/internal/store"
)

type CachedTasksModel interface {
	GetTaskById(context.Context, string) (*data.Task, error)
	SetTask(context.Context, *data.Task) error
	GetTaskOwnerId(context.Context, string) (string, error)
}

type TasksCacheRedis struct {
	DB    store.TasksModel
	Redis *redis.Client
}

func NewTasksCacheRedis(db store.TasksModel, client *redis.Client) CachedTasksModel {
	return &TasksCacheRedis{
		DB:    db,
		Redis: client,
	}
}

func (c *TasksCacheRedis) GetTaskById(ctx context.Context, id string) (*data.Task, error) {
	var task *data.Task

	cacheKey := fmt.Sprintf("task:%s", id)

	err := c.Redis.HGetAll(ctx, cacheKey).Scan(&task)
	if err == nil {
		return task, nil
	}

	task, err = c.DB.GetTaskById(ctx, id)
	if err != nil {
		return nil, err
	}

	// we've grabbed from DB, so set in cache
	_, err = c.Redis.HSet(ctx, cacheKey, map[string]interface{}{
		"id":          task.ID,
		"owner":       task.Owner,
		"title":       task.Title,
		"description": task.Description,
		"completed":   task.Completed,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
	}).Result()

	if err != nil {
		return nil, err
	}

	c.Redis.Expire(ctx, cacheKey, 5*time.Minute)

	return task, nil
}

func (c *TasksCacheRedis) SetTask(ctx context.Context, task *data.Task) error {
	return nil
}

func (c *TasksCacheRedis) GetTaskOwnerId(ctx context.Context, id string) (string, error) {
	return "", nil
}
