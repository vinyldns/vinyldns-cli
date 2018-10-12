load test_helper

@test "zones (when none exist)" {
  result=$($ew zones)
  fixture=$(cat tests/fixtures/zones_none)

  echo $result

  [ "${result}" = "${fixture}" ]
}

@test "zone-create" {
  result=$($ew zone-create \
    --name "ok." \
    --email "test@test.com" \
    --admin-group-id "ok-group" \
    --zone-connection-key-name "vinyldns." \
    --zone-connection-key "nzisn+4G2ldMn0q1CV3vsg==" \
    --zone-connection-primary-server "vinyldns-bind9" \
    --transfer-connection-key-name "vinyldns." \
    --transfer-connection-key "nzisn+4G2ldMn0q1CV3vsg==" \
    --transfer-connection-primary-server "vinyldns-bind9"
  )

  echo $result

  fixture=$(cat tests/fixtures/zone_create)

  [ "${result}" = "${fixture}" ]
}

@test "groups" {
  result=$($ew groups)
  fixture=$(cat tests/fixtures/groups)

  [ "${result}" = "${fixture}" ]
}

@test "group" {
  result=$($ew group --group-id ok-group)
  fixture=$(cat tests/fixtures/group)

  [ "${result}" = "${fixture}" ]
}
