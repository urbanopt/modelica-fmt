# modelica-fmt

The Modelica Formatter provides the ability to automatically format Modelica code providing a consistent file structure to aid in the readability of the files and the ability to compare the difference between similar files.

## Running

```bash
modelica-fmt [-w] [-help] <sources>...
Options:
  -w  overwrite source with formatted output. If flag is not present print to stdout
Arguments:
  sources  one or more files or directories to format
```

To run the examples:

```bash
./modelica-fmt examples/gmt-building.mo > examples/gmt-building-out.mo
./modelica-fmt examples/gmt-coolingtower.mo > examples/gmt-coolingtower-out.mo
```

The resulting .mo file can be diffed to the previous file to compare how the modelica-fmt updates the file.

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
brew install go

# in the repository root directory
go build -o modelicafmt
```


## Updating Parser (Modelica Grammar)

If the grammar file (Modelica.g4) has been edited, you'll need to regenerate the parser by running the following commmand which runs Antlr in a Docker container.
```bash
./generate_parser.sh
```

## Known Issues






