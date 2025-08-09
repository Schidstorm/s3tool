#!/bin/bash

if docker ps -q -f name=minio; then
    echo "Stopping existing MinIO container..."
    docker stop minio
    docker rm minio
fi

echo "Starting MinIO container..."
docker run -d --name minio -p 9000:9000 minio/minio server start

echo "Writing MinIO configuration to ~/.s3tool/minio.yaml..."
mkdir -p ~/.s3tool

cat > ~/.s3tool/minio.yaml <<EOF
access_key_id: minioadmin
secret_access_key: minioadmin
base_endpoint: http://localhost:9000
use_path_style: true
region: ignore

EOF