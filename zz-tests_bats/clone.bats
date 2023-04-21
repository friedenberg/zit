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
	assert_success
	assert_output - <<-EOM
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[one/uno@37d3869e9b1711f009eabf69a2bf294cfd785f5b1c7463cba77d11d5f81f5e09 !md "wow"]
	EOM

	run_zit show one/uno
	assert_success
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
		"$wd1" +zettel,typ

	assert_success
	assert_output_unsorted - <<-EOM
		[!md@b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15]
		[konfig@f6d3d0874fd9475c2b7ac150f366cd211d847a8676ccabc35111cb357fd0c3b9]
		[one/uno@37d3869e9b1711f009eabf69a2bf294cfd785f5b1c7463cba77d11d5f81f5e09 !md "wow"]
	EOM
}
