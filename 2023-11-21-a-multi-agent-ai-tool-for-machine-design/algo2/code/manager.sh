#!/bin/bash

set -x

storiesFn=$1
bugsFn=$2

grok msg "you are a project manager" <<! 
here are the user stories and bugs -- translate these into tasks
in markdown format:

# user stories

$(cat $storiesFn)

# bugs

$(cat $bugsFn)
!

