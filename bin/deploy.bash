#! /bin/bash -e

#TODO don't reset options on lazy build?
rm -rf build_options
gmake
go clean -cache -fuzzcache
git add .
git commit -m update || true
git push
