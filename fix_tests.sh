#!/bin/bash

# Run `go test ./...`
go test ./...
exit_code=$?

# Check if the exit code is not 0
if [ $exit_code -ne 0 ]; then
  # Loop up to a maximum of 25 times
  for i in {1..25}; do
    echo "Repomix..."
    repomix

    echo "Attempt $i: Running aider to fix the error..."

    # Run the aider command with the specified options
    OUTPUT=$(go test ./...) # Capture the output from `go test ./...`
    aider --sonnet --dark-mode --no-cache-prompts --no-auto-commits \
          --no-suggest-shell-commands --yes --yes-always --read repomix-output.txt\
          --message "Fix Error: $OUTPUT"

    # Re-run the tests after fixing
    go test ./...
    exit_code=$?

    # If the test passes, break out of the loop
    if [ $exit_code -eq 0 ]; then
      echo "Tests passed on attempt $i."
      break
    fi

    # If maximum attempts are reached, exit the loop
    if [ $i -eq 25 ]; then
      echo "Maximum attempts reached. Tests still failing."
      exit 1
    fi
  done
else
  echo "Tests passed successfully."
fi