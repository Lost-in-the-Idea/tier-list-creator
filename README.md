# tier-list-creator

.env file will be required to run this program locally with the following:

```
DISCORD_CLIENT_ID=your_discord_client_id
DISCORD_CLIENT_SECRET=your_discord_client_secret
DB_NAME=your_database_name
DB_USER=your_postgres_user
DB_PASSWORD=your_postgres_password
DB_HOST=your_postgres_host
DB_PORT=your_postgres_port
APP_ENV=dev            # Use 'dev' for development, 'prod' for production

# Optional: Run database actions on startup (dev only). After running the server will close, so you will need to remove the values from the .env file after completion.
# Comma separated values: clear, migrate, seed
# Example: DB_ACTION=clear,migrate,seed
DB_ACTION=
```
