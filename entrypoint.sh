set -e

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "db" -U "$POSTGRES_USER" -c '\q'; do
 echo "Postgres is unavailable - sleeping"
 sleep 1
done

echo "Postgres is up - executing migrations"

psql -h "db" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /schema/000001_init.up.sql

exec postgres