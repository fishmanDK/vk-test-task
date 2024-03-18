set -e

# Wait for PostgreSQL to be ready
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "localhost" -U "$POSTGRES_USER" -c '\q'; do
 echo "Postgres is unavailable - sleeping"
 sleep 1
done

# Execute SQL migration files
echo "Postgres is up - executing migrations"
psql -h "localhost" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/000001_init.up.sql
psql -h "localhost" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/000001_init.down.sql

# Start PostgreSQL
exec postgres