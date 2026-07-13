# Unreleased

* Add native support for formatting templated Modelica (`.mot`) files that embed Jinja
  constructs (`{{ ... }}`, `{% ... %}`, `{% raw %}`), ported from geojson-modelica-translator's
  `format_modelica_files.py`. `.mot` files are detected automatically (including in directory
  walks) and formatted via a substitute → format → reverse round trip. Adds a `-template` flag
  (currently `jinja`) to select the template dialect.
* Refuse to overwrite a file with empty output: if a non-empty input produces no
  non-whitespace output (e.g. a non-Modelica file matched as an empty definition), the
  formatter now returns an error and leaves the file unchanged instead of emptying it.

# Version 0.1

* Support formatting of Modelica files
