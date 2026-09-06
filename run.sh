#!/bin/bash -e

pushd www/; python3 -m http.server; popd
