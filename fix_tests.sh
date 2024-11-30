#!/bin/bash

# Run `go test ./...`
go test ./...
exit_code=$?

# Check if the exit code is not 0
if [ $exit_code -ne 0 ]; then
  # Loop up to a maximum of 25 times
  for i in {1..50}; do
    echo "Attempt $i: Running aider to fix the error..."

    # Capture the output from `go test ./...`
    OUTPUT=$(go test ./... 2>&1)

    # Run the aider command with the captured output
    aider --sonnet --dark-mode --no-cache-prompts --no-auto-commits \
          --no-suggest-shell-commands --yes --yes-always \
          --message "Fix Error: $OUTPUT"

    # Re-run the tests after fixing
    go test ./...
    exit_code=$?

    # If the test passes, break out of the loop
    if [ $exit_code -eq 0 ]; then
      echo "Tests passed on attempt $i."
      break
    fi

    # Calculate the backoff time (exponential with a cap at 10 seconds)
    backoff_time=$((2 ** i))
    if [ $backoff_time -gt 10 ]; then
      backoff_time=10  # Cap at 10 seconds
    fi

    echo "Tests failed. Waiting for $backoff_time seconds before the next attempt..."
    sleep $backoff_time

    # If maximum attempts are reached, exit the loop
    if [ $i -eq 50 ]; then
      echo "Maximum attempts reached. Tests still failing."
      exit 1
    fi
  done
else
  echo "Tests passed successfully."
fi