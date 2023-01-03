#!/bin/sh
GO="${GO:-go}"
LDFLAGS="${LDFLAGS:--w -s}"

echo "Using:"
echo "GO='$GO'"
echo "LDFLAGS='$LDFLAGS'"
echo ""

mkdir -p dist/
for cpu in amd64 386 arm arm64; do
  for os in windows linux darwin; do
    DISTFILE="dist/majokko-${os}-${cpu}"
    CMD="GOOS='${os}' GOARCH='${cpu}' ${GO} build -o '$DISTFILE' -ldflags '$LDFLAGS' -trimpath ."
    echo "= $CMD"
    eval "$CMD" && echo "+ $(du -hs "$DISTFILE")" || echo "- Error: ${os} on ${cpu}"
  done
done
