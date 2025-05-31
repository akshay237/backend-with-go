#!/bin/sh

set -e

echo "load the env vars"
source /app/app.env
echo "run databse migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"