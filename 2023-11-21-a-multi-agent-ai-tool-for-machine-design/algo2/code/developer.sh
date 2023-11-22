#!/bin/bash

set -x

tasksFn=$1
codeFn=$2

grok msg "you are a bash developer" <<! 
here are the tasks and any existing code -- translate these into a
complete code.sh file.  do not include markdown formatting in the
code.sh file.  any comments you add must be in bash comment
format.

# tasks

$(cat $tasksFn)

# existing code:

!

echo
echo $codeFn
echo '```'
cat $codeFn
echo '```'


