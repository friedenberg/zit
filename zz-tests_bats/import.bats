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

function import { # @test
	run_zit show -format bestandsaufnahme +:z,e,t,k
	assert_success
	echo -n "$output" | gzip >besties

	besties="$(realpath besties)"
	akten="$(realpath .zit/Objekten2/Akten)"

	wd1="$(mktemp -d)"
	cd "$wd1" || exit 1

	run_zit_init

	run_zit import -bestandsaufnahme "$besties" -akten "$akten"
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[-tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/uno@3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md "wow ok" tag-1 tag-2]
		copied Akte 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 (10 bytes)
		copied Akte 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 (16 bytes)
		copied Akte 3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 (27 bytes)
	EOM

	run_zit show one/uno
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
}

function import_age { # @test
	skip
	run_zit show -format bestandsaufnahme +:z,e,t,k
	assert_success
	echo -n "$output" >besties

	besties="$(realpath besties)"
	akten="$(realpath .zit/Objekten2/Akten)"
	age_id="$(realpath .zit/AgeIdentity)"

	wd1="$(mktemp -d)"
	cd "$wd1" || exit 1

	run_zit_init

	run_zit import -bestandsaufnahme "$besties" -akten "$akten" -age-identity "$(cat "$age_id")"
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[-tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/uno@3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md "wow ok" tag-1 tag-2]
		copied Akte 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 (10 bytes)
		copied Akte 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 (16 bytes)
		copied Akte 3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 (27 bytes)
	EOM

	run_zit show one/uno
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
}