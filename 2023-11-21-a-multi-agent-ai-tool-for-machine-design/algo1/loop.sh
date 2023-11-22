#!/bin/bash 

set -x

while true
do

echo enter new user story
read story
echo "- $story" >> stories.md

grok msg "you are a project manager" <<! > tasks.md
here are the user stories and bugs -- translate these into tasks
in markdown format:

# user stories

$(cat stories.md)

# bugs

$(cat bugs.md)
!

grok msg "you are a bash developer" <<! > code.sh
here are the tasks and any existing code -- translate these into a
complete code.sh file.  do not include markdown formatting in the
code.sh file.  any comments you add must be in bash comment
format.

# tasks

$(cat tasks.md)

# existing code:

$(cat code.sh)
!

grok msg "you are a software tester" <<! > bugs.md
list any bugs you find based on the user stories and code.
provide the bug list in markdown format.

# user stories

$(cat stories.md)

# code

$(cat code.sh)
!

done
