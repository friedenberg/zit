#! /bin/bash -x

dir="$(git rev-parse --show-toplevel)"
pushd "$dir" || exit

git reset HEAD zz-tests_bats/migration/v*/
./go/zit/bin/chflags.bash -R nouchg zz-tests_bats/migration/v*/
git clean -fd zz-tests_bats/migration/v*/
git co zz-tests_bats/migration/
