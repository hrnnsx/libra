
# Libra Anime Library API

Libra is a RESTful API for managing a personal anime library. It allows users to discover anime, maintain their watchlists, track their watching progress, organize anime into groups, and identify anime scenes using image URLs.

The application integrates with AniList for anime data and trace.moe for anime scene identification.


[LIVE API PREVIEW](libra-production-e912.up.railway.app)

## Features

### Authentication and User Management
- User registration and login.
- JWT-based authentication.
- Retrieve and update the authenticated user's profile.
- Password hashing using bcrypt.

### Anime Discovery
- Browse popular anime with pagination.
- Search anime by title or AniList ID.
- Filter anime by season, season year, genre, status, and format.
- Retrieve detailed anime information.
- Identify anime from an image URL using trace.moe.

### Personal Anime Library
- Add anime from AniList to a personal library.
- Retrieve all anime in the library.
- Retrieve a specific library entry.
- Update watching status, episode progress, rating, notes, and dates.
- Remove anime from the library.
- Prevent duplicate anime entries in the same user's library.

### Group Management
- Create, retrieve, update, and delete anime groups.
- Add anime from the personal library to a group.
- Retrieve anime assigned to a group.
- Remove anime from a group without deleting it from the personal library.
- Enforce group ownership and prevent duplicate anime assignments.

### API Documentation
- Interactive API documentation using Swagger UI.
- Endpoint documentation generated with Swaggo.
- Bearer token authentication support in Swagger UI.

## Technology Stack

| Technology | Purpose |
|---|---|
| Go | Backend programming language |
| Gin | HTTP routing and request handling |
| GORM | Object-relational mapping |
| PostgreSQL | Relational database |
| Supabase | Hosted PostgreSQL database |
| JWT | Authentication |
| bcrypt | Password hashing |
| AniList GraphQL API | Anime catalog and metadata |
| trace.moe API | Anime scene identification |
| Swaggo | Swagger/OpenAPI documentation |

## Project Structure

```text
libra/
├── cmd/
│   └── api/
│       └── main.go
├── config/
├── docs/
├── external/
│   ├── anilist/
│   └── tracemoe/
├── internal/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   └── service/
├── migrations/
├── documentation/
├── go.mod
├── go.sum
└── README.md
```

### Architecture

The application separates HTTP handling, business logic, database access, and external API integration.

- **Handler:** Handles HTTP requests, validates input, and returns responses.
- **Service:** Implements application use cases and business rules.
- **Repository:** Handles database operations.
- **Model:** Defines the application's data structures.
- **Middleware:** Handles authentication and other request processing.
- **External:** Contains clients for third-party APIs.
- **Config:** Contains application configuration.
- **Migrations:** Contains SQL scripts for database schema changes.
- **Docs:** Contains generated Swagger documentation.

## Prerequisites

Before running the application, ensure that you have:

- Go installed.
- Access to a PostgreSQL database.
- A configured database connection supported by the application.
- A JWT secret configured through the application's environment configuration.

The application uses AniList and trace.moe as external services. Their APIs must be reachable for the corresponding features to work.

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd libra
```

Replace `<repository-url>` with the repository's Git URL.

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Environment Variables

Configure the environment variables required by the application's configuration and database connection.

At minimum, provide a secure `JWT_SECRET` for signing authentication tokens.

Example:

```env
JWT_SECRET=replace-with-a-secure-random-secret
```

Configure the database connection using the variable names expected by the project's existing configuration. Do not commit real credentials or production secrets to version control.

### 4. Set Up the Database

Create a PostgreSQL database and apply the SQL migration scripts in the `migrations/` directory.

The database schema includes the following tables:

- `users`
- `anime`
- `library_anime`
- `groups`
- `group_library_anime`

Apply migrations before starting the application. The database schema is managed through SQL migrations rather than automatic schema generation.

### 5. Run the Application

From the project root:

```bash
go run ./cmd/api
```

The application uses the host and port configured in its server setup.

### 6. Access Swagger UI

When the server is running, open:

```text
http://localhost:8080/swagger/index.html
```

Swagger UI provides an interactive interface for exploring and testing the documented API endpoints.

## API Endpoints

All routes are listed without an `/api/v1` prefix.

### Authentication

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| POST | `/auth/register` | Register a new user | No |
| POST | `/auth/login` | Log in and obtain a token | No |

### User Profile

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| GET | `/users/me` | Retrieve the current user's profile | Yes |
| PATCH | `/users/me` | Update the current user's profile | Yes |

### Anime Discovery

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| GET | `/animes` | Browse popular anime | No |
| GET | `/animes/search` | Search and filter anime | No |
| GET | `/animes/:external_id` | Retrieve anime details by AniList ID | No |
| POST | `/animes/identify` | Identify anime from an image URL | No |

Example browse request:

```http
GET /animes?page=1&per_page=20
```

Example search request:

```http
GET /animes/search?title=mahouka&season=SPRING&season_year=2014
```

Example anime identification request:

```http
POST /animes/identify
Content-Type: application/json

