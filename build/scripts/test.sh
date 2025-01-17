#!/bin/bash

[ -n "${DEBUG_SCRIPTS}" ] && set -x

set -euo pipefail

cd "${SRC_DIR}"
mkdir -p "${TEST_COVERAGE_DIR}"

LIST_OF_FILES=$(go list ./internal/... ./pkg/... ./test/... | grep -Ev "${GO_TEST_EXCLUDES}")
echo $LIST_OF_FILES

go version
# shellcheck disable=SC2086
CGO_ENABLED=1 go test \
	-gcflags=-l \
	-v \
	-count=1 \
	-race \
	-coverpkg ./... \
	-coverprofile="${TEST_COVERAGE_DIR}/coverage.out" \
	"${GO_TEST_FLAGS}" \
	${LIST_OF_FILES}
