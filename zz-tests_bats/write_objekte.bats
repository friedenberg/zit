#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

function write_objekte_none { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	run_zit write-objekte
	assert_output ''
}

function write_objekte_null { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	run_zit write-objekte - </dev/null
	assert_output 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -'
}

function write_objekte_one_file { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	run_zit write-objekte <(echo wow)
	assert_output --partial 'f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63 /dev/fd/'

	run_zit cat-objekte "f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63"
	assert_output "$(printf "%s\n" wow)"

	run_zit cat -gattung akte
	assert_output --partial "f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63"
}

function write_objekte_one_file_one_stdin { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	run_zit write-objekte <(echo wow) - </dev/null
	assert_output --partial 'f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63 /dev/fd/'
	assert_output --partial 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -'
}
