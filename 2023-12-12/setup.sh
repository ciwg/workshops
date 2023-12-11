#!/bin/bash

set -ex

# clean
rm -rf origin.git agents control

. ./config.sh

# set up git
git init --bare origin.git

# set up control repo
git clone origin.git control
cd control
# put requirements on main branch
cp ../requirements.md .
# touch outputs
for output in ${agent_outputs[@]}
do
    touch $output
done
git add *
git commit -m "set up control repo"
git push
grok init
grok add *
cd -

# set up agents
mkdir agents
cd agents
for agent in $agents
do 
    git clone ../origin.git $agent
    cd $agent
    git checkout -b $agent/main
    git push --set-upstream origin $agent/main
    grok init
    grok add *
    cd -
done
cd -


