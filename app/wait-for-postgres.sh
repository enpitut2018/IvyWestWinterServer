#!/bin/sh

set -e

host="$1"
shift
cmd="$@"

until psql -h "$host" -U "postgres" &> /dev/null; do
  sleep 1
done

>&2 echo "Postgres is up"
exec $cmd
