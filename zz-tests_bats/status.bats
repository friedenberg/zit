#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	run_zit checkout :z,t,e
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [md.typ @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		      checked out [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		      checked out [tag-1.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag-2.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag-3.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag-4.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM
}

teardown() {
	rm_from_version "$version"
}

function dirty_new_zettel() {
	run_zit new -edit=false - <<-EOM
		---
		# the new zettel
		- etikett-one
		! txt
		---

		with a different typ
	EOM

	assert_success
	assert_output --partial - <<-EOM
		[!txt @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett-one @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @aeb82efa111ccb5b8c5ca351f12d8b2f8e76d8d7bd0ecebf2efaaa1581d19400 !txt "the new zettel" etikett-one]
	EOM
}

function dirty_existing_akte() {
	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		- etikett-one
		! one/uno.md
		---
	EOM

	cat >one/uno.md <<-EOM
		newest body but even newer
	EOM
}

function dirty_one_uno() {
	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		- etikett-one
		! md
		---

		newest body
	EOM
}

function dirty_one_dos() {
	cat >one/dos.zettel <<-EOM
		---
		# dos wildly different
		- etikett-two
		! md
		---

		dos newest body
	EOM
}

function dirty_md_typ() {
	cat >md.typ <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM
}

function dirty_da_new_typ() {
	cat >da-new.typ <<-EOM
		inline-akte = true
		vim-syntax-type = "da-new"
	EOM
}

function dirty_zz_archive_etikett() {
	cat >zz-archive.etikett <<-EOM
		hide = true
	EOM
}

function status_simple_one_zettel { # @test
	run_zit status one/uno.zettel
	assert_success
	assert_output - <<-EOM
		             same [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	dirty_one_uno

	run_zit status one/uno.zettel
	assert_success
	assert_output - <<-EOM
		          changed [one/uno.zettel @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}

function status_simple_one_zettel_akte_separate { # @test
	run_zit status one/uno.zettel
	assert_success
	assert_output - <<-EOM
		             same [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	rm one/uno.zettel

	cat >one/uno.md <<-EOM
		newest body but even newerests
	EOM

	run_zit status one/uno.zettel
	assert_success
	assert_output - <<-EOM
		          changed [one/uno @a958b1c8e2bc817fcb17292f6957c0dfc87c874dc33274f0c4f4efdcdd1429bb !md "wow the first" tag-3 tag-4
		                   one/uno.md]
	EOM
}

function status_simple_one_zettel_akte_only { # @test
	run_zit clean one/uno.zettel
	assert_success
	assert_output - <<-EOM
		          deleted [one/uno.zettel]
	EOM

	run_zit checkout -mode blob one/uno
	assert_output - <<-EOM
		      checked out [one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4
		                   one/uno.md]
	EOM

	run_zit status one/uno.zettel
	assert_success
	assert_output - <<-EOM
		             same [one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4
		                   one/uno.md]
	EOM

	dirty_existing_akte

	# TODO fix non-determinism
	run_zit status one/uno.zettel
	assert_success
	assert_output - <<-EOM
		          changed [one/uno.zettel @e5ef6f74b2707b17d8670e5678151d676655c685c43beaeb6e995c9d127fab85 !md "wildly different" etikett-one
		                   one/uno.md]
	EOM
}

function status_zettel_akte_checkout { # @test
	run_zit clean .
	assert_success

	dirty_new_zettel

	run_zit checkout -mode blob two/uno
	assert_success
	assert_output - <<-EOM
		      checked out [two/uno @aeb82efa111ccb5b8c5ca351f12d8b2f8e76d8d7bd0ecebf2efaaa1581d19400 !txt "the new zettel" etikett-one
		                   two/uno.txt]
	EOM

	run_zit status .z
	assert_success
	assert_output - <<-EOM
		             same [two/uno @aeb82efa111ccb5b8c5ca351f12d8b2f8e76d8d7bd0ecebf2efaaa1581d19400 !txt "the new zettel" etikett-one
		                   two/uno.txt]
	EOM
}

function status_zettel_hidden { # @test
	run_zit schlummernd-add tag-3
	assert_success

	run_zit show :z
	assert_success
	assert_output ''

	run_zit show :?z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit status .z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	run_zit status !md:z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM
}

function status_zettelen_typ { # @test
	run_zit status !md.z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		             same [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	dirty_one_uno
	dirty_one_dos

	run_zit status !md.z
	assert_success
	assert_output_unsorted - <<-EOM
		          changed [one/dos.zettel @b5c4fbaac3b71657edee74de4b947f13dfa104715feb8bab7cfa4dd47cafa3db !md "dos wildly different" etikett-two]
		          changed [one/uno.zettel @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}

function status_complex_zettel_etikett_negation { # @test
	run_zit status ^-etikett-two.z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		             same [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	dirty_one_uno

	run_zit status ^-etikett-two.z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		          changed [one/uno.zettel @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}

function status_simple_all { # @test
	run_zit status
	assert_success
	assert_output_unsorted - <<-EOM
		             same [md.typ @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		             same [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		             same [tag-1.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-2.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-3.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-4.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM

	dirty_one_uno
	dirty_one_dos
	dirty_md_typ
	dirty_zz_archive_etikett
	dirty_da_new_typ

	run_zit status .
	assert_success
	assert_output_unsorted - <<-EOM
		             same [tag-1.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-2.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-3.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-4.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		          changed [md.typ @220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
		          changed [one/dos.zettel @b5c4fbaac3b71657edee74de4b947f13dfa104715feb8bab7cfa4dd47cafa3db !md "dos wildly different" etikett-two]
		          changed [one/uno.zettel @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
		        untracked [da-new.typ @1a4c3a8914d9e5fa1a08cb183bcdf7e6e10aa224f663dc56610a303b10aa0834]
		        untracked [zz-archive.etikett @b8cd0eaa1891284eafdf99d3acc2007a3d4396e8a7282335f707d99825388a93]
	EOM
}

function status_simple_typ { # @test
	run_zit status .t
	assert_success
	assert_output_unsorted - <<-EOM
		             same [md.typ @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
	EOM

	dirty_md_typ
	dirty_da_new_typ

	run_zit status .t
	assert_success
	assert_output_unsorted - <<-EOM
		          changed [md.typ @220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
		        untracked [da-new.typ @1a4c3a8914d9e5fa1a08cb183bcdf7e6e10aa224f663dc56610a303b10aa0834]
	EOM
}

function status_simple_etikett { # @test
	run_zit status .e
	assert_success
	assert_output_unsorted - <<-EOM
		             same [tag-1.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-2.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-3.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-4.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM

	dirty_zz_archive_etikett

	run_zit status .e
	assert_success
	assert_output_unsorted - <<-EOM
		             same [tag-1.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-2.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-3.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag-4.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		             same [tag.etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		        untracked [zz-archive.etikett @b8cd0eaa1891284eafdf99d3acc2007a3d4396e8a7282335f707d99825388a93]
	EOM
}

function status_conflict { # @test
	run_zit checkout one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	cat - >one/dos.zettel <<-EOM
		---
		# wow ok again
		- get_this_shit_merged
		- tag-3
		- tag-4
		! txt
		---

		not another one, conflict time
	EOM

	run_zit organize -mode commit-directly one/dos <<-EOM
		---
		! txt2
		---

		# new-etikett-for-all
		- [one/dos  tag-3 tag-4] wow ok again
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!txt2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett-for @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !txt2 "wow ok again" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !txt2 "wow ok again" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit status one/dos.zettel
	assert_success
	assert_output - <<-EOM
		       conflicted [one/dos.zettel]
	EOM
}
