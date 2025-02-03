#! /bin/bash -e

load "$BATS_CWD/zz-tests_bats/test_helper/bats-support/load"
load "$BATS_CWD/zz-tests_bats/test_helper/bats-assert/load"
load "$BATS_CWD/zz-tests_bats/test_helper/bats-assert-additions/load"

set_xdg() {
  loc="$1"
  export XDG_DATA_HOME="$loc/.xdg/data"
  export XDG_CONFIG_HOME="$loc/.xdg/config"
  export XDG_STATE_HOME="$loc/.xdg/state"
  export XDG_CACHE_HOME="$loc/.xdg/cache"
  export XDG_RUNTIME_HOME="$loc/.xdg/runtime"
}

set_xdg "$BATS_TEST_TMPDIR"

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
  pushd "$BATS_TEST_TMPDIR" >/dev/null || exit 1
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
  -predictable-zettel-ids
  -print-typen=false
  -print-time=false
  -print-etiketten=true
  -print-empty-shas=true
  -print-flush=false
  -print-unchanged=false
  -print-bestandsaufnahme=false
  -boxed-description=true
  -print-colors=false
)

export cmd_zit_def

function copy_from_version {
  DIR="$1"
  version="${2:-v$(zit info store-version)}"
  rm -rf "$BATS_TEST_TMPDIR/.xdg"
  cp -r "$DIR/migration/$version/.xdg" "$BATS_TEST_TMPDIR/.xdg"
}

function rm_from_version {
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
  run_zit init \
    -yin <(cat_yin) \
    -yang <(cat_yang) \
    -lock-internal-files=false \
    "$@"
  assert_success
  assert_output - <<-EOM
[!md @$(get_type_blob_sha) !toml-type-v1]
[konfig @$(get_konfig_sha) !toml-config-v1]
EOM

  run_zit_init_workspace
}

function run_zit_init_workspace {
  run_zit init-workspace
}

function get_konfig_sha() {
  echo -n "9ad1b8f2538db1acb65265828f4f3d02064d6bef52721ce4cd6d528bc832b822"
}

function get_type_blob_sha() {
  echo -n "b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16"
}

run_find() {
  run find . -maxdepth 2 ! -ipath './.xdg*' ! -iname '.zit-workspace'
}

function run_zit_init_disable_age {
  run_zit init \
    -yin <(cat_yin) \
    -yang <(cat_yang) \
    -age-identity none \
    -lock-internal-files=false \
    "$@"

  assert_success
  assert_output - <<-EOM
[!md @$(get_type_blob_sha) !toml-type-v1]
[konfig @$(get_konfig_sha) !toml-config-v1]
EOM

  run_zit cat-blob "$(get_konfig_sha)"
  assert_success
  assert_output

  run_zit init-workspace
}
