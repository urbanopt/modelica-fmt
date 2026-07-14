# Changelog

All notable changes to this project are documented in this file. The format
is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and each
released version corresponds to a tagged
[GitHub release](https://github.com/urbanopt/modelica-fmt/releases).

## [Unreleased]

## [v0.3.1] - 2026-07-14

* Add `-format-html`: pretty-print HTML content embedded in Modelica annotation
  strings, e.g. `Documentation(info="<html>...</html>")` and `revisions="..."`. Opt-in and
  conservative: only whitespace/indentation changes (tags, attributes, entities and text are
  preserved exactly), block-level elements are indented to their nesting depth while
  paragraphs with inline-only content, plain `<h4>` headings, and inline/phrasing elements
  (e.g., `<b>`, `<a>`, `<code>`, `<br/>`) stay condensed on one line, whitespace-sensitive
  elements (`<pre>`, `<textarea>`, `<script>`, `<style>`) are left verbatim, and
  malformed/unbalanced HTML is reported as a clear error (leaving the file unchanged and
  exiting non-zero) rather than being silently emitted unformatted (issue #40).
* Tweak default formatting for readability (issue #26):
  * Keep a function call with a single argument on one line instead of breaking
    the argument onto its own line, e.g. `pre(x)`, `der(y)`, `sin(z)`,
    `sum(a .* b)`. Calls with two or more arguments are unchanged.
  * Insert a blank line before `equation`/`initial equation`,
    `algorithm`/`initial algorithm`, and `public`/`protected` section headers to
    separate them from the preceding declarations. The blank line is not
    duplicated when `-extra-padding` already adds one.
  * Collapse the stray empty line often left between the final `</ul>` and
    `</html>` of a `revisions` docstring (`</ul>\n\n</html>` becomes
    `</ul>\n</html>`).
* Fix `.mot`/`.mopt` formatting collapsing the space in a single-element array
  literal whose element is a Jinja expression (e.g. `{ {{ data[...] }} }`),
  which produced an ambiguous `{{{ ... }}}` triple-brace sequence. A single
  disambiguating space is now preserved whenever a restored Jinja expression's
  own `{{`/`}}` delimiters would otherwise merge with an adjacent Modelica
  array brace.

## [v0.3.1-pr1] - 2026-07-13

* Broaden `.mot` template support to handle more difficult real-world files
  (additional Jinja constructs and edge cases).

## [v0.3.0] - 2026-07-13

* Add native support for formatting templated Modelica (`.mot`/`.mopt`) files that embed Jinja
  constructs (`{{ ... }}`, `{% ... %}`, `{% raw %}`), ported from geojson-modelica-translator's
  `format_modelica_files.py`. `.mot` and `.mopt` files are detected automatically (including in directory
  walks) and formatted via a substitute → format → reverse round trip. Adds a `-template` flag
  (currently `jinja`) to select the template dialect.
* Refuse to overwrite a file with empty output: if a non-empty input produces no
  non-whitespace output (e.g., a non-Modelica file matched as an empty definition), the
  formatter now returns an error and leaves the file unchanged instead of emptying it.
* Add `-wrap-arrays` option to format multidimensional arrays (`{...}`) across
  multiple lines outside of annotations (issue #29).

## [v0.2-pr.2] - 2021-02-19

* Add a `-line-length` CLI flag and break lines longer than 80 characters.

## [v0.2-pr.1] - 2021-01-28

* Add a `-v`/`-version` flag to display version info.
* Indent arguments of external function calls.
* Show CLI usage when no arguments are provided.
* Move CI from Travis to GitHub Actions.

## [v0.1] - 2020-09-30

* Support formatting of Modelica files.

## [v0.1-pr.3] - 2020-08-17

* Fix handling of comments trailing at the end of a file.

## [v0.1-pr.2] - 2020-07-31

* Always indent modifications and arguments.
* Allow indentation in model annotation vectors.
* Remove extraneous spacing: no space after commas, around array constructors,
  or around the `+` operator.
* Add basic tests for example files.

## [v0.1-pr.1] - 2020-06-22

* Initial formatter implementation for Modelica (`.mo`) files: annotation,
  if-expression, and array-constructor formatting rules.
* Add an option to write the formatted result back to the source file.
* Add rudimentary comment handling.
* Rename the project/binary to `modelica-fmt`/`modelicafmt`.
* Set up Travis CI and GoReleaser.
