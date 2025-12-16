#!/bin/sh

# set -e: Exit immediately if a command exits with a non-zero status.
# set -u: Treat unset variables as an error when substituting.
set -eu

# --- Constants ---
readonly INSTALL_BASE_DIR="/opt/memory-calculator"
readonly BIN_DIR="${INSTALL_BASE_DIR}/bin"
readonly SYMLINK_PATH="/usr/local/bin/memory-calculator"
readonly TMP_DIR="/tmp"
readonly RELEASES_URL="https://github.com/softleader/memory-calculator/releases/latest/download"
readonly SUPPORTED_ARCHS="amd64 arm64"

# --- Logging ---
# Prints an error message to stderr and exits.
# Usage: error "Something went wrong"
error() {
  echo "Error: $*" >&2
  exit 1
}

# --- Core Functions ---

check_privileges() {
  if [ "$(id -u)" -ne 0 ]; then
    error "This script requires root or sudo privileges. Please try: sudo $0"
  fi
}

check_dependencies() {
  for cmd in curl unzip; do
    if ! command -v "$cmd" > /dev/null 2>&1; then
      error "This script requires '$cmd', but it is not installed."
    fi
  done
}

# Downloads and extracts a release for a specific architecture.
# Usage: setup_arch "amd64"
setup_arch() {
  arch=$1
  zip_file="${TMP_DIR}/linux-${arch}.zip"
  url="${RELEASES_URL}/linux-${arch}.zip"
  extract_dir="${BIN_DIR}/${arch}"
  binary_path="${extract_dir}/memory-calculator"

  echo "Setting up ${arch} architecture..."
  
  mkdir -p "$extract_dir"
  curl -L -s -o "$zip_file" "$url"
  unzip -o -q "$zip_file" -d "$extract_dir"
  rm -f "$zip_file"
  
  # Set executable permissions for the binary
  chmod 755 "$binary_path"
}

create_wrapper_script() {
  wrapper_path="${BIN_DIR}/memory-calculator"
  
  # Using a here document to create the script content
  cat << 'EOF' > "$wrapper_path"
#!/bin/sh
set -e

# Get the directory where this script is located
SCRIPT_DIR=$(dirname "$(readlink -f "$0" 2>/dev/null || realpath "$0" 2>/dev/null || echo "$0")")

# Determine CPU architecture
ARCH=$(uname -m)

case "$ARCH" in
  "x86_64")
    BIN_PATH="${SCRIPT_DIR}/amd64/memory-calculator"
    ;;
  "aarch64")
    BIN_PATH="${SCRIPT_DIR}/arm64/memory-calculator"
    ;;
  *)
    echo "Error: Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

if [ ! -f "$BIN_PATH" ]; then
    echo "Error: Binary not found for architecture $ARCH at $BIN_PATH" >&2
    exit 1
fi

# If MEM_CALC_DEBUG is true, print the binary being used to stderr.
if [ "${MEM_CALC_DEBUG:-false}" = "true" ]; then
  echo "DEBUG: Using binary: $BIN_PATH" >&2
fi

# Execute the target binary, passing all arguments
exec "$BIN_PATH" "$@"
EOF

  chmod +x "$wrapper_path"
}

create_symlink() {
  ln -sf "${BIN_DIR}/memory-calculator" "$SYMLINK_PATH"
}

# Copies entrypoint.sh if the --entrypoint flag was used.
# It iterates through supported architectures and copies from the first one found.
copy_entrypoint() {
  if [ -n "$ENTRYPOINT_TARGET_PATH" ]; then
    copied_successfully=0 # Flag to track if copy was successful
    for arch in $SUPPORTED_ARCHS; do
      source_entrypoint="${BIN_DIR}/${arch}/entrypoint.sh"
      if [ -f "$source_entrypoint" ]; then
        echo "Copying entrypoint.sh to $ENTRYPOINT_TARGET_PATH"
        mkdir -p "$ENTRYPOINT_TARGET_PATH"
        cp "$source_entrypoint" "$ENTRYPOINT_TARGET_PATH/"
        copied_successfully=1
        break # Stop after the first successful copy
      fi
    done

    if [ "$copied_successfully" -eq 0 ]; then
      echo "Warning: entrypoint.sh not found in any supported archive, cannot copy." >&2
    fi
  fi
}

# --- Main Logic ---
main() {
  check_privileges
  check_dependencies

  for arch in $SUPPORTED_ARCHS; do
    setup_arch "$arch"
  done

  create_wrapper_script
  create_symlink

  copy_entrypoint

  echo "Installation complete! 'memory-calculator' is ready to use."
}

# --- Script Entrypoint ---
ENTRYPOINT_TARGET_PATH=""
for arg in "$@"; do
  case "$arg" in
    --entrypoint=*)
      ENTRYPOINT_TARGET_PATH="${arg#*=}"
      ;;
  esac
done

main
