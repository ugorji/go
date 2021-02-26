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

codecgen caveats:

- Canonical option.  
  If Canonical=true, codecgen'ed code will delegate encoding maps to reflection-based code.  
  This is due to the runtime work needed to marshal a map in canonical mode.
- CheckCircularRef option.  
  When encoding a struct, a circular reference can lead to a stack overflow.  
  If CheckCircularRef=true, codecgen'ed code will delegate encoding structs to reflection-based code.
- MissingFielder implementation.  
  If a type implements MissingFielder, a Selfer is not generated (with a warning message).
  Statically reproducing the runtime work needed to extract the missing fields and marshal them along with the struct fields,
  while handling the Canonical=true special case, was onerous to implement.

**More Information**

Please see the [blog article](http://ugorji.net/blog/go-codecgen)
for more information on how to use the tool.
