if ($args[0] -eq "build") {
  go build -a -installsuffix cgo -ldflags "-s -X github.com/deis/workflow-cli/version.BuildVersion=$(git rev-parse --short HEAD)" -o deis.exe .
} elseif ($args[0] -eq "test") {
  go test --cover --race -v $(glide novendor)
} elseif ($args[0] -eq "bootstrap") {
  glide install
} else {
  echo "Unknown command: '$args'"
  exit 1
}

exit $LASTEXITCODE
