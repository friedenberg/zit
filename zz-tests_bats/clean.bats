#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	run_zit checkout :z,t,e
}

teardown() {
	rm_from_version "$version"
}

function clean_all { # @test
	run_zit clean .
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [md.typ]
		          deleted [one/]
		          deleted [one/dos.zettel]
		          deleted [one/uno.zettel]
		          deleted [tag-1.etikett]
		          deleted [tag-2.etikett]
		          deleted [tag-3.etikett]
		          deleted [tag-4.etikett]
		          deleted [tag.etikett]
	EOM

	run find . -maxdepth 2 ! -ipath './.zit*'
	assert_output '.'
}

function clean_zettels { # @test
	run_zit clean .z
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [one/dos.zettel]
		          deleted [one/uno.zettel]
		          deleted [one/]
	EOM

	run find . -maxdepth 2 ! -ipath './.zit*'
	assert_success
	assert_output_unsorted - <<-EOM
		.
		./md.typ
		./tag-1.etikett
		./tag-2.etikett
		./tag-3.etikett
		./tag-4.etikett
		./tag.etikett
	EOM
}

function clean_all_dirty_wd { # @test
	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		- etikett-one
		! md
		---

		newest body
	EOM

	cat >one/dos.zettel <<-EOM
		---
		# dos wildly different
		- etikett-two
		! md
		---

		dos newest body
	EOM

	cat >md.typ <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM

	cat >da-new.typ <<-EOM
		inline-akte = true
		vim-syntax-type = "da-new"
	EOM

	cat >zz-archive.etikett <<-EOM
		hide = true
	EOM

	run_zit clean .
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [tag-1.etikett]
		          deleted [tag-2.etikett]
		          deleted [tag-3.etikett]
		          deleted [tag-4.etikett]
		          deleted [tag.etikett]
	EOM

	run find . -maxdepth 2 ! -ipath './.zit*'
	assert_success
	assert_output_unsorted - <<-EOM
		.
		./md.typ
		./one
		./one/uno.zettel
		./one/dos.zettel
		./da-new.typ
		./zz-archive.etikett
	EOM
}

function clean_all_force_dirty_wd { # @test
	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		- etikett-one
		! md
		---

		newest body
	EOM

	cat >one/dos.zettel <<-EOM
		---
		# dos wildly different
		- etikett-two
		! md
		---

		dos newest body
	EOM

	cat >md.typ <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM

	cat >da-new.typ <<-EOM
		inline-akte = true
		vim-syntax-type = "da-new"
	EOM

	cat >zz-archive.etikett <<-EOM
		hide = true
	EOM

	run_zit clean -force .
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [da-new.typ]
		          deleted [md.typ]
		          deleted [one/dos.zettel]
		          deleted [one/uno.zettel]
		          deleted [one/]
		          deleted [tag-1.etikett]
		          deleted [tag-2.etikett]
		          deleted [tag-3.etikett]
		          deleted [tag-4.etikett]
		          deleted [tag.etikett]
		          deleted [zz-archive.etikett]
	EOM

	run find . -maxdepth 2 ! -ipath './.zit*'
	assert_success
	assert_output '.'
}

function clean_mode_akte { # @test
	skip
	run_zit organize -mode commit-directly :z <<-EOM
		- [one/uno  !md zz-archive tag-3 tag-4] wow the first
	EOM
	assert_success
	assert_output - <<-EOM
		[zz@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz-archive@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4 zz-archive]
	EOM

	run_zit schlummernd-add zz-archive
	assert_success
	assert_output ''

	run_zit checkout -mode akte one/uno
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4 zz-archive
		                   one/uno.md]
	EOM

	run_zit clean !md:z
	assert_success
	assert_output - <<-EOM
		          deleted [one/uno.md]
		          deleted [one/dos.zettel]
		          deleted [one/]
	EOM
}
