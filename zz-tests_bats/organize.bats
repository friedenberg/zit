#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

cmd_def_organize=(
	"${cmd_zit_def[@]}"
	-prefix-joints=false
	-refine=true
)

cmd_def_organize_prefix_joints=(
	"${cmd_zit_def[@]}"
	-prefix-joints=true
	-refine=true
)

function organize_empty { # @test
	run_zit organize "${cmd_def_organize[@]}" -mode output-only
	assert_success
	assert_output_unsorted - <<-EOM
	EOM
}

function organize_empty_commit { # @test
	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly <<-EOM
		- test
	EOM

	assert_success
	assert_output - <<-EOM
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "test"]
	EOM
}

function organize_simple { # @test
	actual="$(mktemp)"
	run_zit organize "${cmd_def_organize[@]}" -mode output-only :z,e,t >"$actual"
	assert_success
	assert_output_unsorted - <<-EOM

		- [!md]
		- [one/dos !md tag-3 tag-4] wow ok again
		- [one/uno !md tag-3 tag-4] wow the first
		- [tag-1]
		- [tag-2]
		- [tag-3]
		- [tag-4]
		- [tag]
	EOM
}

function organize_simple_commit { # @test
	run_zit checkout one/uno
	assert_success

	run_zit organize -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all, %virtual_etikett
		- [   !md   ]
		- [   tag  ]
		- [   tag-1]
		- [   tag-2]
		- [   tag-3]
		- [   tag-4]
		- [one/dos   !md tag-3 tag-4] wow ok again
		- [one/uno   !md tag-3 tag-4] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 %virtual_etikett new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett-for @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" %virtual_etikett new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" %virtual_etikett new-etikett-for-all tag-3 tag-4]
		[tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 %virtual_etikett new-etikett-for-all]
		[tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 %virtual_etikett new-etikett-for-all]
		[tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 %virtual_etikett new-etikett-for-all]
		[tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 %virtual_etikett new-etikett-for-all]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 %virtual_etikett new-etikett-for-all]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" new-etikett-for-all tag-3 tag-4]
		[tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
	EOM
}

function organize_simple_checkedout_matchesmutter { # @test
	run_zit checkout one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	run_zit organize -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos   !md tag-3 tag-4] wow ok again
		- [one/uno   !md tag-3 tag-4] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 new-etikett-for-all]
		[-tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett-for @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 new-etikett-for-all]
		[-tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit status one/dos.zettel
	assert_success
	assert_output - <<-EOM
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" new-etikett-for-all tag-3 tag-4]
	EOM
}

function organize_simple_checkedout_merge_no_conflict { # @test
	run_zit checkout one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	cat - >one/dos.zettel <<-EOM
		---
		# wow ok again
		- get_this_shit_merged
		- tag-3
		- tag-4
		! md
		---

		not another one, now with a different body
	EOM

	run_zit organize -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos   !md tag-3 tag-4] wow ok again
		- [one/uno   !md tag-3 tag-4] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 new-etikett-for-all]
		[-tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett-for @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 new-etikett-for-all]
		[-tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit status one/dos.zettel
	assert_success
	assert_output - <<-EOM
		          changed [one/dos.zettel @7ac3bdeb0ac8fd96cd7f8700a4bbc7a5d777fe26c50b52c20ecd726b255ec3d0 !md "wow ok again" get_this_shit_merged new-etikett-for-all tag-3 tag-4]
	EOM
}

