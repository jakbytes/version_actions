#!/bin/bash

# Record the current stash list
BEFORE_STASH=$(git rev-parse --verify --quiet refs/stash)

# Stash unstaged changes - keep the staged ones
git stash push --keep-index --quiet

# Record the new stash list
AFTER_STASH=$(git rev-parse --verify --quiet refs/stash)

# Run tests
go test ./...

# Capture the go test exit code
TEST_EXIT_CODE=$?

# Pop the stashed changes back, only if a new stash was created
if [ "$BEFORE_STASH" != "$AFTER_STASH" ]; then
    git stash pop --quiet
fi

# Exit with the captured exit code
# (non-zero if go test failed)
exit $TEST_EXIT_CODE
