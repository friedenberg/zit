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

function organize_simple { # @test
	actual="$(mktemp)"
	run_zit organize -mode output-only :z,e,t >"$actual"
	assert_success
	assert_output_unsorted - <<-EOM

		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos  ] wow ok again
		- [one/uno  ] wow the first
	EOM
}

function organize_simple_commit { # @test
	run_zit organize -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos  ] wow ok again
		- [one/uno  ] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@d6f48deae132aa17acd1cea0dffbbf6f76776835fc5db48620b8e90e3ee10a33]
		[-new-etikett-for-all@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-new-etikett-for@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-new-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-new@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag-2@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag-3@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag-4@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[one/dos@591fd7a0cdad23ddac675cc66f3fe004080c2cdff64a43feb1ec0c02f2dae7a1 !md "wow ok again"]
		[one/uno@8fa484873bb584d5d8e8e0121d54d28a821080ceaf67399a3ba891ab82d9d54f !md "wow the first"]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@d6f48deae132aa17acd1cea0dffbbf6f76776835fc5db48620b8e90e3ee10a33 new-etikett-for-all]
		[-new-etikett-for-all@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag-2@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag-3@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag-4@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[-tag@55001f3c8d717cfc6b6a9f6620ecdc006f8ffd2fe440a740cc754f3238a57ebc]
		[one/dos@591fd7a0cdad23ddac675cc66f3fe004080c2cdff64a43feb1ec0c02f2dae7a1 !md "wow ok again"]
		[one/uno@8fa484873bb584d5d8e8e0121d54d28a821080ceaf67399a3ba891ab82d9d54f !md "wow the first"]
	EOM
}

function organize_hides_hidden_etiketten_from_organize { # @test
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

function organize_dry_run { # @test
	expected_show="$(mktemp)"
	# shellcheck disable=SC2154
	zit show "${cmd_zit_def[@]}" -format log :z,e,t >"$expected_show"

	run_zit organize -dry-run -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos  ] wow ok again
		- [one/uno  ] wow the first
	EOM
	assert_success

	run_zit show -format log :z,e,t
	assert_success
	assert_output_unsorted "$(cat "$expected_show")"
}
