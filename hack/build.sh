#!/usr/bin/env bash
GIT_IMPORT=github.com/DanStough/epok/internal/buildinfo
GIT_COMMIT=$(git rev-parse --short HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
#DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

GOLDFLAGS="-X ${GIT_IMPORT}.gitCommit=${GIT_COMMIT}${GIT_DIRTY} -X ${GIT_IMPORT}.version=${VERSION}"

go run -ldflags "${GOLDFLAGS}" main.go "$@"