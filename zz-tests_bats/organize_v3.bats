#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output
}

teardown() {
	rm_from_version
}

cmd_def_organize_v3=(
	-prefix-joints=true
	-metadatei-header=false
	-refine=true
)

function organize_v3_outputs_organize_one_etikett { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success
	assert_output - <<-EOM
		[-ok@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow"]
	EOM

	run zit expand-hinweis o/u
	assert_success
	assert_output 'one/uno'

	expected_organize="$(mktemp)"
	{
		echo
		echo "          # ok"
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize_v3[@]}" -mode output-only ok
	assert_success
	assert_output "$(cat "$expected_organize")"
}

function organize_v3_outputs_organize_two_etiketten { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "- brown"
		echo "! md"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success
	assert_output_unsorted - <<-EOM
		[-brown@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-ok@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow"]
	EOM

	expected_organize="$(mktemp)"
	{
		echo
		echo "          # brown, ok"
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize_v3[@]}" -mode output-only ok brown
	assert_success
	assert_output "$(cat "$expected_organize")"

	{
		echo
		echo "      # ok"
		echo
		echo "- [o/u] wow"
		echo
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize_v3[@]}" -mode commit-directly ok brown <"$expected_organize"
	assert_success
	assert_output - <<-EOM
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow"]
	EOM

	expected_zettel="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "! md"
		echo "---"
	} >"$expected_zettel"

	run_zit show one/uno
	assert_success
	assert_output "$(cat "$expected_zettel")"
}

function organize_v3_outputs_organize_one_etiketten_group_by_one { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

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

	run_zit new -edit=false "$to_add"
	assert_success
	assert_output - <<-EOM
		[-priority@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-priority-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-priority-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-task@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow"]
	EOM

	expected_organize="$(mktemp)"
	{
		echo
		echo "          # task"
		echo
		echo "         ## priority"
		echo
		echo "        ###         -1"
		echo
		echo "- [one/uno] wow"
		echo
		echo "        ###         -2"
		echo
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize_v3[@]}" -mode output-only -group-by priority task
	assert_success
	assert_output "$(cat "$expected_organize")"
}

function organize_v3_outputs_organize_two_zettels_one_etiketten_group_by_one { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# one/uno"
		echo "- task"
		echo "- priority-1"
		echo "! md"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success
	assert_output_unsorted - <<-EOM
		[-priority-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-priority@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-task@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "one/uno"]
	EOM

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# two/dos"
		echo "- task"
		echo "- priority-2"
		echo "! md"
		echo "---"
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success
	assert_output - <<-EOM
		[-priority-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "two/dos"]
	EOM

	expected_organize="$(mktemp)"
	{
		echo
		echo "          # task"
		echo
		echo "         ## priority"
		echo
		echo "        ###         -1"
		echo
		echo "- [one/uno] one/uno"
		echo
		echo "        ###         -2"
		echo
		echo "- [one/dos] two/dos"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize_v3[@]}" -mode output-only -group-by priority task
	assert_success
	assert_output "$(cat "$expected_organize")"
}

function organize_v3_commits_organize_one_etiketten_group_by_two { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

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

	run_zit new -edit=false "$to_add"
	assert_success

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

	run_zit new -edit=false "$to_add"
	assert_success

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
	assert_success

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

	run_zit organize "${cmd_def_organize_v3[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
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

	run_zit show one/uno
	assert_success
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

	run_zit show two/uno
	assert_success
	assert_output "$(cat "$to_add")"
}

function organize_v3_commits_organize_one_etiketten_group_by_two_new_zettels { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

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

	run_zit new -edit=false "$to_add"
	assert_success

	expected="$(mktemp)"
	{
		echo priority-1
		echo task
		echo w-2022-07-07
	} >"$expected"

	run_zit cat-etiketten-schwanzen
	assert_success
	assert_output_unsorted "$(cat "$expected")"

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

	run_zit new -edit=false "$to_add"
	assert_success

	{
		echo priority-1
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	run_zit cat-etiketten-schwanzen
	assert_success
	assert_output_unsorted "$(cat "$expected")"

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
	assert_success

	{
		echo priority-1
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	run_zit cat-etiketten-schwanzen
	assert_success
	assert_output_unsorted "$(cat "$expected")"

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

	run_zit organize "${cmd_def_organize_v3[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
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

	run_zit show one/uno
	assert_success
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

	run_zit show two/uno
	assert_success
	assert_output "$(cat "$to_add")"

	run_zit show one/tres
	assert_success

	run_zit show two/dos
	assert_success

	run_zit show three/uno
	assert_success

	{
		echo priority-1
		echo priority-2
		echo task
		echo w-2022-07-06
		echo w-2022-07-07
	} >"$expected"

	# TODO
	# run zit cat-etiketten-schwanzen
	# assert_output "$(cat "$expected")"
}

function organize_v3_commits_no_changes { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
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

	run_zit new -edit=false "$one"
	assert_success
	assert_output_unsorted - <<-EOM
		[-priority-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-priority@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-task@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-w-2022-07-07@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-w-2022-07@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-w-2022@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-w@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "one/uno"]
	EOM

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
	assert_success
	assert_output_unsorted - <<-EOM
		[-w-2022-07-06@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "two/dos"]
	EOM

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
	assert_success
	assert_output_unsorted - <<-EOM
		[two/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "3"]
	EOM

	expected_organize="$(mktemp)"
	{
		echo
		echo "# task"
		echo
		echo " ## priority-1"
		echo
		echo "  ### w-2022-07-06"
		echo
		echo "  - [two/uno] 3"
		echo "  - [one/dos] two/dos"
		echo
		echo "  ### w-2022-07-07"
		echo
		echo "  - [one/uno] one/uno"
		echo
	} >"$expected_organize"

	# run_zit organize "${cmd_def_organize_v3[@]}" -prefix-joints=false -mode output-only -group-by priority,w task
	run_zit organize "${cmd_def_organize_v3[@]}" -mode commit-directly -group-by priority,w task <"$expected_organize"
	assert_success
	# assert_output "$(cat "$expected_organize")"
	assert_output "no changes"

	run_zit show one/uno
	assert_success
	assert_output "$(cat "$one")"

	run_zit show one/dos
	assert_success
	assert_output "$(cat "$two")"

	run_zit show two/uno
	assert_success
	assert_output "$(cat "$three")"
}

function organize_v3_commits_dependent_leaf { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
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

	run_zit new -edit=false "$one"
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

	run_zit new -edit=false "$two"
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

	run_zit new -edit=false "$three"
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

	run_zit organize "${cmd_def_organize_v3[@]}" -verbose -mode commit-directly -group-by priority,w task <"$expected_organize"
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

	run_zit show one/uno
	assert_success
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

	run_zit show one/dos
	assert_success
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

	run_zit show two/uno
	assert_success
	assert_output "$(cat "$three")"
}

function organize_v3_zettels_in_correct_places { # @test
	cd "$BATS_TEST_TMPDIR" || exit 1
	run_zit_init_disable_age

	one="$(mktemp)"
	{
		echo "---"
		echo "# jabra coral usb_a-to-usb_c cable"
		echo "- inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2"
		echo "---"
	} >"$one"

	run_zit new -edit=false "$one"

	run_zit organize "${cmd_def_organize_v3[@]}" \
		-mode output-only -group-by inventory \
		inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2
	assert_success

	assert_output - <<-EOM

		          # inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2

		         ## inventory

		        ###          -pipe_shelves-atheist_shoes_box-jabra_yellow_box_2

		- [one/uno] jabra coral usb_a-to-usb_c cable
	EOM
}

function organize_v3_etiketten_correct { # @test
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

	run_zit organize "${cmd_def_organize_v3[@]}" -mode commit-directly <"$first_organize"
	assert_success

	expected_etiketten="$(mktemp)"
	{
		echo test1-wow
	} >"$expected_etiketten"

	run_zit cat-etiketten-schwanzen
	assert_success
	assert_output "$(cat "$expected_etiketten")"

	mkdir -p one
	{
		echo "---"
		echo "- test4"
		echo "! md"
		echo "---"
	} >"one/uno.zettel"

	run_zit checkin one/uno.zettel
	assert_success
	assert_output - <<-EOM
		[-test4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md test4]
	EOM

	expected_etiketten="$(mktemp)"
	{
		echo test4
	} >"$expected_etiketten"

	run_zit cat-etiketten-schwanzen
	assert_output "$(cat "$expected_etiketten")"

	mkdir -p one
	{
		echo "---"
		echo "- test4"
		echo "- test1-ok"
		echo "! md"
		echo "---"
	} >"one/uno.zettel"

	run_zit checkin one/uno.zettel
	assert_output - <<-EOM
		[-test1-ok@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md test1-ok test4]
	EOM

	expected_etiketten="$(mktemp)"
	{
		echo test1-ok
		echo test4
	} >"$expected_etiketten"

	run zit cat-etiketten-schwanzen
	assert_output_unsorted "$(cat "$expected_etiketten")"
}
