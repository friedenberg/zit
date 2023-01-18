#! /bin/bash -e

ag '\b(\w+)\b "github.com/friedenberg/zit/src/\w+/\1"' -l0 |
	xargs -0 gsed -E -i 's#(\w+) ("github.com/friedenberg/zit/src/\w+/\1")#\2#g'

goimports -w ./
