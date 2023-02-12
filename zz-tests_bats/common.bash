#! /bin/bash -e

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
	-abbreviate-hinweisen=false
	-abbreviate-shas=false
	-predictable-hinweisen
	-print-typen=false
	-print-time=false
)

function run_zit {
	cmd="$1"
	shift
	#shellcheck disable=SC2068
	run zit "$cmd" ${cmd_zit_def[@]} "$@"
}

function run_zit_init {
	run_zit init -yin <(cat_yin) -yang <(cat_yang)
	assert_success
}

function run_zit_init_disable_age {
	run_zit init -yin <(cat_yin) -yang <(cat_yang) -disable-age
	assert_success
}
