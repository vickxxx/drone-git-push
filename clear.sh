#!/bin/sh

git filter-branch --force --index-filter "git rm -rf --cached --ignore-unmatch $1" --tag-name-filter cat -- --all
rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now
git gc --aggressive --prune=now
git push --force
git remote prune origin