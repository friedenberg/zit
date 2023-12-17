#! /bin/bash -ex

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
"$zit" init -yin "$dir_base/yin" -yang "$dir_base/yang" -disable-age -compression-type none

[ "$(zit show !md:t)" = "[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]" ]
[ "$(zit show :konfig)" = "[konfig@bb61ffad0cd4354654743ec604066a0a02db9ef188f695ce856acd284ee0612d]" ]

"$zit" new -predictable-hinweisen -edit=false - <<EOM
---
# wow ok
- tag-1
- tag-2
! md
---

this is the body aiiiiight
EOM

[ "$(zit show -format etiketten one/uno)" = "tag-1, tag-2" ]

"$zit" new -predictable-hinweisen -edit=false - <<EOM
---
# wow ok again
- tag-3
- tag-4
! md
---

not another one
EOM

[ "$(zit show -format etiketten one/dos)" = "tag-3, tag-4" ]

"$zit" checkout one/uno
cat >one/uno.zettel <<EOM
---
# wow the first
- tag-3
- tag-4
! md
---

last time
EOM

"$zit" checkin -delete one/uno.zettel

[ "$(zit show -format etiketten one/uno)" = "tag-3, tag-4" ]
