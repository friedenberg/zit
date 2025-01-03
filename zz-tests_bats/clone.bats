#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output
	export BATS_TEST_BODY=true
}

teardown() {
	chflags_and_rm
}

# bats file_tags=user_story:clone,user_story:repo,user_store:xdg,user_story:remote

function bootstrap {
	mkdir -p "$1"
	pushd "$1" || exit 1
	run_zit_init -override-xdg-with-cwd

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

	popd || exit 1
}

function print_their_xdg() {
	pushd "$1" >/dev/null || exit 1
	zit info xdg
}

function run_clone_default_with() {
	run_zit clone \
		-yin <(cat_yin) \
		-yang <(cat_yang) \
		-lock-internal-files=false \
		"$@"
}

function try_add_new_after_clone {
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

function clone_history_zettel_typ_etikett { # @test
	them="them"
	bootstrap "$them"
	assert_success

	us="us"
	set_xdg "$us"
	run_clone_default_with \
		-remote-type native-dotenv-xdg \
		<(print_their_xdg them) \
		+zettel,typ,etikett

	assert_success
	assert_output_unsorted - <<-EOM
		[!md @b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 !toml-type-v1]
		[konfig @ef6b910d71068d5cb0598abeaea21140b44da67a4a4a9eca6485e8a2906ca483 !toml-config-v1]
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		copied Blob 024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b (36 bytes)
		copied Blob 9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 (5 bytes)
		copied Blob b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 (51 bytes)
	EOM

	try_add_new_after_clone
}

function clone_history_zettel_typ_etikett_stdio_local { # @test
	them="them"
	bootstrap "$them"
	assert_success

	us="us"
	set_xdg "$us"
	run_clone_default_with \
		-remote-type stdio-local \
		"$(realpath them)" \
		+zettel,typ,etikett

	assert_success
	assert_output_unsorted - <<-EOM
		[!md @b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 !toml-type-v1]
		[konfig @ef6b910d71068d5cb0598abeaea21140b44da67a4a4a9eca6485e8a2906ca483 !toml-config-v1]
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		copied Blob 024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b (36 bytes)
		copied Blob 9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 (5 bytes)
		copied Blob b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 (51 bytes)
	EOM

	try_add_new_after_clone
}

function clone_history_zettel_typ_etikett_stdio_ssh { # @test
	skip
	them="them"
	bootstrap "$them"
	assert_success

	us="us"
	set_xdg "$us"
	run_clone_default_with \
		-remote-type stdio-local \
		"$(realpath them)" \
		+zettel,typ,etikett

	assert_success
	assert_output_unsorted - <<-EOM
		[!md @b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 !toml-type-v1]
		[konfig @ef6b910d71068d5cb0598abeaea21140b44da67a4a4a9eca6485e8a2906ca483 !toml-config-v1]
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		copied Blob 024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b (36 bytes)
		copied Blob 9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 (5 bytes)
		copied Blob b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 (51 bytes)
	EOM

	try_add_new_after_clone
}

function clone_history_default_allow_conflicts { # @test
	them="them"
	bootstrap "$them"
	assert_success

	us="us"
	set_xdg "$us"
	run_clone_default_with \
		-remote-type native-dotenv-xdg \
		<(print_their_xdg them)

	assert_success

	run_zit show +?z,t,e
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 !toml-type-v1]
		[one/dos @024948601ce44cc9ab070b555da4e992f111353b7a9f5569240005639795297b !md "zettel with multiple etiketten" this_is_the_first this_is_the_second]
		[one/uno @9e2ec912af5dff2a72300863864fc4da04e81999339d9fac5c7590ba8a3f4e11 !md "wow" tag]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_first @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[this_is_the_second @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM

	try_add_new_after_clone
}
