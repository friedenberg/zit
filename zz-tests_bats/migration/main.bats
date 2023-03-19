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
	chflags -R nouchg "$BATS_TEST_TMPDIR/$version"
}

function init_and_deinit { # @test
	skip
	run_zit status
	assert_success
}

function validate_contents { # @test
	skip
	run_zit show -format log +@z,e,t
	assert_output_unsorted - <<-EOM
		[one/uno@797cbdf8448a2ea167534e762a5025f5a3e9857e1dd06a3b746d3819d922f5ce !md "wow ok"]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]
	EOM
}
