#!/bin/bash -e

repo_root=$(git rev-parse --show-toplevel)
current_branch=$(git rev-parse --abbrev-ref HEAD)
twig=$(basename $current_branch)

usage() {
    echo "Usage: $0 [-cFn] [-b BASE_BRANCH] [OTHER_BRANCH]

    With no arguments, compare $current_branch with other branches named "'*'"/$twig.

    If OTHER_BRANCH is given, do a manual merge of OTHER_BRANCH into $current_branch.
    
    -F force run even if not on a $USER/ branch
    -b create new branch named $USER/... based on BASE_BRANCH 
    -n no automatic push after commit
    -c commit existing uncommitted changes
    "
    exit 1
}

ck_clean() {
    if [ -n "$(git status --porcelain)" ]
    then 
        set +x
        echo "you have uncommitted changes"
        if [ -z "$commit" ] 
        then
            exit 1
        fi
        git diff HEAD
        git commit -a
        if [ -z "$nopush" ]
        then
            git push
        fi
    fi
}

ck_user() {
    if [ -z "$force" ] && ! echo $current_branch | egrep -q "^$USER/"
    then
        echo "$0: you aren't on a '$USER/' branch"
        echo
        usage
    fi
}

unset force
unset base_branch
unset nopush
unset commit
while getopts Fchnb: opt 
do
    case $opt in
        F)
            force=y
            ;;
        b)
            force=y
            base_branch=$OPTARG
            ;;
        n)
            nopush=y
            ;;
        c)
            commit=y
            ck_clean
            ;;
        h)
            usage
            exit 0
            ;;
        *)
            echo
            usage
            ;;
    esac
done
shift $(( OPTIND-1 ))

other_branch=$1

git fetch

if [ -n "$base_branch" ]
then
    ck_clean
    twig=$(basename $base_branch)
    new_branch=$USER/$twig
    if echo  $base_branch | grep -q remotes/
    then
        git checkout -b $twig $base_branch
        base_branch=$twig
    fi
    git checkout -b $new_branch $base_branch
    git push --set-upstream origin $new_branch
elif [ -z "$other_branch" ]
then
    ck_user

    echo
    echo
    echo "Related branches and their diffs (if any):"
    echo

    for b in $(git branch -a | tr -d '*' | egrep "/$twig$")
    do 
        if [ "$b" == "$current_branch" ]
        then
            continue
        fi
        ahead=$(git diff --shortstat ...${b})
        behind=$(git diff --shortstat ${b}...)
        if [ -n "$ahead" ]
        then
            # printf "%40s is ahead:\n" $b 
            # git diff --stat ...${b} 
            printf "%40s is ahead: %s\n" $b "$ahead"
        elif [ -n "$behind" ]
        then
            printf "%40s is behind: %s\n" $b "$behind"
        else
            printf "%40s is synced\n" $b 
        fi

        # also see https://git.seveas.net/previewing-a-merge-result.html for a
        # different approach with more detailed diffs:
        # git merge-tree $(git merge-base HEAD develop) HEAD develop
    done

    exit 0
else
    ck_user
    ck_clean

    msg=/tmp/mob-consensus.msg

    printf "mob-consensus merge from $other_branch onto $current_branch\n\n" > $msg
    # printf "\n\n" > $msg

    set -x
    git log ..$other_branch --pretty=format:"Co-authored-by: %an <%ae>" \
        | grep -v $(git config --get user.email) \
        | sort -u >> $msg

    if ! git merge --no-commit --no-ff $other_branch
    then
        # deal with merge conflicts
        git mergetool -t vimdiff 
    fi

    # commit during merge will ignore `-t`, so we replace the
    # default commit message here as well
    # XXX might instead want to use -m or -F flag on the above merge 
    cp $msg $repo_root/.git/MERGE_MSG

    # review all changes
    git difftool -t vimdiff HEAD

    if ! git commit -e -F $msg
    then
        set +x
        echo "don't forget to push"
        exit 1
    fi

    if [ -n "$nopush" ]
    then
        set +x
        echo "skipping automatic push -- don't forget to push later"
    else
        git push
    fi
fi


