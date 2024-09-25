#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output
}

function bootstrap {
	run_zit_init

	{
		echo "---"
		echo "# wow"
		echo "- tag"
		echo "! md"
		echo "---"
		echo
		echo "body"
	} >to_add

	run_zit new -edit=false to_add
	assert_success
	assert_output - <<-EOM
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
	EOM

	run_zit show one/uno
	assert_success
	assert_output "$(cat to_add)"
}

function clone { # @test
	skip
	wd1="$(mktemp -d)"
	cd "$wd1" || exit 1
	bootstrap "$wd1"
	assert_success

	run_zit new -edit=false - <<-EOM
		---
		# zettel with multiple etiketten
		- this_is_the_first
		- this_is_the_second
		! md
		---

		zettel with multiple etiketten body
	EOM

	assert_success
	assert_output - <<-EOM
		[-this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
	EOM

	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit clone "$wd1" +zettel,typ

	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[konfig @c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685]
		[konfig @c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685]
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
	EOM

	run_zit show one/dos
	assert_success
	assert_output - <<-EOM
		---
		# zettel with multiple etiketten
		- this_is_the_first
		- this_is_the_second
		! md
		---

		zettel with multiple etiketten body
	EOM
}
