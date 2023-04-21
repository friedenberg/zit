#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

function init_and_deinit { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	[[ -f .zit/KonfigAngeboren ]]

	run_zit deinit
	assert_success
}

# function init_and_init { ## @test
# 	wd="$(mktemp -d)"
# 	cd "$wd" || exit 1

# 	run_zit_init
# 	assert_success

# 	{
# 		echo "---"
# 		echo "# wow"
# 		echo "- tag"
# 		echo "! md"
# 		echo "---"
# 		echo
# 		echo "body"
# 	} >to_add

# 	run_zit new -edit=false to_add
# 	assert_success
# 	assert_output - <<-EOM
# 		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
# 		[one/uno@37d3869e9b1711f009eabf69a2bf294cfd785f5b1c7463cba77d11d5f81f5e09 !md "wow"]
# 	EOM

# 	run_zit show one/uno
# 	assert_success
# 	assert_output "$(cat to_add)"

# 	run_zit init -yin <(cat_yin) -yang <(cat_yang)
# 	assert_failure
# 	assert_output --partial '.zit/Kennung/Yin: file exists'

# 	run_zit init
# 	assert_success
# 	assert_output --partial '.zit/KonfigAngeboren already exists, not overwriting'
# 	assert_output --partial '.zit/KonfigErworben already exists, not overwriting'

# 	run zit show -format log :
# 	assert_success
# 	assert_output "$(cat to_add)"
# }
