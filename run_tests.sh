#!/bin/bash

JSONNET_CMD=jsonnet
JSONNET_ARGS="--ext-str GITLAB_GKE_DOMAIN=\"gitlab.io\" --ext-str release.name=3 --ext-str release.service=tiller --ext-str serviceName=cow"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
for JSONNET_FILE in $(find $DIR/examples/ -name "*.jsonnet") ; do
  TEST_OUTPUT="$($JSONNET_CMD $JSONNET_ARGS "$JSONNET_FILE" 2>&1)"
  TEST_EXIT_CODE="$?"

  if [ "$TEST_EXIT_CODE" -gt 0 ] ; then
    echo -e "FAIL: $JSONNET_FILE"
    echo "This run's output:"
    echo "$TEST_OUTPUT"
    exit 1
  fi

  GOLDEN_OUTPUT=$(cat "$JSONNET_FILE.golden")
  if [ "$TEST_OUTPUT" != "$GOLDEN_OUTPUT" ] ; then
    echo -e "FAIL: $JSONNET_FILE"
    echo "This run's output:"
    echo "$TEST_OUTPUT"
    echo "Expected:"
    echo "$GOLDEN_OUTPUT"
    exit 1
  fi
done
