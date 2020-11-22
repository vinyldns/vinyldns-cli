load test_helper

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
