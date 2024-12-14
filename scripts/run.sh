#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

APP_BIN_PATH="bin/main"

if [ ! -f $APP_BIN_PATH ]; then
  echo -e "${RED}Приложение не скомпилировано или не найдено по пути $APP_BIN_PATH${NC}"
  exit 1
fi

echo -e "${YELLOW}Запуск приложения...${NC}"
./$APP_BIN_PATH &
if [ $? -ne 0 ]; then
  echo -e "${RED}Ошибка при запуске приложения${NC}"
  exit 1
fi
echo -e "${GREEN}Приложение запущено${NC}"