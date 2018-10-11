load test_helper

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
