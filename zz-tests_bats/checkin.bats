#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	run_zit checkout :z,t,e
	assert_success

	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		- etikett-one
		! md
		---

		newest body
	EOM

	cat >one/dos.zettel <<-EOM
		---
		# dos wildly different
		- etikett-two
		! md
		---

		dos newest body
	EOM

	cat >md.typ <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM

	cat >zz-archive.etikett <<-EOM
		hide = true
	EOM
}

teardown() {
	rm_from_version "$version"
}

function dirty_one_virtual() {
	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		- etikett-one
		- %virtual
		! md
		---

		newest body
	EOM
}

function checkin_simple_one_zettel { # @test
	run_zit checkin one/uno.zettel
	assert_success
	assert_output - <<-EOM
		[etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}

function checkin_simple_one_zettel_virtual_etikett { # @test
	dirty_one_virtual
	run_zit checkin one/uno.zettel
	assert_success
	assert_output - <<-EOM
		[etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" %virtual etikett-one]
	EOM

	run_zit show one/uno
	assert_success
	assert_output - <<-EOM
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}

function checkin_complex_zettel_etikett_negation { # @test
	run_zit checkin ^etikett-two.z
	assert_success
	assert_output - <<-EOM
		[etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}

function checkin_simple_all { # @test
	# TODO: modify this to support "." for everything
	run_zit checkin .z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
		[etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett-two@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz-archive@b8cd0eaa1891284eafdf99d3acc2007a3d4396e8a7282335f707d99825388a93]
		[one/dos@b5c4fbaac3b71657edee74de4b947f13dfa104715feb8bab7cfa4dd47cafa3db !md "dos wildly different" etikett-two]
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM

	run_zit show -format log :?z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
		[etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett-two@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz-archive@b8cd0eaa1891284eafdf99d3acc2007a3d4396e8a7282335f707d99825388a93]
		[one/dos@b5c4fbaac3b71657edee74de4b947f13dfa104715feb8bab7cfa4dd47cafa3db !md "dos wildly different" etikett-two]
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}

function checkin_simple_all_dry_run { # @test
	# TODO fix race condition allowing for one double entry
	# TODO: modify this to support "." for everything
	run_zit checkin -dry-run .z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
		[etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett-two@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[zz-archive@b8cd0eaa1891284eafdf99d3acc2007a3d4396e8a7282335f707d99825388a93]
		[one/dos@b5c4fbaac3b71657edee74de4b947f13dfa104715feb8bab7cfa4dd47cafa3db !md "dos wildly different" etikett-two]
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM

	run_zit show -format log :z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function checkin_simple_typ { # @test
	run_zit checkin .t
	assert_success
	assert_output - <<-EOM
		[!md@220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
	EOM

	run_zit last -format bestandsaufnahme-sans-tai
	assert_success
	assert_output - <<-EOM
		---
		Akte 220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217
		Gattung Typ
		Kennung md
		---
	EOM

	run_zit show !md:t
	assert_success
	assert_output - <<-EOM
		[!md@220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
	EOM

	run_zit show -format typ.vim-syntax-type !md:typ
	assert_success
	assert_output 'toml'

	run_zit show -format typ.vim-syntax-type one/uno
	assert_success
	assert_output 'test'
}

function checkin_simple_etikett { # @test
	run_zit checkin zz-archive.e
	assert_success
	assert_output - <<-EOM
		[zz-archive@b8cd0eaa1891284eafdf99d3acc2007a3d4396e8a7282335f707d99825388a93]
	EOM

	run_zit last -format bestandsaufnahme-sans-tai
	assert_success
	assert_output - <<-EOM
		---
		Akte b8cd0eaa1891284eafdf99d3acc2007a3d4396e8a7282335f707d99825388a93
		Gattung Etikett
		Kennung zz-archive
		---
	EOM

	run_zit show -format akte zz-archive?e
	assert_success
	assert_output - <<-EOM
		hide = true
	EOM
}

function checkin_zettel_typ_has_commit_hook { # @test
	cat >typ_with_hook.typ <<-EOM
		hooks = """
		return {
		  on_new = function (kinder)
		    kinder["Etiketten"]["on_new"] = true
		    return nil
		  end,
		  on_pre_commit = function (kinder, mutter)
		    kinder["Etiketten"]["on_pre_commit"] = true
		    return nil
		  end,
		}
		"""
	EOM

	run_zit checkin -delete typ_with_hook.typ
	assert_success
	assert_output - <<-EOM
		[!typ_with_hook@1f6b9061059a83822901612bc050dd7d966bb5a2ceb917549ca3881728854477]
		          deleted [typ_with_hook.typ]
	EOM

	run_zit new -edit=false - <<-EOM
		---
		# test lua
		! typ_with_hook
		---

		should add new etikett
	EOM
	assert_success
	assert_output - <<-EOM
		[two/uno@edf7b6df934442ad0d6ac9fe4132c5e588391eb307fbbdc3ab6de780e17245a5 !typ_with_hook "test lua" on_new on_pre_commit]
	EOM
}

function checkin_zettel_with_komment { # @test
	run_zit checkin -print-bestandsaufnahme=true -komment "message" one/uno.zettel
	assert_success
	assert_output --regexp - <<-'EOM'
		\[etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\]
		\[etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\]
		\[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one\]
		\[[0-9]+\.[0-9]+@[0-9a-f]{64} "message"\]
	EOM
}
