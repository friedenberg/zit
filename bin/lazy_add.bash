#! /bin/bash -e

rm -rf build_options
gmake
go clean
git add .
git commit -m update || true
git push
