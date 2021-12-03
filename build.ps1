docker build . -f .\build\input.Dockerfile -t ghcr.io/vediatoni/input:latest
docker push ghcr.io/vediatoni/input:latest

docker build . -f .\build\input.Dockerfile -t ghcr.io/vediatoni/background:latest
docker push ghcr.io/vediatoni/background:latest