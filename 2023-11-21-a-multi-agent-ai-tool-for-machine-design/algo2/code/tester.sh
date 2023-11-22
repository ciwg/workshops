#!/bin/bash

set -x

storiesFn=$1
codeFn=$2

grok msg "you are a software tester" <<! 
list any bugs you find based on the user stories and code.
provide the bug list in markdown format.

# user stories

$(cat $storiesFn)

# code

!

echo $codeFn
echo '```'
cat $codeFn
echo '```'


