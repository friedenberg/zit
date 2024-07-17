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
  -abbreviate-hinweisen=false
  -abbreviate-shas=false
  -predictable-hinweisen
  -print-typen=false
  -print-time=false
  -print-etiketten=true
  -print-empty-shas=true
  -print-flush=false
  -print-unchanged=false
  -print-bestandsaufnahme=false
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
  # chflags -R nouchg "$BATS_TEST_TMPDIR/$version"
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

function run_zit_init {
  run_zit init -yin <(cat_yin) -yang <(cat_yang)
  assert_success
}

function get_konfig_sha() {
  echo -n "d79235cbe5153286a03aa9ca62539297156e839d28b422b114fae59bedea6f40"
}

function run_zit_init_disable_age {
  run_zit init -yin <(cat_yin) -yang <(cat_yang) -age none "$@"
  assert_success
  assert_output - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[konfig@$(get_konfig_sha)]
	EOM

  # run bash -c 'find .zit/Objekten/Bestandsaufnahme -type f | wc -l'
  # assert_success
  # assert_output '2'

  # run cat .zit/Objekten/Bestandsaufnahme/*/*
  # assert_success
  # assert_output --regexp 'Tai [[:digit:]]+\.[[:digit:]]+'
  # assert_output --regexp 'Akte'

  # run bash -c "cat .zit/Objekten/Bestandsaufnahme/*/* | grep Akte | cut -f2 -d' ' | xargs zit cat-objekte"
  # assert_success
  # assert_output_cut -d' ' -f2- -- - <<-EOM
  # 2061821648.550326 Typ md b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
  # 2061821648.554151 Konfig konfig 62c02b6f59e6de576a3fcc1b89db6e85b75c2ff7820df3049a5b12f9db86d1f5 c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685
  # EOM
}
