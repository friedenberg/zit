#! /bin/bash -e

rm -rf build_options
gmake
git add .
git commit -m update || true
git push
