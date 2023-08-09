#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output
}

function pull { # @test
	wd="$(mktemp -d)"

	(
		cd "$wd" || exit 1
		run_zit_init_disable_age
		assert_success
	)

	wd1="$(mktemp -d)"

	(
		cd "$wd1" || exit 1
		run_zit_init_disable_age
		assert_success
	)

	cd "$wd" || exit 1

	expected="$(mktemp)"
	{
		echo '---'
		echo '# to_add.md'
		echo '- zz-inbox-2022-11-14'
		echo '! md'
		echo '---'
		echo ''
		echo 'test file'
	} >"$expected"

	run_zit new \
		-edit=false \
		"$expected"

	assert_success
	assert_output - <<-EOM
		[-zz@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11-14@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add.md" zz-inbox-2022-11-14]
	EOM

	cd "$wd1" || exit 1

	# TODO fix race condition
	run_zit pull "$wd" :
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[-zz-inbox-2022-11-14@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022-11@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox-2022@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz-inbox@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-zz@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@55f8718109829bf506b09d8af615b9f107a266e19f7a311039d1035f180b22d4 !md "to_add.md" zz-inbox-2022-11-14]
	EOM

	run_zit show one/uno:z
	assert_success
	assert_output "$(cat "$expected")"

	cd "$wd" || exit 1

	run_zit show one/uno:z
	assert_success
	assert_output "$(cat "$expected")"

	run_zit pull "$wd" :
	assert_success
	assert_output ''
}
