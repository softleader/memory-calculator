#!/usr/bin/env bash
set -e

DEFAULT_BIN_PATH="/usr/local/bin"
TMP_ENV="/tmp/.env"
JIB_CLASSPATH_FILE="/app/jib-classpath-file"
JIB_MAIN_CLASS_FILE="/app/jib-main-class-file"
JIB_JVM_FLAGS_FILE="/app/jib-jvm-flags-file"

read_jvm_flags() {
  if [ -f "$JIB_JVM_FLAGS_FILE" ]; then
    cat "$JIB_JVM_FLAGS_FILE"
  else
    echo ""
  fi
}

execute_memory_calculator() {
  if [ "${MEM_CALC_ENABLED}" = "false" ]; then
    echo "Memory calculator is disabled"
    return
  fi

  local bin_path=$1
  local debug=$2
  $bin_path/memory-calculator -o="$TMP_ENV" -v="$debug" || {
    echo "Memory calculator failed (version: $($bin_path/memory-calculator --version))";
    exit 1;
  }
  if [ -f "$TMP_ENV" ]; then
    source "$TMP_ENV"
  fi
}

DEBUG="${MEM_CALC_DEBUG:-false}"
BIN="${MEM_CALC_HOME:-$DEFAULT_BIN_PATH}"
JVM_FLAGS=$(read_jvm_flags)

execute_memory_calculator "$BIN" "$DEBUG"

[ "$DEBUG" = true ] && set -x
exec java $JVM_FLAGS -cp $(cat "$JIB_CLASSPATH_FILE") $(cat "$JIB_MAIN_CLASS_FILE") "$@"
