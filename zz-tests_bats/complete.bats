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

cmd_zit_def=(
	# -abbreviate-hinweisen=false
	-predictable-hinweisen
	-print-typen=false
)

function complete_show { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	expected="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$expected"

	run zit new "${cmd_zit_def[@]}" -edit=false -predictable-hinweisen -bezeichnung wow -etiketten ok
	assert_output '[o/u@5 "wow"] (created)'

	run zit show "${cmd_zit_def[@]}" one/uno
	assert_output "$(cat "$expected")"

	{
		echo "one/uno	Zettel: !md wow"
		echo "ok	Etikett"
	} >"$expected"

	run zit show -complete
	assert_output "$(cat "$expected")"
}