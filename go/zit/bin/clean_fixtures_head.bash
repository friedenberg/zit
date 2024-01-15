#! /bin/bash -e
#
dir="$(git rev-parse --show-toplevel)"

git reset HEAD "$dir"/zz-tests_bats/migration/v*/
chflags -R nouchg "$dir"/zz-tests_bats/migration/v*/
git clean -fd "$dir"/zz-tests_bats/migration/v*/
git co "$dir"/zz-tests_bats/migration/v*/

