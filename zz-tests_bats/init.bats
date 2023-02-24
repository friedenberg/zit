#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

function init_and_deinit { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	[[ -f .zit/KonfigAngeboren ]]

	run_zit deinit
	assert_success
}

function init_and_init { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init

	{
		echo "---"
		echo "# wow"
		echo "- tag"
		echo "! md"
		echo "---"
		echo
		echo "body"
	} >to_add

	run_zit new -edit=false to_add

	run_zit show one/uno
	assert_output "$(cat to_add)"

	run_zit init -yin <(cat_yin) -yang <(cat_yang)
	assert_failure

	run_zit init
	assert_output --partial '.zit/KonfigAngeboren already exists, not overwriting'
	assert_output --partial '.zit/KonfigErworben already exists, not overwriting'

	run zit show one/uno
	assert_output "$(cat to_add)"
}
