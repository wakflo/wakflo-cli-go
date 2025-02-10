#!/bin/bash

# Define the binary location and target link
BINARY_PATH="../bin/wakflo"
TARGET_PATH="/usr/local/bin/wakflo"

# Function to check if the binary exists at the source location
check_binary_exists() {
  if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: CLI binary not found at '$BINARY_PATH'. Please build the binary first and try again."
    exit 1
  fi
}

# Function to handle creating the symbolic link
create_symlink() {
  # Check if the target link already exists
  if [ -L "$TARGET_PATH" ]; then
    echo "Symbolic link already exists at '$TARGET_PATH'."
    read -p "Do you want to overwrite it? (y/N): " response
    if [[ "$response" =~ ^[Yy]$ ]]; then
      echo "Removing existing symlink..."
      sudo rm "$TARGET_PATH"
    else
      echo "Aborting. No changes were made."
      exit 1
    fi
  elif [ -e "$TARGET_PATH" ]; then
    # If a file (not a symlink) exists at the target location
    echo "A file already exists at '$TARGET_PATH'. Please remove it manually before proceeding."
    exit 1
  fi

  # Create the symbolic link
  echo "Creating symbolic link from '$BINARY_PATH' to '$TARGET_PATH'..."
  sudo ln -s "$(realpath $BINARY_PATH)" "$TARGET_PATH"
  echo "Symbolic link created successfully. You can now run 'wakflo' globally."
}

# Main function
main() {
  echo "Linking Wakflo CLI binary..."
  check_binary_exists
  create_symlink
}

# Run the main function
main