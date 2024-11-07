#!/usr/bin/env bash
set -e

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    # Generate only for files that don't have api in go_package
    if grep -q "option go_package" "$file" && grep -H -o -c 'option go_package.*chain-minimal/api' "$file" | grep -q ':0$'; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

echo "Generating pulsar proto code"
buf generate --template buf.gen.pulsar.yaml

cd ..

# Clean any existing generated files
rm -rf x/checkers/types/*.pb.go
mkdir -p x/checkers/types
mkdir -p api/checkers/v1

# Move generated files to correct locations
if [ -d "chain-minimal/x/checkers/types" ]; then
  mv chain-minimal/x/checkers/types/* x/checkers/types/
  rm -rf chain-minimal
fi

# Clean up any stray files
find proto -name "*.pb.go" -type f -delete