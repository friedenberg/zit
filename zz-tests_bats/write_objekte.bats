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

cmd_zit_write_objekte=(
	zit
	write-objekte
	"${cmd_zit_def[@]}"
)

cmd_zit_cat_objekte=(
	zit
	cat-objekte
	"${cmd_zit_def[@]}"
)

function write_objekte_none { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	run "${cmd_zit_write_objekte[@]}"
	assert_output ''
}

function write_objekte_null { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	run "${cmd_zit_write_objekte[@]}" - </dev/null
	assert_output 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -'
}

function write_objekte_one_file { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	run "${cmd_zit_write_objekte[@]}" <(echo wow)
	assert_output --partial 'f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63 /dev/fd/'

	run "${cmd_zit_cat_objekte[@]}" "f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63"
	assert_output "$(printf "%s\n" wow)"

	run zit cat -gattung akte
	assert_output --partial "f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63"
}

function write_objekte_one_file_one_stdin { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	run "${cmd_zit_write_objekte[@]}" <(echo wow) - </dev/null
	assert_output --partial 'f40cd21f276e47d533371afce1778447e858eb5c9c0c0ed61c65f5c5d57caf63 /dev/fd/'
	assert_output --partial 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -'
}
