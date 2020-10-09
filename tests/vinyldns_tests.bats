load test_helper

@test "groups (when none exist)" {
  fixture="$(cat tests/fixtures/groups_none)"
  $ew groups | grep "${fixture}"
}

@test "groups --output=json (when none exist)" {
  fixture="$(cat tests/fixtures/groups_none_json)"
  $ew --output=json groups | grep "${fixture}"
}

@test "group-create" {
  run $ew group-create \
    --json "$(cat tests/fixtures/group_create_json)"

  fixture="$(cat tests/fixtures/group_create)"

  [ "${output}" = "${fixture}" ]
}

@test "groups (when groups exist)" {
  fixture="$(cat tests/fixtures/groups)"
  $ew groups | grep "${fixture}"
}

@test "groups --output=json (when groups exist)" {
  fixture="$(cat tests/fixtures/groups_json)"

  $ew --output=json groups | grep "${fixture}"
}

@test "group (when the group exists)" {
  fixture="$(cat tests/fixtures/group)"

  $ew group --name "ok-group" | grep "${fixture}"
}

@test "group --output=json (when the group exists)" {
  fixture="$(cat tests/fixtures/group_json)"

  $ew --output=json group --name "ok-group" | grep "${fixture}"
}

@test "zones (when none exist)" {
  run $ew zones

  fixture="$(cat tests/fixtures/zones_none)"

  [ "${output}" = "${fixture}" ]
}

@test "zones --output=json (when none exist)" {
  run $ew --output=json zones

  fixture="$(cat tests/fixtures/zones_none_json)"

  [ "${output}" = "${fixture}" ]
}

@test "zone-create (with connection)" {
  run $ew zone-create \
    --name "ok." \
    --email "test@test.com" \
    --admin-group-name "ok-group" \
    --zone-connection-key-name "vinyldns." \
    --zone-connection-key "nzisn+4G2ldMn0q1CV3vsg==" \
    --zone-connection-primary-server "vinyldns-bind9" \
    --transfer-connection-key-name "vinyldns." \
    --transfer-connection-key "nzisn+4G2ldMn0q1CV3vsg==" \
    --transfer-connection-primary-server "vinyldns-bind9"

  fixture="$(cat tests/fixtures/zone_create_connection)"

  [ "${output}" = "${fixture}" ]
}

@test "zone-create (with no connection)" {
  run $ew zone-create \
    --name "vinyldns." \
    --email "admin@test.com" \
    --admin-group-name "ok-group"

  fixture="$(cat tests/fixtures/zone_create_no_connection)"

  [ "${output}" = "${fixture}" ]
}

@test "zone-create (with invalid zone connection)" {
  run $ew zone-create \
    --name "ok-invalid-connection." \
    --email "test@test.com" \
    --admin-group-name "ok-group" \
    --zone-connection-key "nzisn+4G2ldMn0q1CV3vsg==" \
    --zone-connection-primary-server "vinyldns-bind9"

  fixture="$(cat tests/fixtures/zone_create_invalid_zone_connection)"

  [ "${status}" -eq 1 ]
  [ "${output}" = "${fixture}" ]
}

@test "zone-create (with invalid transfer connection)" {
  run $ew zone-create \
    --name "ok-invalid-connection." \
    --email "test@test.com" \
    --admin-group-name "ok-group" \
    --transfer-connection-key "nzisn+4G2ldMn0q1CV3vsg==" \
    --transfer-connection-primary-server "vinyldns-bind9"

  fixture="$(cat tests/fixtures/zone_create_invalid_transfer_connection)"

  [ "${status}" -eq 1 ]
  [ "${output}" = "${fixture}" ]
}

@test "zone (when the zone exists)" {
  fixture="$(cat tests/fixtures/zone)"

  $ew zone --zone-name "ok." | grep "${fixture}"
}

@test "record-set-create (CNAME)" {
  run $ew record-set-create \
    --zone-name "ok." \
    --record-set-name "some-cname" \
    --record-set-type "CNAME" \
    --record-set-ttl "123" \
    --record-set-data "test.com"

  fixture="$(cat tests/fixtures/record_set_create_cname)"

  [ "${output}" = "${fixture}" ]
}

@test "record-set-create (MX)" {
  run $ew record-set-create \
    --zone-name "ok." \
    --record-set-name "some-mx" \
    --record-set-type "mx" \
    --record-set-ttl "123" \
    --record-set-data "3,test.com"

  fixture="$(cat tests/fixtures/record_set_create_mx)"

  [ "${output}" = "${fixture}" ]
}

@test "record-set-create (TXT)" {
  run $ew record-set-create \
    --zone-name "ok." \
    --record-set-name "some-txt" \
    --record-set-type "TXT" \
    --record-set-ttl "123" \
    --record-set-data "test TXT"

  fixture="$(cat tests/fixtures/record_set_create_txt)"

  [ "${output}" = "${fixture}" ]
}

@test "search-record-sets" {
  run $ew search-record-sets \
    --record-name-filter "so*" \
    --record-type-filter "CNAME" \
    --record-type-filter "mx" \
    --max-items "50" \
    --name-sort "DESC"
  record_sets[0]="some-cname"
  record_sets[1]="some-mx"
  for record_set in ${record_sets}
  do
    [[ "${output}" =~ "${record_set}" ]]
  done
}

@test "zone-sync (when the zone exists)" {
  # wait until the recently-created zone is in a state where it can be synced again
  sleep 10

  fixture="$(cat tests/fixtures/zone-sync)"

  $ew zone-sync --zone-name "ok." | grep "${fixture}"
}

@test "batch-change-create" {
  run $ew batch-change-create \
    --json "$(cat tests/fixtures/batch_change_create_json)"

  [ "$status" -eq 0 ]
}
