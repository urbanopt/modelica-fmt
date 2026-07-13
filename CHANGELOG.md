# Unreleased

* Add `--format-html` (BETA): pretty-print HTML content embedded in Modelica annotation
  strings, e.g. `Documentation(info="<html>...</html>")` and `revisions="..."`. Opt-in and
  conservative: only whitespace/indentation changes (tags, attributes, entities and text are
  preserved exactly), whitespace-sensitive elements (`<pre>`, `<textarea>`, `<script>`,
  `<style>`) are left verbatim, and malformed/unbalanced HTML is left untouched.

# Version 0.1

* Support formatting of Modelica files
