#! /bin/bash -e

git reset HEAD zz-tests_bats/migration/v*/
chflags -R nouchg zz-tests_bats/migration/v*/
git clean -fd zz-tests_bats/migration/v*/
git co zz-tests_bats/migration/v*/

