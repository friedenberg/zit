#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
	load 'common.bash'
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

cmd_def_organize=(
	"${cmd_zit_def[@]}"
	-right-align=false
	-prefix-joints=true
	-metadatei-header=false
	-refine=true
)

function outputs_organize_one_etikett { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/uno@9a638e2b183562da6d3c634d5a3841d64bc337c9cf79f8fffa0d0194659bc564 !md "wow"]'

	run_zit show one/uno
	assert_output "$(cat "$to_add")"

	run_zit show ok
	assert_output "$(cat "$to_add")"

	run_zit expand-hinweis o/u
	assert_output 'one/uno'

	expected_organize="$(mktemp)"
	{
		echo
		echo "# ok"
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize[@]}" -mode output-only ok
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
		echo "- brown"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/uno@f0be3e8072724eee5ea5022db397e20deb739d151abef61d37ed386207e32092 !md "wow"]'

	expected_organize="$(mktemp)"
	{
		echo
		echo "# brown, ok"
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run zit show "${cmd_zit_def[@]}" ok brown
	assert_output "$(cat "$to_add")"

	run_zit organize "${cmd_def_organize[@]}" -mode output-only ok brown
	assert_output "$(cat "$expected_organize")"

	{
		echo
		echo "# ok"
		echo
		echo "- [one/uno] wow"
		echo
	} >"$expected_organize"

	expected_organize_output="$(mktemp)"
	{
		echo "Removed etikett 'brown' from zettel 'one/uno'"
		echo '[one/uno@9a638e2b183562da6d3c634d5a3841d64bc337c9cf79f8fffa0d0194659bc564 !md "wow"]'
	} >"$expected_organize_output"

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly ok brown <"$expected_organize"
	assert_output "$(cat "$expected_organize_output")"

	expected_zettel="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$expected_zettel"

	run zit show brown
	assert_output ""

	run zit show ok
	assert_output "$(cat "$expected_zettel")"

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
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/uno@2df585d527ed7e18b3a9346079335509272f5a197b6a2d864e1b80df5ba627bf !md "wow"]'

	expected_organize="$(mktemp)"
	{
		echo
		echo "# task"
		echo
		echo " ## priority"
		echo
		echo "  ### -1"
		echo
		echo "  - [one/uno] wow"
		echo
		echo "  ### -2"
		echo
		echo "  - [one/uno] wow"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize[@]}" -mode output-only -group-by priority task
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
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/uno@b28b69e2e325ca2c7d0144a5d4db6523c2f241958229678ac39a9c5a200386bc !md "one/uno"]'

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-2"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/dos@2720ade68463c806a1aca98df4325e1904a357c6194bf3a8bc981091890aaeed !md "two/dos"]'

	expected_organize="$(mktemp)"
	{
		echo
		echo "# task"
		echo
		echo " ## priority"
		echo
		echo "  ### -1"
		echo
		echo "  - [one/uno] one/uno"
		echo
		echo "  ### -2"
		echo
		echo "  - [one/dos] two/dos"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize[@]}" -mode output-only -group-by priority task
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
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"

	expected_organize="$(mktemp)"
	{
		echo
		echo "# task"
		echo
		echo " ## priority-1"
		echo
		echo "  ### w-2022-07"
		echo
		echo "   #### -06"
		echo
		echo "   - [one/dos] two/dos"
		echo
		echo "   #### -07"
		echo
		echo "   - [one/uno] one/uno"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize[@]}" -mode output-only -group-by priority,w task
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
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/uno@112894f9e6c0b4eb6d39f70482312303870c85123f393d4ebb5a6b1118980d39 !md "one/uno"]'

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_output '[one/dos@1fe2b8f15cd9ec231a5d82a5f2317bfa090ec46e8d879e623083caaac28d46aa !md "two/dos"]'

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

	run_zit new -edit=false "$to_add"
	assert_output '[two/uno@168dbe89748356f7a3d229cab256a82e94106541e0af94a8695bf17f7a661241 !md "3"]'

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

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
	assert_success

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
	assert_success

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-07"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success

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
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success

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
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success

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

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
	assert_success

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
	assert_success

	run zit show two/dos
	assert_success

	run zit show three/uno
	assert_success

	{
		echo priority-1
		echo priority-2
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	#TODO
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

	run_zit new -edit=false "$one"

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

	run_zit new -edit=false "$two"

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

	run_zit new -edit=false "$three"

	expected_organize="$(mktemp)"
	{
		echo "# task"
		echo
		echo "## priority-1"
		echo
		echo "### w-2022-07-06"
		echo
		echo "- [one/dos] two/dos"
		echo "- [two/uno] 3"
		echo
		echo "### w-2022-07-07"
		echo
		echo "- [one/uno] one/uno"
		echo
		echo "###"
		echo
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
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

	one="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-07"
		echo "---"
	} >"$one"

	run_zit new -edit=false "$one"

	two="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "---"
	} >"$two"

	run_zit new -edit=false "$two"

	three="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "---"
	} >"$three"

	run_zit new -edit=false "$three"

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

	run_zit organize "${cmd_def_organize[@]}" -verbose -mode commit-directly -group-by priority,w task <"$expected_organize"

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

	run_zit new -edit=false "$one"

	expected_organize="$(mktemp)"
	{
		echo
		echo "# inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2"
		echo
		echo "- [one/uno] jabra coral usb_a-to-usb_c cable"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize[@]}" -mode output-only -group-by inventory \
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

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly <"$first_organize"

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
		echo "---"
	} >"one/uno.zettel"

	run zit checkin "${cmd_zit_def[@]}" one/uno.zettel
	#TODO-P1 fix typ
	assert_output '[one/uno@dc8d9d8e200a9c2f75e375bdfa267f30605f89229d4be70796800479c01ceede ! test4]'

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
		echo "---"
	} >"one/uno.zettel"

	run zit checkin "${cmd_zit_def[@]}" one/uno.zettel
	assert_output '[one/uno@9153182e2be5871aba88bb75f5a317e3f8dd73f8b2040bca4ac446679d17ef18 ! test1-ok, test4]'

	expected_etiketten="$(mktemp)"
	{
		echo test1-ok
		echo test4
	} >"$expected_etiketten"

	run zit cat -gattung etikett
	assert_output "$(cat "$expected_etiketten")"
}