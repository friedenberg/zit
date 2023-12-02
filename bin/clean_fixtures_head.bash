#! /bin/bash -e
#
chflags -R nouchg zz-tests_bats/migration/v*/
git clean -fd zz-tests_bats/migration/v*/
git co zz-tests_bats/migration/v*/

