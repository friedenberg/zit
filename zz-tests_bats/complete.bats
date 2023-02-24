#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

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
