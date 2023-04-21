#! /bin/bash -e

load "$BATS_CWD/zz-tests_bats/test_helper/bats-support/load"
load "$BATS_CWD/zz-tests_bats/test_helper/bats-assert/load"
load "$BATS_CWD/zz-tests_bats/test_helper/bats-assert-additions/load"

# get the containing directory of this file
# use $BATS_TEST_FILENAME instead of ${BASH_SOURCE[0]} or $0,
# as those will point to the bats executable's location or the preprocessed file respectively
DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" >/dev/null 2>&1 && pwd)"
# make executables in build/ visible to PATH
PATH="$BATS_CWD/build:$PATH"

{
	pushd "$BATS_CWD" >/dev/null 2>&1
	gmake build/zit || exit 1
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
	-abbreviate-shas=false
	-predictable-hinweisen
	-print-typen=false
	-print-time=false
)

function copy_from_version {
	DIR="$1"
	version="${2:-v$(zit store-version)}"
	cp -r "$DIR/migration/$version" "$BATS_TEST_TMPDIR"
	cd "$BATS_TEST_TMPDIR/$version" || exit 1
}

function rm_from_version {
	version="${2:-v$(zit store-version)}"
	# chflags -R nouchg "$BATS_TEST_TMPDIR/$version"
	chflags -R nouchg "$BATS_TEST_TMPDIR"
}

function run_zit {
	cmd="$1"
	shift
	#shellcheck disable=SC2068
	run zit "$cmd" ${cmd_zit_def[@]} "$@"
}

function run_zit_init {
	run_zit init -yin <(cat_yin) -yang <(cat_yang)
	assert_success
}

function run_zit_init_disable_age {
	run_zit init -yin <(cat_yin) -yang <(cat_yang) -disable-age
	assert_success
	assert_output - <<-EOM
		[!md@b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15]
		[konfig@f6d3d0874fd9475c2b7ac150f366cd211d847a8676ccabc35111cb357fd0c3b9]
	EOM
}
