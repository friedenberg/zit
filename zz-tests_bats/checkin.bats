#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	run_zit checkout :z,t,e
	assert_success

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

	cat >zz-archive.etikett <<-EOM
		hide = true
	EOM
}

teardown() {
	rm_from_version "$version"
}

function checkin_simple_one_zettel { # @test
	run_zit checkin one/uno.zettel
	assert_success
	assert_output - <<-EOM
		[-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett-one@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function checkin_complex_zettel_etikett_negation { # @test
	run_zit checkin ^-etikett-two.z
	assert_success
	assert_output - <<-EOM
		[-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett-one@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function checkin_simple_all { # @test
	# TODO: modify this to support "." for everything
	run_zit checkin .z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
		[-etikett-one@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett-two@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-zz-archive@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a]
		[one/dos@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM

	run_zit show -format log :z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
		[-etikett-one@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett-two@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-zz-archive@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a]
		[one/dos@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function checkin_simple_all_dry_run { # @test
	# TODO: modify this to support "." for everything
	run_zit checkin -dry-run .z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
		[-etikett-one@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett-two@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-zz-archive@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a]
		[one/dos@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM

	run_zit show -format log :z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
		[-etikett-one@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett-two@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-zz-archive@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}

function checkin_simple_typ { # @test
	run_zit checkin .t
	assert_success
	assert_output - <<-EOM
		[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
	EOM

	run_zit show -format log !md.t
	assert_success
	assert_output - <<-EOM
		[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
	EOM

	run_zit show -format vim-syntax-type !md.typ
	assert_success
	assert_output 'test'
}

function checkin_simple_etikett { # @test
	run_zit checkin zz-archive.e
	assert_success
	assert_output - <<-EOM
		[-zz-archive@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a]
	EOM

	run_zit show -format text zz-archive.e
	assert_success
	assert_output 'hide = true'
}
