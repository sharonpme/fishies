#!/bin/bash

$VERNEMQ_PATH/bin/vernemq console -noshell -noinput $@ &
pid=$!

trap "kill -s TERM $pid" SIGTERM
trap "kill -s INT $pid" SIGINT

wait $pid
