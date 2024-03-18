#set -e
#
#until PGPASSWORD=$POSTGRES_PASSWORD psql -h "localhost" -U "$POSTGRES_USER" -c '\q'; do
# echo "Postgres is unavailable - sleeping"
# sleep 1
#done
#
#echo "Postgres is up - executing migrations"
#psql -h "localhost" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/000001_init.up.sql
#
#exec postgres

#!/bin/sh
set -e

# Ожидание доступности PostgreSQL
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "db" -U "$POSTGRES_USER" -c '\q'; do
 echo "Postgres is unavailable - sleeping"
 sleep 1
done

echo "Postgres is up - executing migrations"

# Выполнение миграций
psql -h "db" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /schema/000001_init.up.sql

# Запуск основного процесса (в данном случае, postgres)
exec postgres