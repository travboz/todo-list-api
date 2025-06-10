# Todo List API
![Superhero Gopher - Project Title Image](https://raw.githubusercontent.com/egonelbre/gophers/63b1f5a9f334f9e23735c6e09ac003479ffe5df5/vector/superhero/standing.svg)

## Description

This is a RESTful API designed to manage personal to-do lists with some basic secure user access. It supports full CRUD functionality for to-do items and includes user registration, login with token-based authentication, and protected routes. To-dos can be listed, filtered, and paginated, with robust validation, error handling, and security best practices.

The brief follows:
[Todo List API Project](https://roadmap.sh/projects/todo-list-api)

## Features

Features
- User registration to create a new account
- Login endpoint to authenticate users and issue tokens
- Token-based user authentication for protected routes
- Create, read, update, and delete to-do items
- List to-do items with pagination and filtering
- Validates input data for users and to-do entries
- Handles common errors with structured responses
- Stores user and to-do data in a persistent database
- Follows RESTful API design principles
- Implements basic security best practices
- Caching for frequently access data


## Getting Started

### Prerequisites
- Docker
- Docker Compose
- Go (1.22+ recommended)

## Installation

1. Clone this repository:
   ```sh
   git clone https://github.com/travboz/todo-list-api.git
   cd todo-list-api
   ```
2. Set up Go modules:
   ```sh
   go mod tidy
   ```   
3. Run docker containers for MongoDB and Redis instances:
    ```sh
    make compose/up
    ```
4. Seed DB instance:
   ```sh
   make seed
   ```
5. Run server:
    ```sh
    make run
    ```
6. Navigate to `http://localhost<SERVER_PORT>` and call an endpoint

I will use example port `":7666"`.

### `.env` file
This server uses a `.env` environment file for configuration.
For an example, see `.env.example`.

## API Endpoints

| Method   | Endpoint                     | Description           | Auth Required |
| -------- | ---------------------------- | --------------------- | ------------- |
| `GET`    | `/api/v1/healthcheck`        | Health check          | No            |
| `POST`   | `/api/v1/users/register`     | Register a new user   | No            |
| `POST`   | `/api/v1/users/login`        | User login            | No            |
| `GET`    | `/api/v1/users/:id`          | Get user by ID        | Yes           |
| `POST`   | `/api/v1/tasks/create`       | Create a new task     | Yes           |
| `GET`    | `/api/v1/tasks`              | Fetch all tasks       | Yes           |
| `GET`    | `/api/v1/tasks/:id`          | Get task by ID        | Yes           |
| `PUT`    | `/api/v1/tasks/:id/complete` | Mark task as complete | Yes           |
| `PATCH`  | `/api/v1/tasks/:id`          | Update task           | Yes           |
| `DELETE` | `/api/v1/tasks/:id`          | Delete task           | Yes           |


## Authentication

This API uses Bearer Token authentication with custom hex-encoded tokens generated using cryptographically secure random bytes. To access protected endpoints, you must:

1. **Register** a new user account using the `/api/v1/users/register` endpoint
2. **Login** with your credentials using the `/api/v1/users/login` endpoint to receive a hex token
3. **Include the token** in the `Authorization` header for all protected endpoints

Example token: `a7f3c9e8d4b2f1a6c8e7d3b9f2e5a8c1d6f4b7e9a2c5f8b3e6d9c2a5f1b4e7d8`

### Authentication Flow

```sh
# Step 1: Register a new user
curl -X POST "http://localhost:8080/api/v1/users/register" \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe", "email": "john@example.com", "password": "securepassword123"}'

# Step 2: Login to get your token (returns a hex-encoded secure random token)
curl -X POST "http://localhost:8080/api/v1/users/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "securepassword123"}'

# Step 3: Use the token in protected endpoints
curl -X GET "http://localhost:8080/api/v1/tasks" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

### Protected Endpoints

The following endpoints require authentication:
- All `/api/v1/users/:id` endpoints
- All `/api/v1/tasks/*` endpoints (except where noted otherwise)

Include your token in the Authorization header as: `Authorization: Bearer YOUR_HEX_TOKEN_HERE`

## JSON Payload Structures

### Register user payload
```json
{
    "name": "trav",
    "email": "trav@trav.com",
    "password": "pa55word"
}
```

### Login payload  
```json
{
    "email": "trav@trav.com",
    "password": "pa55word"
}
```

### Create task payload
```json
{
    "title": "empty bins",
    "description": "I need to take out the garbage, both bins"
}
```

### Update task payload
```json
{
  "title": "new title",
  "description": "new description", 
  "completed": "true"
}
```

## Endpoint Example Usage

### Health check
```sh
curl -X GET "http://localhost:8080/api/v1/healthcheck"
```

### Register a new user
```sh
curl -X POST "http://localhost:8080/api/v1/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "trav",
    "email": "trav@trav.com",
    "password": "pa55word"
}'
```

### User login
```sh
curl -X POST "http://localhost:8080/api/v1/users/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "trav@trav.com",
    "password": "pa55word"
  }'
```

### Get user by ID
```sh
curl -X GET "http://localhost:8080/api/v1/users/nefknvejj423vvefve" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

### Create a new task
```sh
curl -X POST "http://localhost:8080/api/v1/tasks/create" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE" \
  -d '{
    "title": "empty bins",
    "description": "I need to take out the garbage, both bins"
  }'
```

### Fetch all tasks
```sh
curl -X GET "http://localhost:8080/api/v1/tasks" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

### Get task by ID
```sh
curl -X GET "http://localhost:8080/api/v1/tasks/veirvniUGURBVUr45r" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

### Update task
```sh
curl -X PATCH "http://localhost:8080/api/v1/tasks/veirvniUGURBVUr45r" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE" \
  -d '{
  "title": "new title",
  "description": "new description", 
  "completed": false
  }'
```

### Mark task as complete
```sh
curl -X PUT "http://localhost:8080/api/v1/tasks/veirvniUGURBVUr45r/complete" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

### Delete task
```sh
curl -X DELETE "http://localhost:8080/api/v1/tasks/veirvniUGURBVUr45r" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

## Query Parameters (Searching by Title)

The `/api/v1/tasks` endpoint supports filtering and pagination through query parameters.

You can search by terms contained in the `title` of a task.

### Available Parameters

| Parameter   | Type    | Default | Description                 |
| ----------- | ------- | ------- | --------------------------- |
| `page`      | integer | 1       | Page number for pagination  |
| `page_size` | integer | 100     | Number of tasks per page    |
| `search`    | string  | ""      | Search term to filter tasks |

### Examples

#### Basic fetch all (default pagination)
```sh
curl -X GET "http://localhost:8080/api/v1/tasks" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

#### Fetch with pagination
```sh
curl -X GET "http://localhost:8080/api/v1/tasks?page=2&page_size=50" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

#### Search for specific tasks by title
```sh
curl -X GET "http://localhost:8080/api/v1/tasks?search=project" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

#### Combined filtering and pagination
```sh
curl -X GET "http://localhost:8080/api/v1/tasks?search=urgent&page=1&page_size=25" \
  -H "Authorization: Bearer YOUR_HEX_TOKEN_HERE"
```

### Response Format

The response includes both the filtered data and metadata about pagination:

```json
{
    "data": [
        {
            "ID": "68481b48324efcd3aca71764",
            "owner": "68481b2e324efcd3aca71735",
            "title": "Send holiday email",
            "description": "Send holiday greetings emails.",
            "completed": false,
            "created_at": "2025-06-10T11:47:20.734Z"
        }
    ],
    "metadata": {
        "page": 1,
        "limit": 1,
        "total": 9
    }
}
```

## Caching

This API implements Redis-based caching to improve performance and reduce database load for frequently accessed data.

### Cached Operations

- **Individual Task Fetching** - Tasks retrieved by ID are cached to speed up subsequent requests
- **Token Validation** - Authentication tokens are cached to avoid database lookups on every protected endpoint request

### Cache Behavior

The caching implementation helps reduce response times for:
- `/api/v1/tasks/:id` - Individual task lookups
- All protected endpoints - Token validation occurs without database queries for cached tokens

Cache invalidation occurs automatically when relevant data is modified (e.g., when tasks are updated or deleted).

## Contributing
Feel free to fork and submit PRs!

## License: `MIT`

If there are any concerns regarding the licence, please contact me at `travis.bozic@hotmail.com`.


## Image
Image by [Egon Elbre](https://github.com/egonelbre), used under CC0-1.0 license.