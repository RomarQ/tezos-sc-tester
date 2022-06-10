#!/bin/sh

set -e

# Build
yarn build

# Package & Publish
cp package.json dist
cd dist
npm publish "$(npm pack)"
