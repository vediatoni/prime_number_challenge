$v = $args[0]
docker build . -f .\build\input.Dockerfile -t ghcr.io/vediatoni/input:$v
docker push ghcr.io/vediatoni/input:$v

docker build . -f .\build\background.Dockerfile -t ghcr.io/vediatoni/background:$v
docker push ghcr.io/vediatoni/background:$v