#!/bin/bash 

set -x

touch code/stories.md

while true
do
    inotifywait -r -e modify code/*
    sleep 1
    make -C code
done
