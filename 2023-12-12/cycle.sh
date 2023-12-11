#!/bin/bash

# go to sleep if the -e flag causes us to exit
trap "sleep 99999" ERR

set -ex

. ./config.sh

# get absolute path to fixdiff
abspath() {
    echo $(cd $(dirname $1); pwd)/$(basename $1)
}
fixdiff=$(abspath ../fixdiff/fixdiff)

merge() {
    agent=$1
    other=$2
    cd $agent
    pwd

    git fetch

    # merge, preferring other's changes
    git merge --no-commit -Xtheirs $other

    # save the merge diff and undo the merge
    git diff --staged > .diff
    if ! git merge --abort 
    then 
        # we will also get here if the merge was a fast-forward
        git push
        cd -
        return 0
    fi

    # check sanity
    cat .diff | grok chat .chat -s "You are a $agent.  Do the given changes make sense?  Your answer must match the following regex: '^answer=(yes|no)$'" > .answer

    # if answer=yes then do the merge for real
    if [ $(cat .answer) == "answer=yes" ]
    then
        git merge -Xtheirs -m "merge with $other" $other
    fi
    git push

    # update grokker
    grok add *

    cd -
}

edit() {
    agent=$1
    cd $agent
    pwd

    # get the last commit we looked at
    if [ -f .lastcommit ]
    then
        lastcommit=$(cat .lastcommit)
    else
        # start at the beginning of time
        lastcommit=$(git rev-list --max-parents=0 HEAD)
    fi
    # get the files that have changed since then
    git log --pretty=format: --name-only $lastcommit..HEAD > .changed

    # if there are no changes, then return
    if [ ! -s .changed ]
    then
        cd -
        return 0
    fi

    inputs=$(cat .changed)
    # remove blank lines
    inputs=$(echo "$inputs" | sed '/^$/d')
    # replace newlines with commas
    inputs=$(echo "$inputs" | tr '\n' ',' | sed 's/,$//')
    # get output for this agent from agent_outputs array in config.sh
    output=${agent_outputs[$agent]}

    while true
    do
        if time grok chat .chat -s "You are a $agent." --input-files="$inputs" --output-files="$output" -m "Edit $output based on the changes in $inputs."
        then
            break
        fi
        echo retrying...
    done

    # update grokker
    grok add *

    # commit
    git add *
    grok commit > .commitmsg
    git commit -F .commitmsg
    git push

    # save the commit we just looked at
    git rev-parse HEAD > .lastcommit

    cd -
}


agent=$1

cd agents
while true
do
    # merge with main
    merge $agent origin/main

    # merge with other branches
    for other in *
    do
        if [ $agent == $other ]
        then
            continue
        fi
        merge $agent origin/$other/main
    done

    # make a change
    edit $agent
    sleep 10
done

