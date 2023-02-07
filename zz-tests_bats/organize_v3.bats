#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
	# ... the remaining setup is unchanged

	# get the containing directory of this file
	# use $BATS_TEST_FILENAME instead of ${BASH_SOURCE[0]} or $0,
	# as those will point to the bats executable's location or the preprocessed file respectively
	DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" >/dev/null 2>&1 && pwd)"
	# make executables in src/ visible to PATH
	PATH="$DIR/../:$PATH"
	PATH="$DIR/../build/:$PATH"

	# for shellcheck SC2154
	export output
}

cat_yin() (
	echo "one"
	echo "two"
	echo "three"
)

cat_yang() (
	echo "uno"
	echo "dos"
	echo "tres"
)

cmd_zit_def=(
	# -abbreviate-hinweisen=false
	-predictable-hinweisen
	-print-typen=false
)

cmd_zit_new=(
	zit
	new
	"${cmd_zit_def[@]}"
)

cmd_zit_organize=(
	zit
	organize
	-right-align=false
	-refine=true
	-metadatei-header=false
	"${cmd_zit_def[@]}"
)

cmd_zit_organize_v3=(
	zit
	organize
	"${cmd_zit_def[@]}"
	-prefix-joints=true
	-metadatei-header=false
	-refine=true
)

function outputs_organize_one_etikett { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"
	assert_output '          (new) [o/u@9 !md "wow"]'

	run zit expand-hinweis o/u
	assert_output 'one/uno'

	expected_organize="$(mktemp)"
	{
		echo
		echo "      # ok"
		echo
		echo "- [o/u] wow"
	} >"$expected_organize"

	run "${cmd_zit_organize_v3[@]}" -mode output-only ok
	assert_output "$(cat "$expected_organize")"
}

function outputs_organize_two_etiketten { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "- brown"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"
	assert_output '          (new) [o/u@f !md "wow"]'

	expected_organize="$(mktemp)"
	{
		echo
		echo "      # brown, ok"
		echo
		echo "- [o/u] wow"
	} >"$expected_organize"

	run "${cmd_zit_organize_v3[@]}" -mode output-only ok brown
	assert_output "$(cat "$expected_organize")"

	{
		echo
		echo "      # ok"
		echo
		echo "- [o/u] wow"
		echo
	} >"$expected_organize"

	run "${cmd_zit_organize_v3[@]}" -verbose -mode commit-directly ok brown <"$expected_organize"

	expected_zettel="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$expected_zettel"

	run zit show one/uno
	assert_output "$(cat "$expected_zettel")"
}

function outputs_organize_one_etiketten_group_by_one { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- task"
		echo "- priority-1"
		echo "- priority-2"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"
	assert_output '          (new) [o/u@2 !md "wow"]'

	expected_organize="$(mktemp)"
	{
		echo
		echo "      # task"
		echo
		echo "     ## priority"
		echo
		echo "    ###         -1"
		echo
		echo "- [o/u] wow"
		echo
		echo "    ###         -2"
		echo
		echo "- [o/u] wow"
	} >"$expected_organize"

	run "${cmd_zit_organize_v3[@]}" -mode output-only -group-by priority task
	assert_output "$(cat "$expected_organize")"
}

function outputs_organize_two_zettels_one_etiketten_group_by_one { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- task"
		echo "- priority-1"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"
	assert_output '          (new) [o/u@b !md "one/uno"]'

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-2"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"
	assert_output '          (new) [o/d@2 !md "two/dos"]'

	expected_organize="$(mktemp)"
	{
		echo
		echo "      # task"
		echo
		echo "     ## priority"
		echo
		echo "    ###         -1"
		echo
		echo "- [o/u] one/uno"
		echo
		echo "    ###         -2"
		echo
		echo "- [o/d] two/dos"
	} >"$expected_organize"

	run "${cmd_zit_organize_v3[@]}" -mode output-only -group-by priority task
	assert_output "$(cat "$expected_organize")"
}

