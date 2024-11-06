#!/usr/bin/env bash

set -e

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    # this regex checks if a proto file has its go_package set to github.com/chain-minimal/x/checkers/types
    if grep -q "option go_package" "$file" && grep -H -o -c 'option go_package.*github.com/chain-minimal/x/checkers/types' "$file" | grep -q ':0$'; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

echo "Generating pulsar proto code"
buf generate --template buf.gen.pulsar.yaml

cd ..

# Create necessary directories
mkdir -p x/checkers/types
mkdir -p api/checkers/v1
mkdir -p api/checkers/module/v1

# Copy generated files to appropriate locations
# For types, query, tx
cp -r github.com/chain-minimal/x/checkers/types/* ./x/checkers/types/ 2>/dev/null || :
# For module
cp -r github.com/chain-minimal/x/checkers/module/* ./x/checkers/module/ 2>/dev/null || :

# Copy to API directory for external consumption
cp -r x/checkers/types/* ./api/checkers/v1/ 2>/dev/null || :
cp -r x/checkers/module/* ./api/checkers/module/v1/ 2>/dev/null || :

# Clean up
rm -rf github.com chain-minimal