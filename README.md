# modelica-fmt

The Modelica Formatter provides the ability to automatically format Modelica code providing a consistent file structure to aid in the readability of the files and the ability to compare the difference between similar files.

## Running

```bash
modelica-fmt [options] <sources>...
```

### Options

| Flag | Default | Description |
| --- | --- | --- |
| `-w` | `false` | Overwrite the source file(s) with the formatted output. If omitted, the formatted result is printed to stdout. |
| `-v` | `false` | Display the tool version and exit. |
| `-line-length <n>` | `-1` | Maximum number of characters allowed per line before wrapping; `-1` means no limit. |
| `-extra-padding` | `false` | **BETA:** add empty lines for padding to improve visual separation. |
| `-format-html` | `false` | **BETA:** pretty-print HTML content embedded in annotation strings (e.g. `Documentation(info=...)` / `revisions=...`). See [Formatting HTML in annotations](#formatting-html-in-annotations-beta). |
| `-help` | | Print usage information and exit. |

### Arguments

| Argument | Description |
| --- | --- |
| `sources` | One or more `.mo` files or directories to format. Directories are searched for `.mo` files. |

For example, to overwrite a file in place while wrapping lines at 80 characters:

```bash
modelica-fmt -w -line-length 80 examples/gmt-building.mo
```

To run the examples:

```bash
./modelica-fmt examples/gmt-building.mo > examples/gmt-building-out.mo
./modelica-fmt examples/gmt-coolingtower.mo > examples/gmt-coolingtower-out.mo
```

The resulting .mo file can be diffed to the previous file to compare how the modelica-fmt updates the file.

## Formatting HTML in annotations (BETA)

Modelica annotations commonly embed HTML documentation, e.g.
`Documentation(info="<html>...</html>")` and `revisions="..."`. By default the entire
annotation string is emitted verbatim. Pass `--format-html` to pretty-print the embedded
HTML so it is indented consistently with the surrounding Modelica structure:

```bash
./modelica-fmt --format-html examples/gmt-building.mo
```

This feature is **opt-in** and intentionally conservative:

- A string is treated as HTML only when its content begins with `<html>` (the Modelica
  convention), so key names like `info`/`revisions` are not hard-coded.
- Only whitespace/indentation changes — tags, attributes, entities (e.g. `&amp;`) and text
  are preserved exactly, and the `\"` escaping inside Modelica strings is round-tripped
  losslessly. The result is idempotent.
- Block-level elements (e.g. `<p>`, `<ul>`, `<li>`, `<div>`), comments and doctypes are
  placed on their own lines indented to their nesting depth. Inline/phrasing elements
  (e.g. `<b>`, `<i>`, `<a>`, `<code>`, `<span>`, `<br/>`) and the text around them stay
  together on one line, so short markup such as `<b>Example</b>` or
  `<a href=\"...\">link</a>` is kept condensed instead of exploded one tag per line.
- Whitespace-sensitive elements (`<pre>`, `<textarea>`, `<script>`, `<style>`) are left
  verbatim.
- Malformed or unbalanced HTML is never reflowed (so it can't be corrupted). Instead the
  tool reports a clear error naming the offending tag, leaves the file unchanged, and exits
  non-zero (see [Behavior on malformed HTML](#behavior-on-malformed-html)).
- HTML lines do not participate in the `--line-length` limit.

### Behavior on malformed HTML

When `--format-html` is enabled and an annotation string looks like HTML (its content
starts with `<html>`) but is malformed or unbalanced — for example a mismatched or missing
closing tag — the tool does **not** silently emit it unformatted. It fails with a clear
message identifying the file and the problem, and leaves the file unchanged:

```console
$ modelica-fmt --format-html MyModel.mo
error: MyModel.mo: malformed HTML in annotation string: mismatched closing tag </html> (expected </p>)
skipping MyModel.mo (file left unchanged)
```

The process exits non-zero in this case, so the failure is visible in CI and pre-commit
hooks rather than being silently ignored. Fix the HTML (or run without `--format-html`) and
re-run. When `--format-html` is disabled, annotation strings are never inspected, so
malformed HTML has no effect.

### Known limitations

These are intentional v1 trade-offs, documented here so they are easy to revisit later:

- **Inline runs collapse to a single line.** A run of text and inline elements is placed on
  one line (its internal whitespace runs collapsed to single spaces), so a long paragraph
  becomes one long line; it is not wrapped to `--line-length`.
- **Conservative handling of unbalanced HTML.** HTML with omitted closing tags (e.g. a bare
  `<li>` or `<p>`) or mismatched tags is not reflowed; it is reported as an error and the
  file is left unchanged rather than reflowed, to avoid corrupting content.

## Usage with pre-commit framework

After adding modelicafmt to your system path, add the following lines to your .pre-commit-config.yaml file under the `repos:` section.
Also, make sure to allow modelicafmt to run (especially on Mac). 

```yaml
-
  repo: local
  hooks:
  -
    id: modelica-fmt
    name: Modelica Formatter
    types: [file]
    files: \.(mo)$
    entry: modelicafmt
    args: ["-w"]
    language: system
```
See https://pre-commit.com/ for more information about the framework.

## Building

```bash
# install go with brew, or follow instructions here: https://golang.org/doc/install
brew install go

# if you have an older version of the go tool, you may need to explicitly
# download the dependencies (in repo root directory)
go get -d ./...

# in the repository root directory
go build .

# optionally, with `-o` you can specify a custom name for your executable
# (on Windows remember to add `.exe` to the name)
go build -o modelicafmt
```


## Updating Parser (Modelica Grammar)

If the grammar file (Modelica.g4) has been edited, you'll need to regenerate the parser by running the following commmand which runs Antlr in a Docker container.
```bash
./generate_parser.sh
```

## Known Issues






