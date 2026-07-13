# Unreleased

* Add native support for formatting templated Modelica (`.mot`) files that embed Jinja
  constructs (`{{ ... }}`, `{% ... %}`, `{% raw %}`), ported from geojson-modelica-translator's
  `format_modelica_files.py`. `.mot` files are detected automatically (including in directory
  walks) and formatted via a substitute → format → reverse round trip. Adds a `-template` flag
  (currently `jinja`) to select the template dialect.

# Version 0.1

* Support formatting of Modelica files
