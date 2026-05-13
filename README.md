# tier-list-creator

## Setup

### 1. Create a `.env` file at the repo root

```
DISCORD_CLIENT_ID=your_discord_client_id
DISCORD_CLIENT_SECRET=your_discord_client_secret
DB_NAME=your_database_name
DB_USER=your_postgres_user
DB_PASSWORD=your_postgres_password
DB_HOST=postgres
DB_PORT=5432
APP_ENV=dev

# Optional: runs a database action on startup then exits (dev only)
# Values: migrate, seed, clear
DB_ACTION=
```

`DB_HOST` must be `postgres` when running via Docker (it references the postgres service name).

### 2. Run with Docker

```
docker compose up --build
```

- Frontend: http://localhost:4200
- Backend: http://localhost:8080
- Postgres: localhost:5432

Both the backend and frontend support hot reload - edits to source files take effect without restarting containers.

### 3. Database actions

To run a database action (migrate, seed, clear) without touching the `.env`:

```
docker compose run --rm -e DB_ACTION=migrate,seed backend
```
