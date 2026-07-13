#!/bin/bash
set -e

go build .

# Default formatting (no extra flags) for every source file
for file in ./examples/*.mo; do
    if [[ $file == *-out.mo ]]; then
        continue
    fi
    filename=$(basename -- $file)
    outfile="${filename%.*}-out.mo"
    ./modelica-fmt $file > ./examples/${outfile}
done

# Variants exercised by specific test cases in modelicafmt_test.go
./modelica-fmt -line-length 80 ./examples/gmt-building.mo > ./examples/gmt-building-80-out.mo
./modelica-fmt -extra-padding ./examples/gmt-building.mo > ./examples/gmt-building-empty-lines-out.mo
./modelica-fmt -format-html ./examples/gmt-building.mo > ./examples/gmt-building-html-out.mo
./modelica-fmt -format-html ./examples/html-annotation.mo > ./examples/html-annotation-html-out.mo
