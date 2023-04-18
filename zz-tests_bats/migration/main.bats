#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	cp -r "$DIR/$version" "$BATS_TEST_TMPDIR"
	cd "$BATS_TEST_TMPDIR/$version" || exit 1
}

teardown() {
	rm_from_version "$version"
}

function init_and_deinit { # @test
	run_zit status
	assert_success
}

function validate_contents { # @test
	run_zit show -format log +z,e,t
	assert_output_unsorted - <<-EOM
		[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		[one/uno@797cbdf8448a2ea167534e762a5025f5a3e9857e1dd06a3b746d3819d922f5ce !md "wow ok"]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}

function reindex { # @test
	run_zit reindex
	assert_output_unsorted - <<-EOM
		[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[konfig@7a09788554068a2e1012fe0fbd152bb8d24cd95e15407af4b28e753f151e6534]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		[one/uno@797cbdf8448a2ea167534e762a5025f5a3e9857e1dd06a3b746d3819d922f5ce !md "wow ok"]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}
