load test_helper

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

@test "update zone (when the zone exists)" {
  fixture="$(cat tests/fixtures/zone_updated)"
  ok_zone=$($ew --op json zone --zone-name "ok.")
  updated_zone="$(echo ${ok_zone} | sed 's/test@test.com/update@update.com/g')"
  run $ew zone-update \
    --json "${updated_zone}"

  [ "${output}" = "${fixture}" ]
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
