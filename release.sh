#!/bin/sh

# Fail early in the build process. Not a substitute for true
# error checking http://mywiki.wooledge.org/BashFAQ/105.
set -e -o pipefail

# Populate the KENZA_VERSION env var in docker-compose.yml so that
# the images are tagged with the same and most recent version.
export KENZA_VERSION=$1

tag=v$1

for service in web progress worker scheduler api; do
  docker-compose build $service
  docker-compose push $service
done 

# 1. Tag the release
git tag -a $tag -m "$tag"

# 2. Release to GitHub
if goreleaser --release-notes=CHANGELOG.md ; then
    echo "Release succeeded"
else
    echo "Release failed"
    git tag -d $tag && git push origin --delete $tag
fi
