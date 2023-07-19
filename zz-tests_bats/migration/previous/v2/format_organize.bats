#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output
}

teardown() {
	rm_from_version
}

function format_organize_right_align { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age
	assert_success

	to_add="$(mktemp)"
	cat - >"$to_add" <<-EOM
		# task
		## urgency
		### urgency-1
		- [!md]
		- [-zz-archive]
		### -2
	EOM

	expected="$(mktemp)"
	cat - >"$expected" <<-EOM

		              # task

		             ## urgency

		            ###        -1

		- [-zz-archive]
		- [!md        ]

		            ###        -2
	EOM

	run_zit format-organize -prefix-joints=true -refine "$to_add"
	assert_success
	assert_output "$(cat "$expected")"
}

function format_organize_left_align { # @test
	cd "$BATS_TEST_TMPDIR" || exit 1
	run_zit_init_disable_age

	to_add="$(mktemp)"
	cat - >"$to_add" <<-EOM
		# task
		## urgency
		### urgency-1
		### -2
	EOM

	expected="$(mktemp)"
	cat - >"$expected" <<-EOM

		# task

		 ## urgency

		  ### -1

		  ### -2

	EOM

	run_zit format-organize -prefix-joints=true -refine -right-align=false "$to_add"
	assert_success
	assert_output "$(cat "$expected")"
}
