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
	echo "four"
	echo "five"
	echo "six"
)

cat_yang() (
	echo "uno"
	echo "dos"
	echo "tres"
	echo "quatro"
	echo "cinco"
	echo "seis"
)

cmd_zit_def=(
	# -abbreviate-hinweisen=false
	-predictable-hinweisen
	-print-typen=false
)

cmd_zit_organize=(
	zit
	organize
	-predictable-hinweisen
	-right-align=false
	-refine=true
	-metadatei-header=false
	"${cmd_zit_def[@]}"
)

function commits_no_changes { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	expected_organize="$(mktemp)"
	{
		echo "               # project-22q3-purchases"
		echo "              ## project"
		echo "             ###        -22q1-uws-140-mvp"
		echo "            ####                         -bathroom"
		echo "- buy new bath mats"
		echo "- buy new body towels"
		echo "            ####                         -main_room"
		echo "- buy new doormat"
		echo "           #####                                   -bed"
		echo "- buy bed"
		echo "             ###        -22q3-purchases"
		echo "- Buy white leather belt"
		echo "- buy bed"
		echo "- buy more gap underwear"
		echo "- buy new bath mats"
		echo "- buy new body towels"
		echo "- buy new doormat"
		echo "- buy philips hue bulbs"
		echo "- city water test"
		echo "- cord color knots"
		echo "- hat purchases"
		echo "- hats to buy"
		echo "- portable monitor"
		echo "- purchase from schoolhouse"
		echo "- purchase last napkin architecture art frame"
		echo "- resistance bands"
		echo "- return rei mat and buy new yoga mat"
		echo "- vertically-adjustable wall tv mount"
	} >"$expected_organize"

	run "${cmd_zit_organize[@]}" \
		-mode commit-directly \
		-group-by project \
		project-22q3-purchases \
		<"$expected_organize"

	expected="$(mktemp)"
	{
		echo '[o/u@f "Buy white leather belt"] (created)'
		echo '[o/d@1 "buy bed"] (created)'
		echo '[t/u@13 "buy more gap underwear"] (created)'
		echo '[o/t@d "buy new bath mats"] (created)'
		echo '[t/d@a "buy new body towels"] (created)'
		echo '[th/u@6 "buy new doormat"] (created)'
		echo '[o/q@b "buy philips hue bulbs"] (created)'
		echo '[tw/t@c "city water test"] (created)'
		echo '[th/d@16 "cord color knots"] (created)'
		echo '[f/u@4 "hat purchases"] (created)'
		echo '[o/c@d1 "hats to buy"] (created)'
		echo '[tw/q@8 "portable monitor"] (created)'
		echo '[th/t@de "purchase from schoolhouse"] (created)'
		echo '[f/d@a4 "purchase last napkin architecture art frame"] (created)'
		echo '[fi/u@3 "resistance bands"] (created)'
		echo '[o/s@2 "return rei mat and buy new yoga mat"] (created)'
		echo '[tw/c@ca "vertically-adjustable wall tv mount"] (created)'
	} >"$expected"
	assert_output "$(cat "$expected")"

	expected="$(mktemp)"
	{
		echo ''
		echo '   # project-22q3-purchases'
		echo ''
		echo '    ## project'
		echo ''
		echo '     ### -22q1-uws-140-mvp'
		echo ''
		echo '      #### -bathroom'
		echo ''
		echo '      - [o/t] buy new bath mats'
		echo '      - [tw/d] buy new body towels'
		echo ''
		echo '      #### -main_room'
		echo ''
		echo '       - [th/u] buy new doormat'
		echo ''
		echo '       ##### -bed'
		echo ''
		echo '       - [o/d] buy bed'
		echo ''
		echo '     ### -22q3-purchases'
		echo ''
		echo '     - [o/u] Buy white leather belt'
		echo '     - [o/d] buy bed'
		echo '     - [tw/u] buy more gap underwear'
		echo '     - [o/t] buy new bath mats'
		echo '     - [tw/d] buy new body towels'
		echo '     - [th/u] buy new doormat'
		echo '     - [o/q] buy philips hue bulbs'
		echo '     - [tw/t] city water test'
		echo '     - [th/d] cord color knots'
		echo '     - [fo/u] hat purchases'
		echo '     - [o/c] hats to buy'
		echo '     - [tw/q] portable monitor'
		echo '     - [th/t] purchase from schoolhouse'
		echo '     - [fo/d] purchase last napkin architecture art frame'
		echo '     - [fi/u] resistance bands'
		echo '     - [o/s] return rei mat and buy new yoga mat'
		echo '     - [tw/c] vertically-adjustable wall tv mount'
	} >"$expected"

	run "${cmd_zit_organize[@]}" -mode commit-directly \
		-group-by project project-22q3-purchases \
		<"$expected"
	assert_output "no changes"
}
