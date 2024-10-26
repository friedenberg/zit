#! /bin/bash -e

dir_base="$(realpath "$(dirname "$0")")"
zit="$(realpath build/zit)"

if ! v="$("$zit" store-version)"; then
  echo "failed to get store version" >&2
  exit 1
fi

d="${1:-$dir_base/v$v}"

if [[ -d $d ]]; then
  ./bin/chflags.bash -R nouchg "$d"
  rm -rf "$d"
fi

cmd_bats=(
  bats
  --tap
  --no-tempdir-cleanup
  zz-tests_bats/migration/generate_fixture.bats
)

export BATS_TEST_TIMEOUT=3
if ! bats_run="$("${cmd_bats[@]}" 2>&1)"; then
  echo "$bats_run" >&2
  exit 1
else
  bats_dir="$(echo "$bats_run" | grep "BATS_RUN_TMPDIR" | cut -d' ' -f2)"
fi

mkdir -p "$d"
cp -r "$bats_dir/test/1/.xdg" "$d/.xdg"
