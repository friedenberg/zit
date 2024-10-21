#! /bin/bash -e

bats_require_minimum_version 1.5.0
load "$BATS_CWD/zz-tests_bats/test_helper/bats-support/load"
load "$BATS_CWD/zz-tests_bats/test_helper/bats-assert/load"
load "$BATS_CWD/zz-tests_bats/test_helper/bats-assert-additions/load"

export XDG_DATA_HOME="$BATS_TEST_TMPDIR/.xdg/data"
export XDG_CONFIG_HOME="$BATS_TEST_TMPDIR/.xdg/config"
export XDG_STATE_HOME="$BATS_TEST_TMPDIR/.xdg/state"
export XDG_CACHE_HOME="$BATS_TEST_TMPDIR/.xdg/cache"
export XDG_RUNTIME_HOME="$BATS_TEST_TMPDIR/.xdg/runtime"

# get the containing directory of this file
# use $BATS_TEST_FILENAME instead of ${BASH_SOURCE[0]} or $0,
# as those will point to the bats executable's location or the preprocessed file respectively
DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" >/dev/null 2>&1 && pwd)"
# make executables in build/ visible to PATH
PATH="$BATS_CWD/build:$PATH"

# {
#   pushd "$BATS_CWD" >/dev/null 2>&1
#   gmake build/zit || exit 1
# }

{
  pushd "$BATS_TEST_TMPDIR" || exit 1
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
  -debug no-tempdir-cleanup
  -abbreviate-zettel-ids=false
  -abbreviate-shas=false
  -predictable-hinweisen
  -print-typen=false
  -print-time=false
  -print-etiketten=true
  -print-empty-shas=true
  -print-flush=false
  -print-unchanged=false
  -print-bestandsaufnahme=false
  -boxed-description=true
)

export cmd_zit_def

function copy_from_version {
  DIR="$1"
  version="${2:-v$(zit store-version)}"
  cp -r "$DIR/migration/$version" "$BATS_TEST_TMPDIR"
  cd "$BATS_TEST_TMPDIR/$version" || exit 1
}

function rm_from_version {
  version="${2:-v$(zit store-version)}"
  chflags_and_rm
}

function chflags_and_rm {
  "$BATS_CWD/bin/chflags.bash" -R nouchg "$BATS_TEST_TMPDIR"
}

function run_zit {
  cmd="$1"
  shift
  #shellcheck disable=SC2068
  run zit "$cmd" ${cmd_zit_def[@]} "$@"
}

function run_zit_stderr_unified {
  cmd="$1"
  shift
  #shellcheck disable=SC2068
  run zit "$cmd" ${cmd_zit_def[@]} "$@"
}

function run_zit_init {
  run_zit init -yin <(cat_yin) -yang <(cat_yang)
  assert_success
}

function get_konfig_sha() {
  echo -n "facdee599b069eb9dae4b04079fbf1b3aaaed30fe587ccc3e6fa7b6ff680b1f0"
}

function run_zit_init_disable_age {
  run_zit init -yin <(cat_yin) -yang <(cat_yang) -age none "$@"
  assert_success
  assert_output - <<-EOM
	[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
	[konfig @$(get_konfig_sha)]
EOM
}
