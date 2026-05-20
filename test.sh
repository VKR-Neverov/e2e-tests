#!/bin/bash

# Define the number of times to run the tests
num_runs=10

for ((i=1; i<=$num_runs; i++)); do
    echo "Run $i:"
    # Run your test command here
    go clean -testcache
    go test -v ./internal/tests/...
done
