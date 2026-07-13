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

# Templated Modelica (.mot) examples. Some templates cannot be formatted even
# with preprocessing (see SKIP_FILES in geojson-modelica-translator); they are
# kept as fixtures to exercise graceful failure and have no *-out.mot output.
skip_mot=(gmt-district-energy-system.mot gmt-hptrio-variable-dist.mot gmt-run-spawn-building.mot)

for file in ./examples/*.mot; do
    if [[ $file == *-out.mot ]]; then
        continue
    fi
    filename=$(basename -- $file)
    skip=false
    for s in "${skip_mot[@]}"; do
        if [[ $filename == "$s" ]]; then
            skip=true
            break
        fi
    done
    if [[ $skip == true ]]; then
        continue
    fi
    outfile="${filename%.*}-out.mot"
    ./modelica-fmt $file > ./examples/${outfile}
done
