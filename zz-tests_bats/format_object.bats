#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

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
		[!md @21759bebd1a7937005f692b9394c0d2629361286b9fe837617e166c3ded687eb]
		          deleted [md.typ]
	EOM

	run_zit format-object one/uno text
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

	run_zit format-object one/uno.zettel text
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

function show_simple_one_zettel_binary { # @test
	echo "binary file" >file.bin
	run_zit add -delete file.bin
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [file.bin]
		[!bin @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627 !bin "file"]
	EOM

	run_zit format-object two/uno
	assert_success
	assert_output - <<-EOM
		---
		# file
		! b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627.bin
		---
	EOM
}
