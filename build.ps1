$v = $args[0]
$containerRegistry = "ghcr.io/vediatoni"
$packages = @("input","background") #array

$packages | ForEach-Object {
    $package = $_
    $tag = $containerRegistry + "/$package" + ":$v"
    docker build . -f .\build\$package.Dockerfile -t $tag
    docker push $tag
}