# prime
prime number tool for the command line

```
prime: generate a prime number and print to stdout
Example: 'prime -f 16 -b 128' prints: 83b19881300529d1fd4dac680415c60f
Example: 'prime -f 64 -b 256' prints: zUwiK96O4sy6pm3LtQM5YtRP9L4RGxsU/zCNliZZXn0=
Example: 'prime -f  0 -b 1024 > p.bytes' saves the raw bytes to the file 'p.bytes'
Options:
  -b int
    	number of bits [supports: 2,...,128,...] (default 128)
  -f int
    	format of output [supports: 0,2-36,64,85] (default 10)
```

# Examples

```
$prime -f 16 -b 128
83b19881300529d1fd4dac680415c60f
$prime -f 64 -b 256
zUwiK96O4sy6pm3LtQM5YtRP9L4RGxsU/zCNliZZXn0=
$prime -f  0 -b 1024 > p.bytes
saves the raw bytes to the file 'p.bytes'
```
