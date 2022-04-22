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
}

function can_run_zit { # @test
	./build/zit
}

function provides_help_with_no_params { # @test
	run ./build/zit
	assert_output --partial 'No subcommand provided.'
}

function can_initialize_without_age { # @test
  yin="$(mktemp)"
  echo "one" > "$yin"
  echo "two" > "$yin"
  echo "three" > "$yin"

  yang="$(mktemp)"
  echo "uno" > "$yang"
  echo "dos" > "$yang"
  echo "tres" > "$yang"

	run ./build/zit init -disable-age -yin "$yin" -yang "$yang"
	assert_output --partial 'No subcommand provided.'
}
