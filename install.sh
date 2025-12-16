#!/bin/sh

# set -e: Exit immediately if a command exits with a non-zero status.
# set -u: Treat unset variables as an error when substituting.
set -eu

# --- Constants ---
readonly GITHUB_REPO_URL="https://github.com/softleader/memory-calculator"
readonly INSTALL_BIN_PATH="/usr/local/bin/memory-calculator"

# --- Script Variables (set by parsing args or detection) ---
VERSION_TAG="latest"
ENTRYPOINT_TARGET_PATH=""
OS_OVERRIDE=""
ARCH_OVERRIDE=""
PLATFORM=""
TMP_DIR=""

# --- Logging ---
error() {
  echo "Error: $*" >&2
  exit 1
}

# --- Core Functions ---

parse_args() {
  for arg in "$@"; do
    case "$arg" in
      --version=*)    VERSION_TAG="${arg#*=}" ;;
      --entrypoint=*) ENTRYPOINT_TARGET_PATH="${arg#*=}" ;;
      --os=*)         OS_OVERRIDE="${arg#*=}" ;;
      --arch=*)       ARCH_OVERRIDE="${arg#*=}" ;;
    esac
  done
}

check_privileges() {
  local install_dir
  install_dir=$(dirname "${INSTALL_BIN_PATH}")
  if ! [ -d "${install_dir}" ]; then
    error "Installation directory '${install_dir}' does not exist. Please create it or use sudo."
  elif ! [ -w "${install_dir}" ]; then
    error "You do not have write permissions for ${install_dir}. Please run with sudo."
  fi

  if [ -n "$ENTRYPOINT_TARGET_PATH" ]; then
    local entrypoint_dir
    entrypoint_dir=$(dirname "${ENTRYPOINT_TARGET_PATH}")

    # The script will use 'mkdir -p', so we check permissions on the first existing parent.
    local dir_to_check="${entrypoint_dir}"
    while [ "${dir_to_check}" != "." ] && [ "${dir_to_check}" != "/" ] && ! [ -d "${dir_to_check}" ]; do
      dir_to_check=$(dirname "${dir_to_check}")
    done

    if ! [ -w "${dir_to_check}" ]; then
      error "You do not have write permissions for '${dir_to_check}' to create '${entrypoint_dir}'. Please run with sudo."
    fi
  fi
}

check_dependencies() {
  for cmd in curl unzip uname tr; do
    if ! command -v "$cmd" > /dev/null 2>&1; then
      error "This script requires '$cmd', but it is not installed."
    fi
  done
}

determine_platform() {
  local os="${OS_OVERRIDE}"
  local arch="${ARCH_OVERRIDE}"

  if [ -z "$os" ]; then
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
  else
    echo "Using specified OS: ${os}"
  fi

  if [ -z "$arch" ]; then
    arch=$(uname -m)
  else
    echo "Using specified Arch: ${arch}"
  fi

  # Normalize OS
  case "$os" in
    linux) os="linux" ;;
    darwin) os="darwin" ;;
    *) error "Unsupported OS: ${os}" ;;
  esac

  # Normalize Arch
  case "$arch" in
    x86_64) arch="amd64" ;;
    aarch64 | arm64) arch="arm64" ;;
    amd64) arch="amd64" ;; # Allow explicit 'amd64'
    *) error "Unsupported architecture: ${arch}" ;;
  esac

  PLATFORM="${os}-${arch}"
  echo "Determined platform: ${PLATFORM}"
}

download_and_extract() {
  local zip_name="${PLATFORM}.zip"
  local download_url

  if [ "$VERSION_TAG" = "latest" ]; then
    download_url="${GITHUB_REPO_URL}/releases/latest/download/${zip_name}"
  else
    download_url="${GITHUB_REPO_URL}/releases/download/${VERSION_TAG}/${zip_name}"
  fi

  echo "Downloading from: ${download_url}"

  local zip_file="${TMP_DIR}/${zip_name}"
  curl -L -s -o "${zip_file}" "${download_url}"

  if ! unzip -t "${zip_file}" > /dev/null 2>&1; then
    error "Download failed. Check if version '${VERSION_TAG}' and platform '${PLATFORM}' are valid. URL: ${download_url}"
  fi

  unzip -o -q "${zip_file}" -d "${TMP_DIR}"
}

install_files() {
  echo "Installing 'memory-calculator' to ${INSTALL_BIN_PATH}..."
  mv "${TMP_DIR}/memory-calculator" "${INSTALL_BIN_PATH}"
  chmod 755 "${INSTALL_BIN_PATH}"

  if [ -n "$ENTRYPOINT_TARGET_PATH" ]; then
    local source_entrypoint="${TMP_DIR}/entrypoint.sh"
    if [ -f "$source_entrypoint" ]; then
      echo "Copying entrypoint.sh to ${ENTRYPOINT_TARGET_PATH}"
      mkdir -p "$(dirname "${ENTRYPOINT_TARGET_PATH}")"
      cp "${source_entrypoint}" "${ENTRYPOINT_TARGET_PATH}"
      chmod 755 "${ENTRYPOINT_TARGET_PATH}"
    else
      echo "Warning: entrypoint.sh not found in archive, cannot copy." >&2
    fi
  fi
}

# --- Main Logic ---
main() {
  parse_args "$@"
  check_privileges
  check_dependencies
  determine_platform

  # Setup temp dir and cleanup
  TMP_DIR="/tmp/mc-install-$$"
  mkdir -p "${TMP_DIR}"
  trap 'rm -rf "${TMP_DIR}"' EXIT

  download_and_extract
  install_files

  echo "Installation complete! 'memory-calculator' version ${VERSION_TAG} is ready to use."
}

# --- Script Entrypoint ---
main "$@"

