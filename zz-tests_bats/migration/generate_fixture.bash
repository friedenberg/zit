#! /bin/bash -xe

dir_base="$(realpath "$(dirname "$0")")"
make build/zit
zit="$(realpath build/zit)"
v="$("$zit" store-version)"
d="$dir_base/v$v"

if [[ -d "$d" ]]; then
  chflags -R nouchg "$d"
  rm -rf "$d"
fi

mkdir -p "$d"

pushd "$d"
"$zit" init -yin "$dir_base/yin" -yang "$dir_base/yang" -disable-age

"$zit" new -predictable-hinweisen -edit=false - <<EOM
---
# wow ok
- tag-1
- tag-2
! md
---

this is the body aiiiiight
EOM

"$zit" new -predictable-hinweisen -edit=false - <<EOM
---
# wow ok again
- tag-3
- tag-4
! md
---

not another one
EOM

"$zit" checkout o/u
cat > one/uno.zettel <<EOM
---
# wow the first
- tag-3
- tag-4
! md
---

last time
EOM

"$zit" checkin -delete one/uno.zettel
