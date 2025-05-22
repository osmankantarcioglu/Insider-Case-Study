#!/bin/bash

# Set environment variables
export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Run go build with verbose output
echo "Running go build with verbose output..."
go build -v -a -installsuffix cgo -o main ./cmd

# Check if build was successful
if [ $? -eq 0 ]; then
  echo "Build successful!"
else
  echo "Build failed. Check the error message above."
  # Try to get more information about dependencies
  echo ""
  echo "Checking dependencies..."
  go mod tidy -v
  go list -m all
fi 