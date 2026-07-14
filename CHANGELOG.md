# Unreleased

* Add native support for formatting templated Modelica (`.mot`/`.mopt`) files that embed Jinja
  constructs (`{{ ... }}`, `{% ... %}`, `{% raw %}`), ported from geojson-modelica-translator's
  `format_modelica_files.py`. `.mot` and `.mopt` files are detected automatically (including in directory
  walks) and formatted via a substitute → format → reverse round trip. Adds a `-template` flag
  (currently `jinja`) to select the template dialect.
* Refuse to overwrite a file with empty output: if a non-empty input produces no
  non-whitespace output (e.g., a non-Modelica file matched as an empty definition), the
  formatter now returns an error and leaves the file unchanged instead of emptying it.
* Add `-wrap-arrays` option to format multidimensional arrays (`{...}`) across
  multiple lines outside of annotations (issue #29)
* Add `-format-html`: pretty-print HTML content embedded in Modelica annotation
  strings, e.g. `Documentation(info="<html>...</html>")` and `revisions="..."`. Opt-in and
  conservative: only whitespace/indentation changes (tags, attributes, entities and text are
  preserved exactly), block-level elements are indented to their nesting depth while
  paragraphs with inline-only content, plain `<h4>` headings, and inline/phrasing elements
  (e.g., `<b>`, `<a>`, `<code>`, `<br/>`) stay condensed on one line, whitespace-sensitive
  elements (`<pre>`, `<textarea>`, `<script>`, `<style>`) are left verbatim, and
  malformed/unbalanced HTML is reported as a clear error (leaving the file unchanged and
  exiting non-zero) rather than being silently emitted unformatted (issue #40).

# Version 0.1

* Support formatting of Modelica files
