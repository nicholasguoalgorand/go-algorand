#!/usr/bin/env bash

set -e
set -x

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
OS=$("${SCRIPTPATH}/../ostype.sh")
ARCH=$("${SCRIPTPATH}/../archtype.sh")

curl -sL -o ~/gimme https://raw.githubusercontent.com/travis-ci/gimme/master/gimme
chmod +x ~/gimme
eval $(~/gimme $("${SCRIPTPATH}/../get_golang_version.sh"))

if [ "${OS}-${ARCH}" = "linux-arm" ]; then
    # for arm, no tests need to be invoked.
    exit 0
fi

GOPATHBIN=$(go env GOPATH)/bin
export PATH=$PATH:$GOPATHBIN

make fixcheck
scripts/travis/run_tests.sh;
scripts/travis/after_build.sh;
