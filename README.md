# modelica-fmt

The Modelica Formatter provides the ability to automatically format Modelica code providing a consistent file structure to aid in the readability of the files and the ability to compare the difference between similar files.

## Running

```bash
modelica-fmt [-w] [-help] <sources>...
Options:
  -w            overwrite source with formatted output. If flag is not present print to stdout
  -template     template dialect used for .mot/.mopt files (currently: jinja)
  -wrap-arrays  wrap multidimensional arrays ({...}) across multiple lines (outside annotations)
Arguments:
  sources  one or more files or directories to format
```

### Array formatting (`-wrap-arrays`)

By default arrays (`{...}`) are kept on a single line. With `-wrap-arrays`,
multidimensional arrays outside of annotations are broken across lines in a
compact style: the innermost two dimensions are kept inline while the outer
`max(N-2, 1)` dimension levels are broken. For example:

```modelica
parameter Integer array2D[2,2]={
  {1,2},
  {3,4}};
parameter Integer array3D[2,2,2]={
  {{1,2},{3,4}},
  {{5,6},{7,8}}};
parameter Integer array4D[2,2,2,2]={
  {
    {{1,2},{3,4}},
    {{5,6},{7,8}}},
  {
    {{9,10},{11,12}},
    {{13,14},{15,16}}}};
```

1D arrays, iterator constructors, matrices (`[...]`) and arrays inside
annotations are left unchanged.

To try the formatter against the bundled test data:

```bash
./modelica-fmt internal/format/testdata/gmt-building.mo > internal/format/testdata/gmt-building-out.mo
./modelica-fmt internal/format/testdata/gmt-coolingtower.mo > internal/format/testdata/gmt-coolingtower-out.mo
```

The resulting .mo file can be diffed to the previous file to compare how the modelica-fmt updates the file.

## Templated Modelica (`.mot`/`.mopt`) files

`modelica-fmt` can also format Modelica **template** files (`.mot` and `.mopt`) — as used by
[geojson-modelica-translator (GMT)](https://github.com/urbanopt/geojson-modelica-translator) —
which embed [Jinja](https://jinja.palletsprojects.com/) constructs (`{{ ... }}`, `{% ... %}`)
that aren't valid Modelica on their own. Files ending in `.mot` or `.mopt` are detected automatically,
including during directory walks, so no extra flag is required:

```bash
./modelica-fmt -w path/to/Template.mot
./modelica-fmt -w path/to/Template.mopt
./modelica-fmt -w path/to/templates/   # formats .mo, .mot, and .mopt files found in the tree
```

Internally this uses a substitute → format → reverse round trip: template constructs are
temporarily replaced with placeholders that the Modelica lexer accepts (control statements
`{% ... %}` are commented out, inline expressions such as `{{ ... }}` and `${...}` become
bare identifiers, and GMT generated-snippet expressions such as `{{ model.instance }}` are
commented out), the file is formatted, and then the original template constructs are
restored. Formatting is idempotent and only affects whitespace/layout — template constructs
are preserved exactly.

The template dialect is selectable with `-template` (currently only `jinja` is supported):

```bash
./modelica-fmt -w -template jinja path/to/Template.mot
```

Some template files are not Modelica classes at all, such as Dymola run scripts. When a
templated `.mot` or `.mopt` file still cannot be made parseable through preprocessing, `modelica-fmt`
emits the original content unchanged rather than failing the whole directory run or writing
empty output.

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
    files: \.(mo|mot|mopt)$
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




