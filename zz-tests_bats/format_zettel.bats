#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function format_simple { # @test
	run_zit checkout !md:t
	assert_success

	cat >md.typ <<-EOM
		inline-akte = true
		[formatters.text]
		shell = [
		  "cat",
		]
	EOM

	run_zit checkin -delete .t
	assert_success

	run_zit show -format debug !md:t

	run_zit format-zettel one/uno
	assert_success
	assert_output - <<-EOM
		---
		# wow the first
		- tag-3
		- tag-4
		! md
		---

		last time
	EOM

	run_zit checkout one/uno
	assert_success
	cat >one/uno.zettel <<-EOM
		---
		# wow the second
		- tag-3
		- tag-4
		! md
		---

		last time but new
	EOM

	run_zit format-zettel one/uno
	assert_success
	assert_output - <<-EOM
		---
		# wow the second
		- tag-3
		- tag-4
		! md
		---

		last time but new
	EOM
}
