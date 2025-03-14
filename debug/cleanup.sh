# Clean up all docker images, containers, and volumes

docker rm -vf $(docker ps -aq)
docker system prune

# Free port 27017 from host mongo instance

systemctl stop mongod
