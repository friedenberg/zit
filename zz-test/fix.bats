#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
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

cat_yin() (
	echo "one"
	echo "two"
	echo "three"
	echo "four"
	echo "five"
	echo "six"
)

cat_yang() (
	echo "uno"
	echo "dos"
	echo "tres"
	echo "quatro"
	echo "cinco"
	echo "seis"
)

cmd_zit_def=(
	-abbreviate-hinweisen=false
	-predictable-hinweisen
	-print-typen=false
)

cmd_zit_new=(
	zit
	new
	"${cmd_zit_def[@]}"
)

cmd_zit_organize=(
	zit
	organize
	"${cmd_zit_def[@]}"
	-right-align=false
	-prefix-joints=true
	-metadatei-header=false
	-refine=true
)

function outputs_organize_one_etikett { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"
	assert_output '          (new) [one/uno@5 !md "wow"]'

	run zit show one/uno
	# assert_output "$(cat "$to_add")"

	run zit show ok
	assert_output "$(cat "$to_add")"

	run zit expand-hinweis o/u
	assert_output 'one/uno'

	expected_organize="$(mktemp)"
	{
		echo
		echo "# ok"
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run "${cmd_zit_organize[@]}" -mode output-only ok
	assert_output "$(cat "$expected_organize")"
}
