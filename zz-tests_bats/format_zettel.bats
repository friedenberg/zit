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

	# run cat .zit/Objekten/Akten/*/*
	# assert_output ''

	run_zit checkin -delete .t
	assert_success
	assert_output - <<-EOM
		[!md@f493a68ab71003dc5f1aaca8a2c5f90a013c868a77574b7e8f3dfb94f5c8cfd7]
		          deleted [md.typ]
	EOM

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
