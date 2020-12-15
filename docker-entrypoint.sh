#!/bin/sh
if [ $# -eq 0 ]
then
 MODE="static"
else
 MODE=$1
fi
echo "postman-mock-server mode selected: " $MODE
exec /app/postman-mockserver $MODE --config /app/config/config.yaml