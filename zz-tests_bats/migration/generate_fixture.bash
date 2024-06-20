#! /bin/bash -e

dir_base="$(realpath "$(dirname "$0")")"
zit="$(realpath build/zit)"
v="$("$zit" store-version)"
d="${1:-$dir_base/v$v}"

if [[ -d $d ]]; then
  chflags -R nouchg "$d"
  rm -rf "$d"
fi

mkdir -p "$d"

pushd "$d"

function zit() (
  "$zit" "$@"
)

flags=(
  # -verbose
  -chrest-enabled=false
  -predictable-hinweisen 
)

zit init "${flags[@]}" -yin "$dir_base/yin" -yang "$dir_base/yang" -age none -compression-type none

[ "$(zit show "${flags[@]}" !md:t)" = "[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]" ]
[ "$(zit show "${flags[@]}" :konfig)" = "[konfig@e9412d561f5caaa9219ca6983ed842fefedf85c1aa10a98f271226070b9d1351]" ]

zit new "${flags[@]}" -edit=false - <<EOM
---
# wow ok
- tag-1
- tag-2
! md
---

this is the body aiiiiight
EOM

[ "$(zit show "${flags[@]}" -format etiketten one/uno)" = "tag-1, tag-2" ]

zit new "${flags[@]}" -edit=false - <<EOM
---
# wow ok again
- tag-3
- tag-4
! md
---

not another one
EOM

[ "$(zit show "${flags[@]}" -format etiketten one/dos)" = "tag-3, tag-4" ]

zit checkout "${flags[@]}" one/uno
cat >one/uno.zettel <<EOM
---
# wow the first
- tag-3
- tag-4
! md
---

last time
EOM

zit checkin "${flags[@]}" -delete one/uno.zettel

[ "$(zit show -format etiketten one/uno)" = "tag-3, tag-4" ]
