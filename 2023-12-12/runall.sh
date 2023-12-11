#!/bin/bash

term="/usr/bin/uxterm -ms red -bg black -fg white -ls -sb -sl 9999"

agent=$1

. ./config.sh

if [ -n "$agent" ]
then
    # run a single agent
    $term -e "./cycle.sh $agent" &
    exit 0
else
    # run all agents, each in their own xterm
    for agent in $agents
    do
        $term -e "./cycle.sh $agent" &
    done
fi
