## testkube scripts execution

Gets script execution details

### Synopsis

Gets script execution details, you can change output format

```
testkube scripts execution [flags]
```

### Options

```
  -h, --help   help for execution
```

### Options inherited from parent commands

```
  -c, --client string        Client used for connecting to testkube API one of proxy|direct (default "proxy")
      --go-template string   in case of choosing output==go pass golang template (default "{{ . | printf \"%+v\"  }}")
  -s, --namespace string     kubernetes namespace (default "default")
  -o, --output string        output typoe one of raw|json|go  (default "raw")
  -v, --verbose              should I show additional debug messages
```

### SEE ALSO

* [testkube scripts](testkube_scripts.md)	 - Scripts management commands
