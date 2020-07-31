#!/bin/bash
set -e

go build .

for file in ./examples/*.mo; do
    if [[ $file == *-out.mo ]]; then
        continue
    fi
    filename=$(basename -- $file)
    outfile="${filename%.*}-out.mo"
    ./modelica-fmt $file > ./examples/${outfile}
done
