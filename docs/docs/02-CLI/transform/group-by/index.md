# group-by

group string array by prefix to json

```
awesome-ci transform group-by [flags]
```

## Examples

```
awesome-ci tf group-by -p 3 --sub-prefix 4 st1_infrastructure-base  st2_ubi9-openjdk-11  st2_ubi9-openjdk-17 
   produces: {"st1":["infrastructure-base"],"st2":["ubi9-openjdk-11","ubi9-openjdk-17"]}
```

## Options

```
  -h, --help             help for group-by
  -p, --prefix int       group by prefix until index -- eg.: 'st1_base-image' and int 3, will be grouped by 'st1' (default 3)
      --sub-prefix int   remove prefix until index number -- eg.: 'st1_base-image' and 4 will be 'base-image'
```

## Options inherited from parent commands

```
  -v, --verbose   verbose output
```

## SEE ALSO

* **awesome-ci transform**	 - transform given input to json

##### Auto generated on 18-Mar-2025
