#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

function bootstrap {
	run_zit_init

	{
		echo "---"
		echo "# wow"
		echo "- tag"
		echo "! md"
		echo "---"
		echo
		echo "body"
	} >to_add

	run_zit new -edit=false to_add
	assert_output '[one/uno@37d3869e9b1711f009eabf69a2bf294cfd785f5b1c7463cba77d11d5f81f5e09 !md "wow"]'

	run zit show one/uno
	assert_output "$(cat to_add)"
}

function clone { # @test
	wd1="$(mktemp -d)"
	cd "$wd1" || exit 1
	bootstrap "$wd1"
	assert_success

	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit clone \
		"$wd1" :+zettel,typ

	assert_success
	assert_output --partial '[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]'
	assert_output --partial '[konfig@7633ddede5e2acbc3fb490417c0504790a94ea8cd726e4f9156115a8fa61173e]'
	assert_output --partial '[one/uno@37d3869e9b1711f009eabf69a2bf294cfd785f5b1c7463cba77d11d5f81f5e09 !md "wow"]'
}
