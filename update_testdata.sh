#!/bin/bash
set -e

go build .

for file in ./internal/format/testdata/*.mo; do
    if [[ $file == *-out.mo ]]; then
        continue
    fi
    filename=$(basename -- $file)
    outfile="${filename%.*}-out.mo"
    ./modelica-fmt $file > ./internal/format/testdata/${outfile}
done

# Golden files for non-default formatter configurations used by
# internal/format/modelicafmt_test.go.
./modelica-fmt -wrap-arrays ./internal/format/testdata/example-arrays.mo > ./internal/format/testdata/example-arrays-wrapped-out.mo
./modelica-fmt -wrap-arrays ./internal/format/testdata/example-arrays-nd.mo > ./internal/format/testdata/example-arrays-nd-out.mo
./modelica-fmt -line-length 80 ./internal/format/testdata/gmt-building.mo > ./internal/format/testdata/gmt-building-80-out.mo
./modelica-fmt -extra-padding ./internal/format/testdata/gmt-building.mo > ./internal/format/testdata/gmt-building-empty-lines-out.mo
./modelica-fmt -format-html ./internal/format/testdata/html-annotation.mo > ./internal/format/testdata/html-annotation-html-out.mo
./modelica-fmt -format-html ./internal/format/testdata/gmt-building.mo > ./internal/format/testdata/gmt-building-html-out.mo

for file in ./internal/format/testdata/*.mot; do
    if [[ $file == *-out.mot ]]; then
        continue
    fi
    filename=$(basename -- $file)
    outfile="${filename%.*}-out.mot"
    ./modelica-fmt $file > ./internal/format/testdata/${outfile}
done

for file in ./internal/format/testdata/*.mopt; do
    if [[ ! -e $file ]]; then
        continue
    fi
    if [[ $file == *-out.mopt ]]; then
        continue
    fi
    filename=$(basename -- $file)
    outfile="${filename%.*}-out.mopt"
    ./modelica-fmt $file > ./internal/format/testdata/${outfile}
done
