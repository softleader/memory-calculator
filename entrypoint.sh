#!/usr/bin/env bash

/tmp/memory-calculator -o /tmp/.env

if [ -f "/tmp/.env" ]; then
  source /tmp/.env
  java -cp $( cat /app/jib-classpath-file ) $( cat /app/jib-main-class-file )
else
  java $JAVA_OPTS -cp $( cat /app/jib-classpath-file ) $( cat /app/jib-main-class-file )
fi

