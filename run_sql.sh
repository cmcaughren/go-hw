#!/bin/bash

# Check if SQL file was provided
if [ $# -eq 0 ]; then
    echo "Usage: ./run_sql <sql_file>"
    exit 1
fi

# Source the .env file
source .env

# Run the SQL file
sqlcmd -S "$SERVER_NAME" \
       -d "$DB_NAME" \
       -U "$ADMIN_USER" \
       -P "$ADMIN_PASSWORD" \
       -i "$1"