function outputs_organize_one_etiketten_group_by_two { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	expected_organize="$(mktemp)"
	{
		echo
		echo "      # task"
		echo
		echo "     ## priority-1"
		echo
		echo "    ### w-2022-07"
		echo
		echo "   ####          -06"
		echo
		echo "- [o/d] two/dos"
		echo
		echo "   ####          -07"
		echo
		echo "- [o/u] one/uno"
	} >"$expected_organize"

	run "${cmd_zit_organize_v3[@]}" -mode output-only -group-by priority,w task
	assert_output "$(cat "$expected_organize")"
}

function commits_organize_one_etiketten_group_by_two { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	expected_organize="$(mktemp)"
	{
		echo "# task"
		echo
		echo "## priority-1"
		echo
		echo "### w-2022-07-06"
		echo
		echo "- [one/dos] two/dos"
		echo
		echo "## priority-2"
		echo
		echo "### w-2022-07-07"
		echo
		echo "- [one/uno] one/uno"
		echo
		echo "###"
		echo
		echo "- [two/uno] 3"
	} >"$expected_organize"

	run "${cmd_zit_organize[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
	echo "$output"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$to_add"

	run zit show one/uno
	assert_output "$(cat "$to_add")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-2"
		echo "- task"
		echo "! md"
		echo "---"
	} >"$to_add"

	run zit show two/uno
	assert_output "$(cat "$to_add")"
}

function commits_organize_one_etiketten_group_by_two_new_zettels { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	expected="$(mktemp)"
	{
		echo priority-1
		echo task
		echo w-2022-07-07
	} >"$expected"

	run zit cat -gattung etikett
	assert_output "$(cat "$expected")"

	{
		echo one/uno
	} >"$expected"

	# run zit cat -gattung hinweis
	# assert_output --partial "$(cat "$expected")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	{
		echo priority-1
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	run zit cat -gattung etikett
	assert_output "$(cat "$expected")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$to_add"

	run "${cmd_zit_new[@]}" -edit=false "$to_add"

	{
		echo priority-1
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	run zit cat -gattung etikett
	assert_output "$(cat "$expected")"

	expected_organize="$(mktemp)"
	{
		echo "# task"
		echo "- new zettel one"
		echo "## priority-1"
		echo "- new zettel two"
		echo "### w-2022-07-06"
		echo "- [one/dos] two/dos"
		echo "## priority-2"
		echo "### w-2022-07-07"
		echo "- [one/uno] one/uno"
		echo "###"
		echo "- new zettel three"
		echo "- [two/uno] 3"
	} >"$expected_organize"

	run "${cmd_zit_organize[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$to_add"

	run zit show one/uno
	assert_output "$(cat "$to_add")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-2"
		echo "- task"
		echo "! md"
		echo "---"
	} >"$to_add"

	run zit show two/uno
	assert_output "$(cat "$to_add")"

	run zit show one/tres
	run zit show two/dos
	run zit show three/uno

	{
		echo priority-1
		echo priority-2
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	# TODO
	# run zit cat -gattung etikett
	# assert_output "$(cat "$expected")"
}

function commits_no_changes { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	one="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$one"

	run "${cmd_zit_new[@]}" -edit=false "$one"
	assert_output '          (new) [o/u@11 !md "one/uno"]'

	two="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$two"

	run "${cmd_zit_new[@]}" -edit=false "$two"
	assert_output '          (new) [o/d@1f !md "two/dos"]'

	three="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$three"

	run "${cmd_zit_new[@]}" -edit=false "$three"
	assert_output '          (new) [t/u@16 !md "3"]'

	expected_organize="$(mktemp)"
	{
		echo
		echo "# task"
		echo
		echo " ## priority-1"
		echo
		echo "  ### w-2022-07-06"
		echo
		echo "  - [t/u] 3"
		echo "  - [o/d] two/dos"
		echo
		echo "  ### w-2022-07-07"
		echo
		echo "  - [o/u] one/uno"
		echo
	} >"$expected_organize"

	# run "${cmd_zit_organize[@]}" -prefix-joints=false -mode output-only -group-by priority,w task
	run "${cmd_zit_organize[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
	# assert_output "$(cat "$expected_organize")"
	assert_output "no changes"

	run zit show one/uno
	assert_output "$(cat "$one")"

	run zit show one/dos
	assert_output "$(cat "$two")"

	run zit show two/uno
	assert_output "$(cat "$three")"
}

function commits_dependent_leaf { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	one="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$one"

	run "${cmd_zit_new[@]}" -edit=false "$one"
	assert_success

	two="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$two"

	run "${cmd_zit_new[@]}" -edit=false "$two"
	assert_success

	three="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "! md"
		echo "---"
	} >"$three"

	run "${cmd_zit_new[@]}" -edit=false "$three"
	assert_success

	expected_organize="$(mktemp)"
	{
		echo "# task"
		echo "## priority-2"
		echo "### w-2022-07"
		echo "#### -07"
		echo "- [one/dos] two/dos"
		echo "- [two/uno] 3"
		echo "#### -08"
		echo "- [one/uno] one/uno"
		echo "###"
	} >"$expected_organize"

	run "${cmd_zit_organize[@]}" -verbose -mode commit-directly -group-by priority,w task <"$expected_organize"
	assert_success

	one="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-08"
		echo "! md"
		echo "---"
	} >"$one"

	run zit show one/uno
	assert_output "$(cat "$one")"

	two="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$two"

	run zit show one/dos
	assert_output "$(cat "$two")"

	three="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-07"
		echo "! md"
		echo "---"
	} >"$three"

	run zit show two/uno
	assert_output "$(cat "$three")"
}

