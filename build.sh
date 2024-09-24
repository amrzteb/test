#!/bin/bash

# Set Go environment variables for cross-compilation
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# Set output directory and project name
BUILD_DIR="build"
PROJECT_NAME="test5000"

# Remote server details
REMOTE_SERVER="test"
REMOTE_PATH="/root/opt/p711co/"

# Create build directory if it doesn't exist
mkdir -p $BUILD_DIR

# Build the project
echo "Building $PROJECT_NAME for $GOOS/$GOARCH..."
go build -o $BUILD_DIR/$PROJECT_NAME .

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build succeeded."

    # Copy the build to the remote server
    echo "Copying $PROJECT_NAME to $REMOTE_SERVER:$REMOTE_PATH"
    scp $BUILD_DIR/$PROJECT_NAME $REMOTE_SERVER:$REMOTE_PATH

    # Check if the SCP was successful
    if [ $? -eq 0 ]; then
        echo "Copy succeeded."
    else
        echo "Copy failed."
        exit 1
    fi

else
    echo "Build failed."
    exit 1
fi
