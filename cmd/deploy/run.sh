SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)
cd "$SCRIPT_DIR/../../"
docker compose down --remove-orphans
git pull
git reset --hard origin/main
docker compose up -d --build
docker compose logs -f