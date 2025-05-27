#!/bin/bash

# Application details
APP_NAME="clir"

# Remote server details
REMOTE_USER="root"
REMOTE_HOST="49.12.214.64"
REMOTE_DIR="/var/www/clir"

# Create directories
echo "Creating directories on remote server, if do not exist..."
ssh ${REMOTE_USER}@${REMOTE_HOST} "mkdir -p ${REMOTE_DIR}"

# RSync files
echo "Copying binary to ${REMOTE_HOST}${REMOTE_DIR}..."
rsync -avz --delete --exclude 'node_modules' --exclude 'package.json' --exclude 'package-lock.json' -e ssh "./site/" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}"

echo "Deployment completed"