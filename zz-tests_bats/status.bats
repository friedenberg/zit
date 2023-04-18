#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	run_zit checkout :z,t,e
	assert_success
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
		[!txt@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-etikett-one@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[two/uno@2e844ebe1018e2071c6f2b6b37a9ea2c1bd69e391d89f54aa4256228a1d49db0 !txt "the new zettel"]
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
		             same [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM

	dirty_one_uno

	run_zit status one/uno.zettel
	assert_success
	assert_output - <<-EOM
		          changed [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function status_zettel_akte_checkout { # @test
	run_zit clean .
	assert_success

	dirty_new_zettel

	run_zit checkout -mode akte two/uno
	assert_success
	assert_output - <<-EOM
		      checked out [two/uno.txt@aeb82efa111ccb5b8c5ca351f12d8b2f8e76d8d7bd0ecebf2efaaa1581d19400 !txt "the new zettel"]
	EOM

	run_zit status .z
	assert_success
	assert_output - <<-EOM
		             same [two/uno@2e844ebe1018e2071c6f2b6b37a9ea2c1bd69e391d89f54aa4256228a1d49db0 !txt "the new zettel"]
		                ↳ [two/uno.txt@aeb82efa111ccb5b8c5ca351f12d8b2f8e76d8d7bd0ecebf2efaaa1581d19400]
	EOM
}

function status_zettelen_typ { # @test
	run_zit status !md.z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		             same [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM

	dirty_one_uno
	dirty_one_dos

	run_zit status !md.z
	assert_success
	assert_output_unsorted - <<-EOM
		          changed [one/dos.zettel@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		          changed [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function status_complex_zettel_etikett_negation { # @test
	run_zit status ^-etikett-two.z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		             same [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	EOM

	dirty_one_uno

	run_zit status ^-etikett-two.z
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		          changed [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function status_simple_all { # @test
	run_zit status .
	assert_success
	#TODO why is fix issue with untracked appear
	assert_output_unsorted - <<-EOM
		             same [md.typ@b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15 !md]
		             same [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		             same [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		             same [tag-1.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-1]
		             same [tag-2.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-2]
		             same [tag-3.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-3]
		             same [tag-4.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-4]
		             same [tag.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag]
	EOM

	dirty_one_uno
	dirty_one_dos
	dirty_md_typ
	dirty_zz_archive_etikett
	dirty_da_new_typ

	run_zit status .
	assert_success
	assert_output_unsorted - <<-EOM
		             same [tag-1.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-1]
		             same [tag-2.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-2]
		             same [tag-3.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-3]
		             same [tag-4.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-4]
		             same [tag.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag]
		          changed [md.typ@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa !md]
		          changed [one/dos.zettel@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		          changed [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
		        untracked [da-new.typ@97c5ec233de522a564d1a6d43ea992b4f92d2bfa439762215b5da85780a4f529 !da-new]
		        untracked [zz-archive.etikett@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a -zz-archive]
	EOM
}

function status_simple_typ { # @test
	run_zit status .t
	assert_success
	assert_output_unsorted - <<-EOM
		             same [md.typ@b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15 !md]
	EOM

	dirty_md_typ
	dirty_da_new_typ

	run_zit status .t
	assert_success
	assert_output_unsorted - <<-EOM
		          changed [md.typ@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa !md]
		        untracked [da-new.typ@97c5ec233de522a564d1a6d43ea992b4f92d2bfa439762215b5da85780a4f529 !da-new]
	EOM
}

function status_simple_etikett { # @test
	run_zit status .e
	assert_success
	assert_output_unsorted - <<-EOM
		             same [tag-1.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-1]
		             same [tag-2.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-2]
		             same [tag-3.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-3]
		             same [tag-4.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-4]
		             same [tag.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag]
	EOM

	dirty_zz_archive_etikett

	run_zit status .e
	assert_success
	assert_output_unsorted - <<-EOM
		             same [tag-1.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-1]
		             same [tag-2.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-2]
		             same [tag-3.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-3]
		             same [tag-4.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag-4]
		             same [tag.etikett@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249 -tag]
		        untracked [zz-archive.etikett@0b7afc0b23d2f265b64bc184728d540cbadd0df54a2ae719e9757bcf17d8548a -zz-archive]
	EOM
}
