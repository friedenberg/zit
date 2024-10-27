#! /bin/bash -e

while IFS= read -r -d '' file
do
  name="$(basename "$file")"
  dir="$(dirname "$file")"
  ln -fs "$(realpath "$file")" "$HOME/.local/share/pandoc/$dir/$name"
done <   <(find filters defaults -type f -print0)
