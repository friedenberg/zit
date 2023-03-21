#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

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
	assert_output '[one/uno@11327fbe60cabd2a9eabf4a37d541cf04b539f913945897efe9bab1e30784781 !md "to_add.md"]'

	cd "$wd1" || exit 1

	run_zit pull "$wd" :
	assert_success
	assert_output '[one/uno@11327fbe60cabd2a9eabf4a37d541cf04b539f913945897efe9bab1e30784781 !md "to_add.md"]'

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
