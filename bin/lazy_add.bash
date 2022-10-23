#! /bin/bash -e

rm -rf build_options
_MR_BUILD_AND_WATCH_ONCE=1 fish -c mr-build-and-watch
git add .
git commit -m update || true
git push
