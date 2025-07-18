#!/bin/bash

ENV_FILE="$HOME/.zshrc" # replace to ~/.bash_profile if using bash

declare -A env_vars=(
  ["LOCAL_SERVER_PORT"]="8765"
  ["BANGUMI_CLIENT_ID"]="bangumi APP ID"
  ["BANGUMI_CLIENT_SECRET"]="bangumi APP Secret"
  ["QBITTORRENT_SERVER"]="http://localhost:8767"
  ["QBITTORRENT_USERNAME"]="admin"
  ["QBITTORRENT_PASSWORD"]=""
  ["MIKAN_IDENTITY_COOKIE"]=""
)

for key in "${!env_vars[@]}"; do
  value="${env_vars[$key]}"
  if grep -q "export $key=" "$ENV_FILE"; then
    echo "Updating existing $key in $ENV_FILE"
    sed -i '' "s|^export $key=.*|export $key=\"$value\"|" "$ENV_FILE"
  else
    echo "Adding $key to $ENV_FILE"
    echo "export $key=\"$value\"" >> "$ENV_FILE"
  fi
done

echo "✅ Successfully added all envs to $ENV_FILE"
echo "🔄 Please run: source $ENV_FILE"
