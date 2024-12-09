#!/bin/bash
sudo rm .env
echo "remove last env"

sudo rm conf/app.conf
echo "remove last app.conf"

SOURCE_CONFIG="../config/smf-config.conf"
DEST_CONFIG="./conf/app.conf"
ENV_FILE="./.env"

# Buat direktori tujuan jika belum ada

# Ekstrak baris yang dimulai dengan "user-" dan tulis ke file tujuan
sudo grep -E '^(smfleadsuser-|db\.|AWS\.|server\.|service\.grpcport|smf\.token)' "$SOURCE_CONFIG" > "$DEST_CONFIG"


echo "app.conf has been added $SOURCE_CONFIG ke $DEST_CONFIG"

#CONFIG ENV FOR DOCKER COMPOSE
sudo grep '^smfleadsuser-port' "$SOURCE_CONFIG" | while IFS='=' read -r key value; do
    echo "SMFLEADS_USER_PORT=${value// /}" >> "$ENV_FILE"
done

echo ".env telah dibuat dengan nilai SMFLEADS_USER_PORT"


sudo docker system prune -A
echo "system prune"

# Remove old binary if exists
rm -f bee-smfleads-user
echo "Old binary removed"

# Build new binary
go build -o bee-smfleads-user
echo "Go build complete"

# Stop running container
sudo docker container stop bee-smfleads-user|| true
echo "Stopped container"

# Remove old container
sudo docker container rm bee-smfleads-user || true
echo "Removed container"

# Remove old image
sudo docker image rm bee-smfleads-user_app:latest || true
echo "Removed image"


# Build and run new container
sudo docker-compose up --build --remove-orphans -d
echo "Docker compose up complete"

if [ $? -eq 0 ]; then
    echo "Docker Compose successfully started."
else
    echo "Error starting Docker Compose."
    exit 1
fi

# Restart nginx
sudo service nginx restart
echo "Nginx restarted"

#sudo rm -rf config/
# List all Docker containers
sudo docker container ls -a
echo "Deployment complete"
