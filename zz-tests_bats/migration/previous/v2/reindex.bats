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

function reindex_simple { # @test
	run_zit reindex
	assert_success
	# assert_output_unsorted - <<-EOM
	# 	[!md@b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15]
	# 	[!md@b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15]
	# 	[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
	# 	[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
	# 	[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
	# 	[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
	# 	[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
	# 	[konfig@62c02b6f59e6de576a3fcc1b89db6e85b75c2ff7820df3049a5b12f9db86d1f5]
	# 	[konfig@62c02b6f59e6de576a3fcc1b89db6e85b75c2ff7820df3049a5b12f9db86d1f5]
	# 	[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	# 	[one/uno@797cbdf8448a2ea167534e762a5025f5a3e9857e1dd06a3b746d3819d922f5ce !md "wow ok"]
	# 	[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	# EOM
}

function reindex_after_changes { # @test
	cat >md.typ <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM

	run_zit checkin .t
	assert_success
	# assert_output - <<-EOM
	# 	[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
	# EOM
	# cat >zz-archive.etikett <<-EOM
	# 	hide = true
	# EOM

	run_zit reindex
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md@220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
		[-tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[konfig@c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685]
		[konfig@c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
		[one/uno@3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md "wow ok"]
	EOM

	run_zit show !md:t
	assert_success
	assert_output - <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM
}