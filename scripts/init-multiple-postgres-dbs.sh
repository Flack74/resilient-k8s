#!/bin/bash

set -e
set -u

# Create the chaos_platform database
echo "Creating database 'chaos_platform'"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE chaos_platform;
    GRANT ALL PRIVILEGES ON DATABASE chaos_platform TO $POSTGRES_USER;
EOSQL

# Create the grafana database
echo "Creating database 'grafana'"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE grafana;
    GRANT ALL PRIVILEGES ON DATABASE grafana TO $POSTGRES_USER;
EOSQL

echo "Databases created successfully"