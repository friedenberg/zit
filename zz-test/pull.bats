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

# cmd_zit_def=(
# 	-abbreviate-hinweisen=false
# 	-predictable-hinweisen
# 	-print-typen=false
# )

# cmd_zit_add=(
# 	zit
# 	add
# 	"${cmd_zit_def[@]}"
# )

function pull { # @test
	wd="$(mktemp -d)"

	(
		cd "$wd" || exit 1

		run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
		assert_success
	)

	wd1="$(mktemp -d)"

	(
		cd "$wd1" || exit 1

		run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
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

	run zit new \
		-predictable-hinweisen \
		-abbreviate-hinweisen=false \
		-edit=false \
		"$expected"

	assert_output --partial '          (new) [one/uno@d !md "to_add.md"]'

	cd "$wd1" || exit 1

	run zit pull -abbreviate-hinweisen=false -all "$wd"
	assert_output '          (new) [one/uno@d !md "to_add.md"]'

	run zit show one/uno
	assert_output "$(cat "$expected")"

	cd "$wd" || exit 1

	run zit show one/uno
	assert_output "$(cat "$expected")"

	run zit pull -abbreviate-hinweisen=false -all "$wd"
	assert_output ''
}
