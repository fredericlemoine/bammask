bammask
========
[![build](https://github.com/fredericlemoine/bammask/actions/workflows/go.yml/badge.svg)](https://github.com/fredericlemoine/bammask/actions)

This tool just reads a bam file and replace bases that are below a given quality threshold by N.

Usage:
-----
```
Usage:
  bammask quality [flags]

Flags:
  -h, --help               help for quality
  -i, --input-bam string   Input bam file (default "stdin")
  -o, --out-bam string     Output bam file (default "stdout")
  -q, --quality int        Quality cutoff below which bases are masked (default 20)
```

Install:
-------

```
go get .
go build .
```

New location:
------------
The functions are now integrated in https://github.com/fredericlemoine/fastqutils .
