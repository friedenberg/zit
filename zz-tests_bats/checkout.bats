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

function checkout_simple_all { # @test
	run_zit checkout @z,t,e
	assert_output_unsorted - <<-EOM
		              (same) [md.typ@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7 !md]
		       (checked out) [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		       (checked out) [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}

function checkout_simple_zettel { # @test
	run_zit checkout @
	assert_output_unsorted - <<-EOM
		       (checked out) [one/dos.zettel@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
		       (checked out) [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM
}

function checkout_simple_typ { # @test
	run_zit checkout @t
	assert_output '              (same) [md.typ@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7 !md]'
}
