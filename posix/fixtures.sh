#!/bin/sh
set -e
set -x

rm -rf /tmp/remote/greeting
mkdir -p /tmp/remote/greeting
pushd /tmp/remote/greeting

git init

echo -n "hi world" > hello.txt
git add hello.txt
git commit -m "say hi"
git tag v1.0.0

echo -n "hello world" > hello.txt
git add hello.txt
git commit -m "say hello"
git tag v1.1.0

git checkout -b fr
echo -n "salut monde" > hello.txt
git add hello.txt
git commit -m "say hello in french"
git tag v2.0.0

echo -n "bonjour monde" > hello.txt
git add hello.txt
git commit -m "say hello en francais"
git tag v2.1.0

git checkout master

popd
tar -cvf fixtures.tar /tmp/remote/greeting
