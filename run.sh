git reset --hard origin/main
git pull
docker compose up -d --build
docker compose logs -f