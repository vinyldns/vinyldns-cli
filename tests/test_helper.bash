setup() {
  ew="bin/vinyldns \
    --host http://localhost:9000 \
    --access-key=okAccessKey \
    --secret-key=okSecretKey"
}

teardown() {
  echo $output
}
