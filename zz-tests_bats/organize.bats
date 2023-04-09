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

function hides_hidden_etiketten_from_organize { # @test
	echo "hide = true" >zz-archive.etikett
	run_zit checkin -delete .e
	assert_success
	assert_output - <<-EOM
		[-zz-archive@cba019d4f889027a3485e56dd2080c7ba0fa1e27499c24b7ec08ad80ef55da9d]
		          deleted [zz-archive.etikett]
	EOM

	to_add="$(mktemp)"
	{
		echo ---
		echo "# split hinweis for usability"
		echo - project-2021-zit
		echo - zz-archive-task-done
		echo ! md
		echo ---
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success
	assert_output - <<-EOM
		[-project@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-project-2021@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-project-2021-zit@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-archive-task@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-zz-archive-task-done@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[two/uno@8a35da296f0c4b007f386ca26553e2fc92c71173cf372e575b3cda857f7fb7e3 !md "split hinweis for usability"]
	EOM

	expected_organize="$(mktemp)"
	{
		echo
		echo "# project-2021-zit"
		echo
	} >"$expected_organize"

	run_zit organize -mode output-only project-2021-zit
	assert_success
	assert_output - <<-EOM
		---
		- project-2021-zit
		---
	EOM
}
