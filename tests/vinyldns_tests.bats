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

@test "group-update (when the group exists)" {
  fixture="$(cat tests/fixtures/group_updated)"
  ok_group=$($ew --op json group --name "ok-group")
  updated_group="$(echo ${ok_group} | sed 's/test@test.com/update@update.com/g')"
  run $ew group-update --json "${updated_group}"

  [ "${output}" = "${fixture}" ]
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

@test "search-record-sets (when the search returns results)" {
  fixture="$(cat tests/fixtures/search_with_results)"
  $ew search-record-sets \
    --record-name-filter "so*" \
    --record-type-filter "CNAME" \
    --record-type-filter "mx" \
    --max-items "50" \
    --name-sort "DESC" | grep "${fixture}"
}

@test "search-record-sets (when the search returns no results)" {
  run $ew search-record-sets \
    --record-name-filter "asdf" \
    --record-type-filter "CNAME" \
    --record-type-filter "mx" \
    --max-items "50" \
    --name-sort "DESC"
  fixture="$(cat tests/fixtures/search_with_no_results)"
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

@test "ping" {
  run $ew ping

  [ "$status" -eq 0 ]
  [ "$output" = "PONG" ]
}

@test "health" {
  run $ew health

  [ "$status" -eq 0 ]
  [ "$output" = "OK" ]
}

@test "color" {
  run $ew color

  [ "$status" -eq 0 ]
  echo "$output" | grep -E "blue|green"
}

@test "metrics-prometheus" {
  run $ew metrics-prometheus

  [ "$status" -eq 0 ]
  [ -n "$output" ]
}

@test "status" {
  run $ew --output=json status

  [ "$status" -eq 0 ]
  echo "$output" | grep '"processingDisabled"'
}

@test "status-update (requires admin)" {
  run $ew status-update --processing-disabled true

  [ "$status" -eq 1 ]
  echo "$output" | grep "/status?processingDisabled=true"
}

@test "zone-backend-ids" {
  run $ew --output=json zone-backend-ids

  [ "$status" -eq 0 ]
  echo "$output" | grep '\['
}

@test "zone-details" {
  run $ew --output=json zone-details --zone-name "ok."

  [ "$status" -eq 0 ]
  echo "$output" | grep '"name":"ok\.'
}

@test "zone-changes-failure" {
  run $ew --output=json zone-changes-failure

  [ "$status" -eq 0 ]
  echo "$output" | grep '"failedZoneChanges"'
}

@test "zones-deleted-changes" {
  run $ew --output=json zones-deleted-changes

  [ "$status" -eq 0 ]
  echo "$output" | grep '"zonesDeletedInfo"'
}

@test "zone-acl-rule-create (unknown zone)" {
  run $ew zone-acl-rule-create \
    --zone-id "does-not-exist" \
    --json '{"accessLevel":"Read","recordTypes":["A"],"groupId":"global-acl-group-id"}'

  [ "$status" -eq 1 ]
  echo "$output" | grep "/zones/does-not-exist/acl/rules"
}

@test "zone-acl-rule-delete (unknown zone)" {
  run $ew zone-acl-rule-delete \
    --zone-id "does-not-exist" \
    --json '{"accessLevel":"Read","recordTypes":["A"],"groupId":"global-acl-group-id"}'

  [ "$status" -eq 1 ]
  echo "$output" | grep "/zones/does-not-exist/acl/rules"
}

@test "record-set-count" {
  run $ew --output=json record-set-count --zone-name "ok."

  [ "$status" -eq 0 ]
  echo "$output" | grep '"count"'
}

@test "record-set-history" {
  run $ew --output=json record-set-history \
    --zone-name "ok." \
    --record-set-fqdn "some-cname.ok." \
    --record-set-type "CNAME"

  [ "$status" -eq 0 ]
  echo "$output" | grep '"recordSetChanges"'
}

@test "record-set-changes-failure" {
  run $ew --output=json record-set-changes-failure --zone-name "ok."

  [ "$status" -eq 0 ]
  echo "$output" | grep '"failedRecordSetChanges"'
}

@test "record-set-update (unknown record set)" {
  run $ew record-set-update \
    --json '{"zoneId":"does-not-exist","id":"does-not-exist","name":"nope","type":"CNAME","ttl":300,"records":[{"cname":"test.com"}]}'

  [ "$status" -eq 1 ]
  echo "$output" | grep "/zones/does-not-exist/recordsets/does-not-exist"
}

@test "group-change (unknown group change)" {
  run $ew group-change --group-change-id "does-not-exist"

  [ "$status" -eq 1 ]
  echo "$output" | grep "/groups/change/does-not-exist"
}

@test "group-valid-domains" {
  run $ew --output=json group-valid-domains

  [ "$status" -eq 0 ]
  echo "$output" | grep '\['
}

@test "user (unknown user)" {
  run $ew user --user-id "does-not-exist"

  [ "$status" -eq 1 ]
  echo "$output" | grep "/users/does-not-exist"
}

@test "user-lock (unknown user)" {
  run $ew user-lock --user-id "does-not-exist"

  [ "$status" -eq 1 ]
  echo "$output" | grep "/users/does-not-exist/lock"
}

@test "user-unlock (unknown user)" {
  run $ew user-unlock --user-id "does-not-exist"

  [ "$status" -eq 1 ]
  echo "$output" | grep "/users/does-not-exist/unlock"
}

@test "batch-change-approve (unknown batch change)" {
  run $ew batch-change-approve --batch-change-id "does-not-exist"

  [ "$status" -eq 1 ]
  echo "$output" | grep "/zones/batchrecordchanges/does-not-exist/approve"
}

@test "batch-change-reject (unknown batch change)" {
  run $ew batch-change-reject --batch-change-id "does-not-exist"

  [ "$status" -eq 1 ]
  echo "$output" | grep "/zones/batchrecordchanges/does-not-exist/reject"
}

@test "batch-change-cancel (unknown batch change)" {
  run $ew batch-change-cancel --batch-change-id "does-not-exist"

  [ "$status" -eq 1 ]
  echo "$output" | grep "/zones/batchrecordchanges/does-not-exist/cancel"
}
