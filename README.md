# tier-list-creator

## Prerequisites

- Docker Desktop
- Git

## Setup

### 1. Discord Application

1. Create a new discord application at the [Discord Developer Portal](https://discord.com/developers/applications).
2. Under "OAuth2" > Add a redirect URI: `http://localhost:8080/api/auth/discord/callback`
3. Copy the Client ID and Client Secret, these are needed for the environment variables.
4. Dev: Redirect Links
http://localhost:8080/api/auth/discord/redirect
http://localhost:8080/api/auth/discord/callback

### 2. Create a `.env` file at the repo root

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
# Values: migrate, seed, clear (comma separated for multiple actions)
DB_ACTION=
```

`DB_HOST` must be `postgres` when running via Docker (it references the postgres service name).

### 3. Run with Docker

```
docker compose up --build
```

- Frontend: http://localhost:4200
- Backend: http://localhost:8080
- Postgres: localhost:5432

Both the backend and frontend support hot reload - edits to source files take effect without restarting containers. Switching Git branches is reflected instantly in running containers via the volume mounts.

### 4. Database actions

`DB_ACTION` runs on startup then exits the server. Do not set this in your `.env` during normal operation.

Available Actions:

- 'clear' - Clears all data from the database (dev only).
- 'migrate' - Runs database migrations to update the schema.
- 'seed' - Seeds the database with initial data (dev only)

To run a database action (migrate, seed, clear) without touching the `.env`:

```
docker compose run --rm -e DB_ACTION=migrate backend
docker compose run --rm -e DB_ACTION=migrate,seed backend
docker compose run --rm -e DB_ACTION=clear,migrate,seed backend
```

`seed` and `clear` are blocked when `APP_ENV` is not `dev`.

### 5. Resetting the database

**This is destructive, all database data will be permanently lost.**

```
docker compose down -v
docker compose up --build
docker compose run --rm -e DB_ACTION=migrate,seed backend
```

`docker compose down -v` wipes the database volume. `migrate` recreates the schema, `seed` is optional.

### 6. Useful notes

Service names are `backend`, `frontend`, and `postgres` (for database commands). Use these in place of `<service_name>` in the commands below.

```
# Start and build containers
docker compose up --build

# Start without rebuilding
docker compose up

# Stop containers, keep database data
docker compose down

# Stop containers and wipe database volume (destructive)
docker compose down -v

# Rebuild and restart a single service
docker compose up --build <service_name>

# Restart a single service without rebuilding
docker compose restart <service_name>

# Tail logs for a service
docker compose logs -f <service_name>
```
