# codecgen tool

Generate is given a list of *.go files to parse, and an output file (fout),
codecgen will create an output file __file.go__ which
contains `codec.Selfer` implementations for the named types found
in the files parsed.

Using codecgen is very straightforward.

**Download and install the tool**

`go get -u github.com/ugorji/go/codec/codecgen`

**Run the tool on your files**

The command line format is:

`codecgen [options] (-o outfile) (infile ...)`

```sh
% codecgen -?
Usage of codecgen:
  -c="github.com/ugorji/go/codec": codec path
  -o="": out file
  -d="": random identifier for use in generated code (help reduce changes when files are regenerated)
  -nx=false: do not support extensions (use when you know no extensions are used)
  -r=".*": regex for type name to match
  -nr="": regex for type name to exclude
  -rt="": tags for go run
  -st="codec,json": struct tag keys to introspect
  -t="": build tag to put in file
  -u=false: Use unsafe, e.g. to avoid unnecessary allocation on []byte->string
  -x=false: keep temp file

% codecgen -o values_codecgen.go values.go values2.go moretypedefs.go
```

**Limitations**

codecgen doesn't support the following:

- Canonical option. 
  This has not been implemented.
- MissingFielder implementation.
  If a type implements MissingFielder, it is completely ignored by codecgen.
- CheckCircularReference.
  When encoding, if a circular reference is encountered, it will cause a stack overflow.
  Only use codecgen on types that you know do have circular references.

**More Information**

Please see the [blog article](http://ugorji.net/blog/go-codecgen)
for more information on how to use the tool.

