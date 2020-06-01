#! /bin/bash

docker build -t antlr4:latest -f Dockerfile-Antlr .

docker run -v "$(pwd)":/var/antlrResult \
        antlr4:latest \
        -Dlanguage=Go -o /var/antlrResult/parser /var/antlrResult/Modelica.g4
