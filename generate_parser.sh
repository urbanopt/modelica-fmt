#!/bin/bash

docker build -t antlr4:latest -f thirdparty/Dockerfile-Antlr .

docker run -v "$(pwd)/thirdparty":/var/antlrResult \
        antlr4:latest \
        -Dlanguage=Go -o /var/antlrResult/parser /var/antlrResult/Modelica.g4
