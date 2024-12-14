#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="project-sem-1"
DB_USER="validator"
DB_PASSWORD="val1dat0r"
COMPILE_TO="bin/main"
COMPILE_FROM="cmd/main.go"

export PGPASSWORD=$DB_PASSWORD

echo -e "${YELLOW}Создание таблицы prices...${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/0001_create_prices.sql &> /dev/null
if [ $? -ne 0 ]; then
  echo -e "${RED}Ошибка при создании таблицы prices${NC}"
  exit 1
fi
echo -e "${GREEN}Таблица prices создана${NC}"

echo -e "${YELLOW}Компиляция приложения...${NC}"
go build -o $COMPILE_TO $COMPILE_FROM
if [ $? -ne 0 ]; then
  echo -e "${RED}Ошибка при компиляции приложения${NC}"
  exit 1
fi
echo -e "${GREEN}Приложение скомпилировано в $COMPILE_TO${NC}"