{
  "url": "https://example.com/anime-scene.jpg"
}
```

The image URL must be accessible to trace.moe.

### Personal Anime Library

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| POST | `/me/library/` | Add anime to the personal library | Yes |
| GET | `/me/library/` | Retrieve the personal library | Yes |
| GET | `/me/library/:id` | Retrieve a library entry | Yes |
| PATCH | `/me/library/:id` | Update a library entry | Yes |
| DELETE | `/me/library/:id` | Remove a library entry | Yes |

Example add-anime request:

```http
POST /me/library/
Content-Type: application/json
Authorization: Bearer <token>

{
  "external_id": "20464"
}
```

Supported watching statuses:

- `WATCHING`
- `COMPLETED`
- `PLAN_TO_WATCH`
- `ON_HOLD`
- `DROPPED`

### Group Management

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| POST | `/me/groups` | Create a group | Yes |
| GET | `/me/groups` | Retrieve all groups | Yes |
| GET | `/me/groups/:id` | Retrieve a group | Yes |
| PATCH | `/me/groups/:id` | Update a group | Yes |
| DELETE | `/me/groups/:id` | Delete a group | Yes |

### Group Anime Management

| Method | Endpoint | Description | Authentication |
|---|---|---|---|
| POST | `/me/groups/:id/animes` | Add a library anime to a group | Yes |
| GET | `/me/groups/:id/animes` | Retrieve anime in a group | Yes |
| DELETE | `/me/groups/:id/animes/:library_anime_id` | Remove an anime from a group | Yes |

Example add-anime-to-group request:

```http
POST /me/groups/1/animes
Content-Type: application/json
Authorization: Bearer <token>

{
  "library_anime_id": 1
}
```

Removing an anime from a group only removes the group association. It does not remove the anime from the user's personal library.

## Authentication

Protected endpoints require a valid JWT access token in the `Authorization` header.

```http
Authorization: Bearer <your-access-token>
```

Obtain a token by registering or logging in. Include that token when calling protected endpoints.

## Error Responses

Errors use a structured response containing an error code and a message.

Example:

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "invalid request body"
  }
}
```

Common error codes include:

| Error Code | HTTP Status | Description |
|---|---:|---|
| `INVALID_REQUEST` | 400 | Invalid request body or parameter |
| `MISSING_AUTHORIZATION` | 401 | Authorization header is missing |
| `INVALID_AUTHORIZATION` | 401 | Authorization header format is invalid |
| `INVALID_TOKEN` | 401 | Authentication token is invalid or expired |
| `INVALID_CREDENTIALS` | 401 | Invalid login credentials |
| `USER_NOT_FOUND` | 404 | User does not exist |
| `ANIME_NOT_FOUND` | 404 | Anime was not found |
| `GROUP_NOT_FOUND` | 404 | Group was not found |
| `USERNAME_EXISTS` | 409 | Username is already registered |
| `EMAIL_EXISTS` | 409 | Email is already registered |
| `ANIME_ALREADY_IN_LIBRARY` | 409 | Anime already exists in the personal library |
| `ANIME_ALREADY_IN_GROUP` | 409 | Anime already exists in the group |
| `EXTERNAL_API_ERROR` | 502 | External API request failed |
| `INTERNAL_SERVER_ERROR` | 500 | Unexpected server error |

## External APIs

### AniList

AniList provides anime catalog information, including titles, synopsis, cover images, episode counts, genres, formats, seasons, release years, and scores.

API endpoint:

```text
https://graphql.anilist.co
```

### trace.moe

trace.moe identifies anime scenes using an image URL and returns matching results, including AniList IDs, episode information when available, timestamps, similarity scores, and preview links.

API endpoint:

```text
https://api.trace.moe/search
```


## Postman Collection

A Postman collection is provided to simplify API testing and development.

### Import the Collection

1. Open Postman.
2. Click **Import**.
3. Select the collection file located at `documentation/postman/Libra.postman_collection.json`.
4. If an environment file is provided, import `documentation/postman/Libra.postman_environment.json` as well.

### Configuration

Configure the following environment variables in Postman:

| Variable | Example Value | Description |
|---|---|---|
| `base_url` | `http://localhost:8080` | Base URL of the running API |
| `token` | Empty initially | JWT access token obtained after login |

Update the variable values according to your local environment.

### Testing Authenticated Endpoints

1. Start the Libra API server.
2. Register a user or log in using the authentication endpoints.
3. Copy the returned JWT access token into the `token` environment variable.
4. Use the collection to test authenticated endpoints, including user profiles, personal libraries, and groups.

Ensure that the API server and database are configured before running the requests.

## Development

Download dependencies:

```bash
go mod download
```

Format Go source files:

```bash
go fmt ./...
```

Run tests:

```bash
go test ./...
```

Regenerate Swagger documentation after changing endpoint annotations:

```bash
swag init -g cmd/api/main.go
```

## License

MIT License

Copyright (c) 2026 Haerunnas

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
