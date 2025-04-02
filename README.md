Goserv is a Go-based web service that provides a RESTful API for managing chirps (short messages) and users.

## Features

- User management (registration, login, update)
- Chirp creation, retrieval, and deletion
- Authentication using JWT
- Chirp filtering by author
- Sorting chirps by creation time
- Admin metrics and reset functionality
- Health check endpoint

## Prerequisites

- Go 1.x (version used in development)
- PostgreSQL database

## Installation

1. Clone the repository:

```json
 git clone https://github.com/your-username/Goserv.git
```

1. CD goserv
2. Install dependencies:

```
go mod tidy
```

1. Set up your environment variables in a `.env` file:
    
    ```json
    DB_URL=your_database_connection_string
    PLATFORM=your_platform
    JWTsecret=your_jwt_secret
    POLKA_KEY=your_polka_key
    ```
    
2. Running The Application
    
    ```json
    go run .
    ```
    

The server will start on port 8080 by default.

## API Endpoints

- `GET /api/healthz`: Health check
- `GET /admin/metrics`: Get admin metrics
- `POST /admin/reset`: Reset admin metrics
- `POST /api/users`: Create a new user
- `POST /api/login`: User login
- `PUT /api/users`: Update user information
- `POST /api/chirps`: Create a new chirp
- `GET /api/chirps`: Get all chirps
- `GET /api/chirps?sort=desc` : Sort by Ascending or Descending
- `GET /api/chirps/{chirpID}`: Get a specific chirp
- `DELETE /api/chirps/{chirpID}`: Delete a chirp
- `GET /api/chirps?author_id={authorID}`: Get chirps by author
- `POST /api/refresh`: Refresh authentication token
- `POST /api/revoke`: Revoke authentication token
- `POST /api/polka/webhooks`: Handle Polka webhooks

## License

This project is licensed under the [MIT License](https://www.notion.so/rajeshmondal/LICENSE).