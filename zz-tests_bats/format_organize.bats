#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

function format_organize_right_align { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "# task"
		echo "## urgency"
		echo "### urgency-1"
		echo "### -2"
	} >"$to_add"

	expected="$(mktemp)"
	{
		echo
		echo "    # task"
		echo
		echo "   ## urgency"
		echo
		echo "  ###        -1"
		echo
		echo "  ###        -2"
		echo
	} >"$expected"

	run_zit format-organize -prefix-joints=true -refine "$to_add"
	assert_output "$(cat "$expected")"
}

function format_organize_left_align { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "# task"
		echo "## urgency"
		echo "### urgency-1"
		echo "### -2"
	} >"$to_add"

	expected="$(mktemp)"
	{
		echo
		echo "# task"
		echo
		echo " ## urgency"
		echo
		echo "  ### -1"
		echo
		echo "  ### -2"
		echo
	} >"$expected"

	run_zit format-organize -prefix-joints=true -refine -right-align=false "$to_add"
	assert_output "$(cat "$expected")"
}
