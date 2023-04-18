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
		[-zz-archive@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a]
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
		[-project@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-project-2021@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-project-2021-zit@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-zz@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-zz-archive-task@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-zz-archive-task-done@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
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
