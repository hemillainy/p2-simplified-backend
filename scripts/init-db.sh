#!/bin/bash
echo "run db migration"
echo $BACKEND_USER
echo $BACKEND_PASSWORD
echo $BACKEND_DB_HOST
echo $BACKEND_DB_PORT
echo $BACKEND_DB_NAME
/app/migrate.linux-amd64 -path /app/migrations/ -database "$DB_SOURCE" up
./backend