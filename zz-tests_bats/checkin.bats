#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	cp -r "$DIR/migration/$version" "$BATS_TEST_TMPDIR"
	cd "$BATS_TEST_TMPDIR/$version" || exit 1
	run_zit checkout @z,t,e

	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		! md
		---

		newest body
	EOM

	cat >one/dos.zettel <<-EOM
		---
		# dos wildly different
		! md
		---

		dos newest body
	EOM

	cat >md.typ <<-EOM
		inline-akte = false
	EOM
}

teardown() {
	chflags -R nouchg "$BATS_TEST_TMPDIR/$version"
}

function checkin_simple_one_zettel { # @test
	run_zit checkin one/uno.zettel
	assert_output '[one/uno@6e82467623a2aef20ec4c2207300d6c2adbc2711ad57a92d38f90946135a661d !md "wildly different"]'
}

function checkin_simple_all { # @test
	run_zit checkin .
	# assert_output --partial '[!md@0966bffa92f9391ec0874fe0bd5ed77b9ceddc45e36a866c71a3ccbb31711a71]'
	assert_output --partial '[!@0966bffa92f9391ec0874fe0bd5ed77b9ceddc45e36a866c71a3ccbb31711a71]'
	assert_output --partial '[one/uno@6e82467623a2aef20ec4c2207300d6c2adbc2711ad57a92d38f90946135a661d !md "wildly different"]'
	assert_output --partial '[one/dos@f69dde187bd082e8366587d2a55d2c7d44a892250acc9748d1aa62b87f0304e2 !md "dos wildly different"]'
}

function checkin_simple_typ { # @test
	skip
	run_zit checkin .t
	assert_output '[!@0966bffa92f9391ec0874fe0bd5ed77b9ceddc45e36a866c71a3ccbb31711a71]'
}
