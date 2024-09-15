#!/bin/bash

# Navigate to the root directory
cd "$(dirname "$0")/.."

# Run Docker Compose
docker compose up --build
