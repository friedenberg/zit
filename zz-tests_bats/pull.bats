#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
	load 'common.bash'
	# ... the remaining setup is unchanged

	# get the containing directory of this file
	# use $BATS_TEST_FILENAME instead of ${BASH_SOURCE[0]} or $0,
	# as those will point to the bats executable's location or the preprocessed file respectively
	DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" >/dev/null 2>&1 && pwd)"
	# make executables in src/ visible to PATH
	PATH="$DIR/../:$PATH"
	PATH="$DIR/../build/:$PATH"

	# for shellcheck SC2154
	export output
}

function pull { # @test
	wd="$(mktemp -d)"

	(
		cd "$wd" || exit 1
		run_zit_init_disable_age
	)

	wd1="$(mktemp -d)"

	(
		cd "$wd1" || exit 1
		run_zit_init_disable_age
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

	assert_output '[one/uno@11327fbe60cabd2a9eabf4a37d541cf04b539f913945897efe9bab1e30784781 !md "to_add.md"]'

	cd "$wd1" || exit 1

	run_zit pull "$wd" @
	assert_output '[one/uno@11327fbe60cabd2a9eabf4a37d541cf04b539f913945897efe9bab1e30784781 !md "to_add.md"]'

	run_zit show one/uno
	assert_output "$(cat "$expected")"

	cd "$wd" || exit 1

	run_zit show one/uno
	assert_output "$(cat "$expected")"

	run_zit pull "$wd" @
	assert_output ''
}
