# Unreleased

* Add `--format-html` (BETA): pretty-print HTML content embedded in Modelica annotation
  strings, e.g. `Documentation(info="<html>...</html>")` and `revisions="..."`. Opt-in and
  conservative: only whitespace/indentation changes (tags, attributes, entities and text are
  preserved exactly), block-level elements are indented to their nesting depth while
  inline/phrasing elements (e.g. `<b>`, `<a>`, `<code>`, `<br/>`) and their surrounding text
  stay condensed on one line, whitespace-sensitive elements (`<pre>`, `<textarea>`,
  `<script>`, `<style>`) are left verbatim, and malformed/unbalanced HTML is reported as a
  clear error (leaving the file unchanged and exiting non-zero) rather than being silently
  emitted unformatted.

# Version 0.1

* Support formatting of Modelica files
