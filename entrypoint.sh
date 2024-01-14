#!/usr/bin/env bash
set -e

DEFAULT_BIN_PATH="/usr/local/bin"
TMP_ENV="/tmp/.env"
JIB_CLASSPATH_FILE="/app/jib-classpath-file"
JIB_MAIN_CLASS_FILE="/app/jib-main-class-file"
JIB_JVM_FLAGS_FILE="/app/jib-jvm-flags-file"

DEBUG=${DEBUG:-false}

get_memory_calculator_bin() {
  if [ -n "$MEMORY_CALCULATOR_HOME" ]; then
    echo "$MEMORY_CALCULATOR_HOME"
  else
    echo "$DEFAULT_BIN_PATH"
  fi
}

read_jvm_flags() {
  if [ -f "$JIB_JVM_FLAGS_FILE" ]; then
    cat "$JIB_JVM_FLAGS_FILE"
  else
    echo ""
  fi
}

execute_java_app() {
  local jvm_flags=$1
  local args=$2

  [ "$DEBUG" = true ] && set -x

  if [ -f "$TMP_ENV" ]; then
    source "$TMP_ENV"
    exec java $jvm_flags $args -cp $(cat "$JIB_CLASSPATH_FILE") $(cat "$JIB_MAIN_CLASS_FILE")
  else
    exec java $jvm_flags $args $JAVA_OPTS -cp $(cat "$JIB_CLASSPATH_FILE") $(cat "$JIB_MAIN_CLASS_FILE")
  fi

  [ "$DEBUG" = true ] && set +x
}

BIN=$(get_memory_calculator_bin)
JVM_FLAGS=$(read_jvm_flags)
ARGS="$@"

$BIN/memory-calculator -o $TMP_ENV || { echo "Memory calculator failed"; exit 1; }

execute_java_app "$JVM_FLAGS" "$ARGS"
