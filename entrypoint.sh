#!/bin/bash

# Run generate docs swagger
go install github.com/swaggo/swag/cmd/swag@latest
swag init --dir cmd/server/,internal --output docs

## Wait for MySQL to be ready
#until mysqladmin ping -h"$DB_HOST" -u"$DB_USER" -p"$DB_PASSWORD" --silent; do
#    echo "Waiting for database..."
#    sleep 2
#done
#
# Run the migrations
go run cmd/migrate.go up

# Air build
air -c .air.toml

# Start the application
exec "$@"