function zettels_in_correct_places { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	one="$(mktemp)"
	{
		echo "---"
		echo "# jabra coral usb_a-to-usb_c cable"
		echo "- inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2"
		echo "---"
	} >"$one"

	run "${cmd_zit_new[@]}" -edit=false "$one"

	expected_organize="$(mktemp)"
	{
		echo
		echo "# inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2"
		echo
		echo "- [one/uno] jabra coral usb_a-to-usb_c cable"
	} >"$expected_organize"

	run "${cmd_zit_organize[@]}" -mode output-only -group-by inventory \
		inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2

	assert_output "$(cat "$expected_organize")"
}

function etiketten_correct { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	first_organize="$(mktemp)"
	{
		echo
		echo "# test1"
		echo "## -wow"
		echo
		echo "- zettel bez"
	} >"$first_organize"

	run "${cmd_zit_organize[@]}" -mode commit-directly <"$first_organize"

	expected_etiketten="$(mktemp)"
	{
		echo test1-wow
	} >"$expected_etiketten"

	run zit cat -gattung etikett
	assert_output "$(cat "$expected_etiketten")"

	mkdir -p one
	{
		echo "---"
		echo "- test4"
		echo "! md"
		echo "---"
	} >"one/uno.zettel"

	run zit checkin "${cmd_zit_def[@]}" one/uno.zettel
	assert_output '      (updated) [o/u@4 !md test4]'

	expected_etiketten="$(mktemp)"
	{
		echo test4
	} >"$expected_etiketten"

	run zit cat -gattung etikett
	assert_output "$(cat "$expected_etiketten")"

	mkdir -p one
	{
		echo "---"
		echo "- test4"
		echo "- test1-ok"
		echo "! md"
		echo "---"
	} >"one/uno.zettel"

	run zit checkin "${cmd_zit_def[@]}" one/uno.zettel
	assert_output '      (updated) [o/u@1a !md test1-ok, test4]'

	expected_etiketten="$(mktemp)"
	{
		echo test1-ok
		echo test4
	} >"$expected_etiketten"

	run zit cat -gattung etikett
	assert_output "$(cat "$expected_etiketten")"
}
