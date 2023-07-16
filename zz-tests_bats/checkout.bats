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

function checkout_simple_all { # @test
	run_zit checkout :z,t,e
	assert_success
	assert_output_unsorted - <<-EOM
		             same [md.typ@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 !md]
		             same [tag-1.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -tag-1]
		             same [tag-2.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -tag-2]
		             same [tag-3.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -tag-3]
		             same [tag-4.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -tag-4]
		             same [tag.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 -tag]
		      checked out [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		      checked out [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}

function checkout_simple_zettel { # @test
	run_zit checkout :
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		      checked out [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}

function checkout_simple_zettel_akte_only { # @test
	run_zit clean .
	assert_success
	# TODO fail checkouts if working directly has incompatible checkout
	run_zit checkout -mode akte-only :z
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.md@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		      checked out [one/uno.md@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
	EOM
}

function checkout_zettel_several { # @test
	run_zit checkout one/uno one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		      checked out [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}

function checkout_simple_typ { # @test
	run_zit checkout :t
	assert_success
	assert_output_unsorted - <<-EOM
		             same [md.typ@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 !md]
	EOM
}
