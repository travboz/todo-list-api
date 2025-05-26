# Todo List API
![Superhero Gopher - Project Title Image](https://raw.githubusercontent.com/egonelbre/gophers/63b1f5a9f334f9e23735c6e09ac003479ffe5df5/vector/superhero/standing.svg)

## Description

This is a RESTful API designed to manage personal to-do lists with secure user access. It supports full CRUD functionality for to-do items and includes user registration, login with token-based authentication, and protected routes. To-dos can be listed, filtered, and paginated, with robust validation, error handling, and security best practices.

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


## Getting Started

### Prerequisites
- Docker
- Docker Compose
- Go (1.18+ recommended)

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
3. Run docker container containing MongoDB instance:
    ```sh
    make compose/up
    ```
4. Seed MongoDB instance:
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

| Method    | Endpoint           | Description                    |
|-----------|--------------------|--------------------------------|
| `GET`     | `/health`          | Health check                   |

## Example usage

### JSON payload structures

#### Create article payload

```json
{
  "content": "this is the content for a new article",
  "tags": ["these", "are", "the", "tags"]
}
```

#### Update user payload

```json
{
  "content": "this is the NEW CONTENT for an existinf article",
  "tags": ["where", "did", "the", "old", "tags", "go?"]
}
```

### Endpoint example usage
#### Create a user
```sh
curl -X POST "http://localhost:8080/articles" \
     -H "Content-Type: application/json" \
     -d '{
        "content": "this is the content for a new article",
        "tags": ["these", "are", "the", "tags"]
     }'
```

#### Update a user
```sh
curl -X POST "http://localhost:8080/users/67a0a3eef39fc03fe52450b5" \
     -H "Content-Type: application/json" \
     -d '{
        "content": "this is the NEW CONTENT for an existinf article",
        "tags": ["where", "did", "the", "old", "tags", "go?"]
      }'
```

#### Get a user by id
```sh
curl -X GET "http://localhost:7666/articles/67a0a3eef39fc03fe52450b5"
```

#### Fetch all users
```sh
curl http://localhost:7666/articles
```

#### Delete a user
```sh
curl -X DELETE "http://localhost:7666/articles/67a0a3eef39fc03fe52450b5"
```

## Contributing
Feel free to fork and submit PRs!

## License: `MIT`


If there are any concerns regarding the licence, please contact me at `travis.bozic@hotmail.com`.


## Image
Image by [Egon Elbre](https://github.com/egonelbre), used under CC0-1.0 license.