function organize_simple_checkedout_merge_conflict { # @test
	#TODO-project-2022-zit-collapse_skus
	cat - >txt.typ <<-EOM
		inline-akte = true
	EOM

	cat - >txt2.typ <<-EOM
		inline-akte = true
	EOM

	run_zit checkin -delete .t
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [txt.typ]
		          deleted [txt2.typ]
		[!txt2 @c9627d6d0f0a88e6cbc93a5ccb4657a7b274655c1b89c53cbff92ecae5f6c583]
		[!txt @c9627d6d0f0a88e6cbc93a5ccb4657a7b274655c1b89c53cbff92ecae5f6c583]
	EOM

	run_zit checkout one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	cat - >one/dos.zettel <<-EOM
		---
		# wow ok again modified
		- get_this_shit_merged
		- tag-3
		- tag-4
		! txt
		---

		not another one, conflict time
	EOM

	run_zit organize -mode commit-directly :z,e,t <<-EOM
		---
		! txt2
		---

		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos   !md tag-3 tag-4] wow ok again different
		- [one/uno   !md tag-3 tag-4] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 !txt2 new-etikett-for-all]
		[-tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett-for @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new-etikett @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[new @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again different" new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 !txt2 new-etikett-for-all]
		[-tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !txt2 new-etikett-for-all]
		[new-etikett-for-all @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again different" new-etikett-for-all tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit status one/dos.zettel
	assert_success
	assert_output - <<-EOM
		       conflicted [one/dos.zettel]
	EOM
}

function organize_hides_hidden_etiketten_from_organize { # @test
	run_zit schlummernd-add zz-archive
	assert_success
	assert_output ''

	to_add="$(mktemp)"
	{
		echo ---
		echo "# split hinweis for usability"
		echo - project-2021-zit
		echo - zz-archive-task-done
		echo ! md
		echo ---
	} >"$to_add"

	run_zit new -edit=false "$to_add"
	assert_success
	assert_output - <<-EOM
		[project @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[project-2021 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[project-2021-zit @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz-archive @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz-archive-task @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz-archive-task-done @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "split hinweis for usability" project-2021-zit zz-archive-task-done]
	EOM

	run_zit show two/uno
	assert_success
	assert_output - <<-EOM
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "split hinweis for usability" project-2021-zit zz-archive-task-done]
	EOM

	run_zit organize -mode output-only project-2021-zit:z
	assert_success
	assert_output - <<-EOM
		---
		- project-2021-zit
		---
	EOM
}

function organize_dry_run { # @test
	expected_show="$(mktemp)"
	# shellcheck disable=SC2154
	zit show "${cmd_zit_def[@]}" -format log :z,e,t >"$expected_show"

	run_zit organize -dry-run -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos  ] wow ok again
		- [one/uno  ] wow the first
	EOM
	assert_success

	run_zit show -format log :z,e,t
	assert_success
	assert_output_unsorted "$(cat "$expected_show")"
}

function organize_with_typ_output { # @test
	run_zit organize "${cmd_def_organize[@]}" -mode output-only !md:z
	assert_success
	assert_output - <<-EOM
		---
		! md
		---

		- [one/dos tag-3 tag-4] wow ok again
		- [one/uno tag-3 tag-4] wow the first
	EOM
}

function organize_with_typ_commit { # @test
	run_zit organize -mode commit-directly !md:z <<-EOM
		---
		! txt
		---

		- [one/dos tag-3 tag-4] wow ok again
		- [one/uno tag-3 tag-4] wow the first
	EOM

	assert_success
	assert_output_unsorted - <<-EOM
		[!txt @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !txt "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !txt "wow the first" tag-3 tag-4]
	EOM
}

function modify_bezeichnung { # @test
	run_zit organize -mode commit-directly :z,e,t <<-EOM

		- [   !md   ]
		- [   tag  ]
		- [   tag-1]
		- [   tag-2]
		- [   tag-3]
		- [   tag-4]
		- [one/dos   !md tag-3 tag-4] wow ok again was modified
		- [one/uno   !md tag-3 tag-4] wow the first was modified too
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again was modified" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first was modified too" tag-3 tag-4]
	EOM
}

function add_named { # @test
	# TODO modify organize to not require query group or else accidentally output unchanged objektes
	run_zit organize -mode commit-directly :e <<-EOM
		# with-tag
		- [added_tag]
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[added_tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 with-tag]
		[with-tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[with @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM
}

function organize_v5_outputs_organize_one_etikett { # @test
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
		[ok @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow" ok]
	EOM

	run_zit expand-zettel-id o/u
	assert_success
	assert_output 'one/uno'

	run_zit organize "${cmd_def_organize[@]}" -mode output-only ok
	assert_success
	assert_output - <<-EOM
		---
		- ok
		---

		- [two/uno !md] wow
	EOM
}

function organize_v5_outputs_organize_two_etiketten { # @test
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
		[brown @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[ok @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow" brown ok]
	EOM

	run_zit organize "${cmd_def_organize[@]}" -mode output-only ok brown
	assert_success
	assert_output - <<-EOM
		---
		- brown
		- ok
		---

		- [two/uno !md] wow
	EOM

	run_zit organize "${cmd_def_organize[@]}" \
		-mode commit-directly \
		ok brown <<-EOM
			      # ok

			- [two/uno !md] wow
		EOM

	assert_success
	assert_output - <<-EOM
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow" ok]
	EOM

	run_zit show -format text two/uno
	assert_success
	assert_output - <<-EOM
		---
		# wow
		- ok
		! md
		---
	EOM
}

function organize_v5_outputs_organize_one_etiketten_group_by_one { # @test
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
		[priority @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[priority-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[priority-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[task @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "wow" priority-1 priority-2 task]
	EOM

	run_zit organize "${cmd_def_organize[@]}" \
		-mode output-only \
		-group-by priority task
	assert_success
	assert_output - <<-EOM
		---
		- task
		---

		          # priority-1

		- [two/uno !md priority-2] wow

		          # priority-2

		- [two/uno !md priority-1] wow
	EOM

	return

	# shellcheck disable=2317
	run_zit organize "${cmd_def_organize_prefix_joints[@]}" \
		-mode output-only \
		-group-by priority task

	# shellcheck disable=2317
	assert_success
	# shellcheck disable=2317
	assert_output - <<-EOM
		---
		- task
		---

		          # priority

		         ##         -1

		- [two/uno  !md] wow

		         ##         -2

		- [two/uno  !md] wow
	EOM
}

function organize_v5_outputs_organize_two_zettels_one_etiketten_group_by_one { # @test
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
		[priority-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[priority @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[task @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "one/uno" priority-1 task]
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
		[priority-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/tres @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "two/dos" priority-2 task]
	EOM

	# add prefix joints
	run_zit organize "${cmd_def_organize[@]}" -mode output-only -group-by priority task
	assert_success
	assert_output - <<-EOM
		---
		- task
		---

		           # priority-1

		- [two/uno !md] one/uno

		           # priority-2

		- [one/tres !md] two/dos
	EOM
}

function organize_v5_commits_organize_one_etiketten_group_by_two { # @test
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
		echo "- [one/dos !md] two/dos"
		echo
		echo "## priority-2"
		echo
		echo "### w-2022-07-07"
		echo
		echo "- [one/uno !md] one/uno"
		echo
		echo "###"
		echo
		echo "- [two/uno !md] 3"
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

	run_zit show -format text one/uno
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

	run_zit show -format text two/uno
	assert_success
	assert_output "$(cat "$to_add")"
}

function organize_v5_commits_organize_one_etiketten_group_by_two_new_zettels { # @test
	skip
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
		echo "- new zettel one"
		echo "## priority-1"
		echo "- new zettel two"
		echo "### w-2022-07-06"
		echo "- [one/dos !md] two/dos"
		echo "## priority-2"
		echo "### w-2022-07-07"
		echo "- [one/uno !md] one/uno"
		echo "###"
		echo "- new zettel three"
		echo "- [two/uno !md] 3"
	} >"$expected_organize"

	run_zit organize \
		"${cmd_def_organize[@]}" \
		-mode commit-directly \
		-group-by priority,w \
		task <"$expected_organize"
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

	run_zit show -format text one/uno
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

	run_zit show -format text two/uno
	assert_success
	assert_output "$(cat "$to_add")"

	run_zit show -format text one/tres
	assert_success

	run_zit show -format text two/dos
	assert_success

	run_zit show -format text three/uno
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

function organize_v5_commits_no_changes { # @test
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
		[priority-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[priority @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[task @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "one/uno" priority-1 task w-2022-07-07]
		[w-2022-07-07 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[w-2022-07 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[w-2022 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[w @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
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
		[w-2022-07-06 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/tres @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "two/dos" priority-1 task w-2022-07-06]
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
		[two/dos @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "3" priority-1 task w-2022-07-06]
	EOM

	# TODO add prefix joints
	run_zit organize "${cmd_def_organize[@]}" \
		-mode output-only \
		-group-by priority,w task
	assert_success
	assert_output - <<-EOM
		---
		- task
		---

		           # priority-1

		          ## w-2022-07-06

		- [one/tres !md] two/dos
		- [two/dos !md] 3

		          ## w-2022-07-07

		- [two/uno !md] one/uno
	EOM

	run_zit organize "${cmd_def_organize[@]}" \
		-mode commit-directly \
		-group-by priority,w task \
		<<-EOM
			---
			- task
			---

			           # priority

			          ##         -1

			         ### w

			        ####  -2022-07

			       #####          -06

			- [two/uno   !md] one/uno

			       #####          -07

			- [one/tres  !md] two/dos
			- [two/dos   !md] 3

		EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[one/tres @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "two/dos" priority-1 task w-2022-07-07]
		[two/dos @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "3" priority-1 task w-2022-07-07]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "one/uno" priority-1 task w-2022-07-06]
	EOM
}

function organize_v5_commits_dependent_leaf { # @test
	skip
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
		echo "- [one/dos !md] two/dos"
		echo "- [two/uno !md] 3"
		echo "#### -08"
		echo "- [one/uno !md] one/uno"
		echo "###"
	} >"$expected_organize"

	run_zit organize "${cmd_def_organize[@]}" -verbose -mode commit-directly -group-by priority,w task <"$expected_organize"
	assert_success
}

function organize_v5_zettels_in_correct_places { # @test
	one="$(mktemp)"
	{
		echo "---"
		echo "# jabra coral usb_a-to-usb_c cable"
		echo "- inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2"
		echo "---"
	} >"$one"

	run_zit new -edit=false "$one"

	run_zit organize "${cmd_def_organize[@]}" \
		-mode output-only -group-by inventory \
		inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2
	assert_success

	# TODO add prefix joints
	assert_output - <<-EOM
		---
		- inventory-pipe_shelves-atheist_shoes_box-jabra_yellow_box_2
		---

		- [two/uno !md] jabra coral usb_a-to-usb_c cable
	EOM
}

function organize_v5_etiketten_correct { # @test

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly <<-EOM
		# test1
		## -wow

		- zettel bez
	EOM
	assert_success

	assert_output - <<-EOM
		[test1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[test1-wow @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "zettel bez" test1-wow]
	EOM

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
		[test4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md test4]
	EOM

	# TODO-P2 fix issue with kennung schwanzen
	# run_zit cat-etiketten-schwanzen
	# assert_output - <<-EOM
	# EOM

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
		[test1-ok @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md test1-ok test4]
	EOM
}

function organize_remove_anchored_metadata { # @test
	run_zit show tag-3:z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly tag-3 <<-EOM
		---
		- tag-3
		---
	EOM

	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-4]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-4]
	EOM

	run_zit show tag-3:z
	assert_success
	assert_output_unsorted - <<-EOM
	EOM
}

function organize_update_checkout { # @test
	run_zit checkout one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly :z <<-EOM
		---
		- test
		---

		- [one/dos  !md tag-3 tag-4] wow ok again
		- [one/uno  !md tag-3 tag-4] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4 test]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4 test]
		[test @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM

	run_zit status
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4 test]
	EOM
}

function organize_update_checkout_remove_etiketten { # @test
	run_zit checkout one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly :z <<-EOM
		- [one/dos  !md] wow ok again
		- [one/uno  !md] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
	EOM

	run_zit status
	assert_success
	assert_output_unsorted - <<-EOM
		             same [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
	EOM
}

function create_structured_zettels { # @test
	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly <<-EOM
		---
		- test
		---

		- [/] first
		- [/  !task tag-3] second
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!task @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/tres @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !task "second" tag-3 test]
		[test @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "first" test]
	EOM
}

function description_with_literal_characters { # @test
	run_zit organize "${cmd_def_organize[@]}" -mode commit-directly <<-EOM
		- [terb/ala !md payee] thoughts on quincey's contract / scope of work
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[payee @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[terb/ala @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md "thoughts on quincey's contract / scope of work" payee]
	EOM
}
