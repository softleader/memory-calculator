#!/usr/bin/env bash

/tmp/memory-calculator

if [ -f "/tmp/.env" ]; then
  source /tmp/.env
fi

java $JAVA_OPTS -cp $( cat /app/jib-classpath-file ) $( cat /app/jib-main-class-file )
