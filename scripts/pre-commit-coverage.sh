#!/bin/bash

# Run tests and generate coverage
go test -coverprofile=coverage.out ./... > /dev/null 2>&1

# Extract total coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

if [ -z "$COVERAGE" ]; then
    echo "Failed to calculate coverage"
    exit 0
fi

# Determine color based on coverage percentage
COLOR="red"
if (( $(echo "$COVERAGE > 90" | bc -l) )); then
    COLOR="brightgreen"
elif (( $(echo "$COVERAGE > 70" | bc -l) )); then
    COLOR="yellow"
fi

# Update README.md badge
# Expected format: ![Coverage](https://img.shields.io/badge/coverage-VAL%25-COLOR)
SED_EXPR="s|coverage-[0-9.]*%25-[a-z]*|coverage-${COVERAGE}%25-${COLOR}|g"

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "$SED_EXPR" README.md
else
    sed -i "$SED_EXPR" README.md
fi

# Stage README.md if it was changed
git add README.md
