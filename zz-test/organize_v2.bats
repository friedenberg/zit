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

function outputs_organize_one_etikett { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"
	assert_output --partial '[one/uno '

	expected_organize="$(mktemp)"
	{
		echo
		echo "# ok"
		echo
		echo "- [one/uno] wow"
	} >>"$expected_organize"

	run zit organize -mode output-only ok
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
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"
	assert_output --partial '[one/uno '

	expected_organize="$(mktemp)"
	{
		echo
		echo "# brown, ok"
		echo
		echo "- [one/uno] wow"
	} >>"$expected_organize"

	run zit organize -mode output-only ok brown
	assert_output "$(cat "$expected_organize")"

	{
		echo
		echo "# ok"
		echo
		echo "- [one/uno] wow"
		echo
	} >"$expected_organize"

	run zit organize -verbose -mode commit-directly ok brown <"$expected_organize"

	expected_zettel="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$expected_zettel"

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
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"
	assert_output --partial '[one/uno '

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
	} >>"$expected_organize"

	run zit organize -mode output-only -group-by priority task
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
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"
	assert_output --partial '[one/uno '

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-2"
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"
	assert_output --partial '[one/dos '

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
	} >>"$expected_organize"

	run zit organize -mode output-only -group-by priority task
	assert_output "$(cat "$expected_organize")"
}

function outputs_organize_one_etiketten_group_by_two { # @test
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
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

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
	} >>"$expected_organize"

	run zit organize -mode output-only -group-by priority,w task
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
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

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
	} >>"$expected_organize"

	run zit organize -predictable-hinweisen -mode commit-directly -group-by priority,w task <"$expected_organize"
	echo "$output"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-07"
		echo "---"
	} >>"$to_add"

	run zit show one/uno
	assert_output "$(cat "$to_add")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-2"
		echo "- task"
		echo "---"
	} >>"$to_add"

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
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

	expected="$(mktemp)"
	{
		echo priority-1
		echo task
		echo w-2022-07-07
	} >"$expected"

	run zit cat -type etikett
	assert_output "$(cat "$expected")"

	{
		echo one/uno
	} >"$expected"

	# run zit cat -type hinweis
	# assert_output --partial "$(cat "$expected")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

	{
		echo priority-1
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	run zit cat -type etikett
	assert_output "$(cat "$expected")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- task"
		echo "- priority-1"
		echo "- w-2022-07-06"
		echo "---"
	} >>"$to_add"

	run zit new -edit=false -predictable-hinweisen "$to_add"

	{
		echo priority-1
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	run zit cat -type etikett
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
	} >>"$expected_organize"

	run zit organize -predictable-hinweisen -mode commit-directly -group-by priority,w task <"$expected_organize"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-07"
		echo "---"
	} >>"$to_add"

	run zit show one/uno
	assert_output "$(cat "$to_add")"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-2"
		echo "- task"
		echo "---"
	} >>"$to_add"

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

	run zit cat -type etikett
	assert_output "$(cat "$expected")"
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
		echo "---"
	} >"$one"

	run zit new -edit=false -predictable-hinweisen "$one"

	two="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "---"
	} >"$two"

	run zit new -edit=false -predictable-hinweisen "$two"

	three="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "---"
	} >"$three"

	run zit new -edit=false -predictable-hinweisen "$three"

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
	} >>"$expected_organize"

	run zit organize -predictable-hinweisen -mode commit-directly -group-by priority,w task <"$expected_organize"
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

	run zit new -edit=false -predictable-hinweisen "$one"

	two="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "---"
	} >"$two"

	run zit new -edit=false -predictable-hinweisen "$two"

	three="$(mktemp)"
	{
		echo "---"
		echo "# 3"
		echo "- priority-1"
		echo "- task"
		echo "- w-2022-07-06"
		echo "---"
	} >"$three"

	run zit new -edit=false -predictable-hinweisen "$three"

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
	} >>"$expected_organize"

	run zit organize -predictable-hinweisen -verbose -mode commit-directly -group-by priority,w task <"$expected_organize"

	one="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- priority-2"
		echo "- task"
		echo "- w-2022-07-08"
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
		echo "---"
	} >"$three"

	run zit show two/uno
	assert_output "$(cat "$three")"
}

function zettels_in_correct_places { # @test
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

	run zit new -edit=false -predictable-hinweisen "$one"

	expected_organize="$(mktemp)"
	{
		echo
		echo "# inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2"
		echo
		echo "- [one/uno] jabra coral usb_a-to-usb_c cable"
	} >"$expected_organize"

	run zit organize -predictable-hinweisen -mode output-only -group-by inventory \
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

	run zit organize -predictable-hinweisen -mode commit-directly <"$first_organize"

	expected_etiketten="$(mktemp)"
	{
    echo test1-wow
	} >"$expected_etiketten"

  run zit cat -type etikett
  assert_output "$(cat "$expected_etiketten")"

  mkdir -p one
	{
		echo "---"
		echo "- test4"
		echo "---"
	} >"one/uno.md"

	run zit checkin one/uno.md
  assert_output --partial "[one/uno "
  assert_output --partial "(updated)"

	expected_etiketten="$(mktemp)"
	{
    echo test4
	} >"$expected_etiketten"

  run zit cat -type etikett
  assert_output "$(cat "$expected_etiketten")"

  mkdir -p one
	{
		echo "---"
		echo "- test4"
    echo "- test1-ok"
		echo "---"
	} >"one/uno.md"

	run zit checkin one/uno.md
  assert_output --partial "[one/uno "
  assert_output --partial "(updated)"

	expected_etiketten="$(mktemp)"
	{
    echo test1-ok
    echo test4
	} >"$expected_etiketten"

  run zit cat -type etikett
  assert_output "$(cat "$expected_etiketten")"
}
