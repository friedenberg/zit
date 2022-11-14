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

function add { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	f=to_add.md
	{
		echo test file
	} >"$f"

	run zit add \
		-abbreviate-hinweisen=false \
		-predictable-hinweisen \
		-dedupe \
		-delete \
		-etiketten zz-inbox-2022-11-14 \
		"$f"

	assert_output --partial '(created) [one/uno@b !md "to_add.md"]'
	assert_output --partial '(updated) [one/uno@d !md "to_add.md"]'
	assert_output --partial '[to_add.md] (deleted)'
}

function add_1 { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	f=to_add.md
	{
		echo test file
	} >"$f"

	run zit add \
		-predictable-hinweisen \
		-dedupe \
		-delete \
		-etiketten zz-inbox-2022-11-14 \
		"$f"

	assert_output --partial '(created) [o/u@b !md "to_add.md"]'
	assert_output --partial '(updated) [o/u@d !md "to_add.md"]'
	assert_output --partial '[to_add.md] (deleted)'
}
