#!/bin/bash

echo "Current working directory: $(pwd)"
echo "Directory contents:"
ls -la

echo "Waiting for ClickHouse to be ready (max ${MAX_RETRIES} attempts)..."
for i in $(seq 1 $MAX_RETRIES); do
  if curl -s "http://${CLICKHOUSE_HOST}:8123/ping" | grep -q "Ok"; then
    echo "ClickHouse is ready!"
    break
  fi
  echo "Attempt ${i}/${MAX_RETRIES}: ClickHouse not ready - sleeping ${RETRY_DELAY}s"
  sleep $RETRY_DELAY
  
  if [ $i -eq $MAX_RETRIES ]; then
    echo "ClickHouse not ready after ${MAX_RETRIES} attempts, aborting"
    exit 1
  fi
done



migrate -path . -database "clickhouse://${CLICKHOUSE_HOST}:${CLICKHOUSE_PORT}?username=${CLICKHOUSE_USER}&password=${CLICKHOUSE_PASSWORD}&database=${CLICKHOUSE_DB}&x-multi-statement=true" -verbose up 

# #!/bin/bash

# echo "Starting ClickHouse migrations..."
# echo "Current working directory: $(pwd)"
# echo "Directory contents:"
# ls -la

# # Ждем пока ClickHouse станет доступен
# until clickhouse client --host "${CLICKHOUSE_HOST}" --user "${CLICKHOUSE_USER}" --password "${CLICKHOUSE_PASSWORD}" --query "SELECT 1" >/dev/null 2>&1; do
#   echo "Waiting for ClickHouse to be ready..."
#   sleep 2
# done

# echo "ClickHouse is ready, applying migrations..."

# # Применяем все SQL файлы из текущей директории
# for migration_file in *.sql; do
#   echo "Applying migration: $migration_file"
#   clickhouse client --host "${CLICKHOUSE_HOST}" \
#                    --user "${CLICKHOUSE_USER}" \
#                    --password "${CLICKHOUSE_PASSWORD}" \
#                    --database "${CLICKHOUSE_DB}" \
#                    --multiquery < "$migration_file"
  
#   if [ $? -ne 0 ]; then
#     echo "Migration failed: $migration_file"
#     exit 1
#   fi
# done

# echo "All ClickHouse migrations applied successfully!"