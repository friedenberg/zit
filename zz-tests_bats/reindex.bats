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

function status_simple_one_zettel { # @test
	run_zit reindex
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]
		[konfig@7a09788554068a2e1012fe0fbd152bb8d24cd95e15407af4b28e753f151e6534]
		[one/uno@797cbdf8448a2ea167534e762a5025f5a3e9857e1dd06a3b746d3819d922f5ce !md "wow ok"]
		[-tag-2@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-tag-1@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-tag@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		[-tag-4@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-tag-3@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}