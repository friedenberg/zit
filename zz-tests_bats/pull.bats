#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output
}

teardown() {
	chflags_and_rm
}

# bats file_tags=user_story:pull,user_story:repo,user_store:xdg

function bootstrap {
	set_xdg "$1"
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
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
	EOM

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
		[this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
	EOM
}

function try_add_new_after_pull {
	run_zit new -edit=false - <<-EOM
		---
		# zettel after clone description
		! md
		---

		zettel after clone body
	EOM

	assert_success
	assert_output - <<-EOM
		[two/uno @13af191e86dcd8448565157de81919f19337656787f3d0fdd90b5335d2170f3f !md "zettel after clone description"]
	EOM
}

function pull_history_zettel_typ_etikett_no_conflicts { # @test
	them="them"
	bootstrap "$them"
	assert_success

	function print_their_xdg() (
		set_xdg "$them"
		zit info xdg
	)

	set_xdg "$BATS_TEST_TMPDIR"

	run_zit_init_disable_age
	run_zit pull \
		-xdg-dotenv <(print_their_xdg) \
		+zettel,typ,etikett

	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		copied Blob 024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b (36 bytes)
		copied Blob 9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 (5 bytes)
	EOM

	try_add_new_after_pull
}

function pull_history_default { # @test
	them="them"
	bootstrap "$them"
	assert_success

	function print_their_xdg() (
		set_xdg "$them"
		zit info xdg
	)

	set_xdg "$BATS_TEST_TMPDIR"

	run_zit_init_disable_age
	run_zit pull -xdg-dotenv <(print_their_xdg)

	assert_success
	assert_output_unsorted --regexp - <<-EOM
		\[.+\..+ @.+ !inventory_list-v1]
		\[.+\..+ @.+ !inventory_list-v1]
		\[.+\..+ @.+ !inventory_list-v1]
		\[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
		\[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
		\[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		\[this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		\[this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		copied Blob .+ \(.+ bytes)
		copied Blob .+ \(.+ bytes)
		copied Blob .+ \(.+ bytes)
		copied Blob .+ \(.+ bytes)
		copied Blob .+ \(.+ bytes)
	EOM

	try_add_new_after_pull
}
