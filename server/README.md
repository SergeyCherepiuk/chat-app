# Server Microservice

The Server microservice is responsible for authenticating users and performing CRUD operations for users and roles.

## Installation and usage

1. Clone the repository

2. Make sure your PostgreSQL and Redis databases are running

3. Porvide the way to connect to your databases in `.env` file (you can find example in `example.env`)

4. Run the server
```shell
go run main.go
```

## API Endpoints

The microservice exposes the following API endpoints:

### Authentication and users
- `POST /api/auth/login` (<i>username</i> and <i>password</i> required in the request body)

- `POST /api/auth/signup` (<i>username</i> and <i>password</i> required in the request body)

- `POST /api/auth/logout`

- `GET /user/me`

### Roles

- `GET /roles/{id}`

- `GET /roles`

- `POST /roles` (<i>name</i> required in the request body)

- `PUT /roles/{id}` (<i>name</i> required in the request body)

- `DELETE /roles/{id}`

- `DELETE /roles`