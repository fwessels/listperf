## listperf

Parallel lister of objects in an S3 bucket, especially for objects with hexadecimal key names.

As an example use it like this (using one of the AWS Public Datasets, namely the [multimedia-commons](http://aws.amazon.com/public-data-sets/multimedia-commons/))

```
$ ./listperf -p "data/images/08" -b "multimedia-commons" -r "us-west-2" -e "https://s3-us-west-2.amazonaws.com"
Number of replies: 25743
```

or

```
$ ./listperf -p "data/images/80" -b "multimedia-commons" -r "us-west-2" -e "https://s3-us-west-2.amazonaws.com"
Number of replies: 388277
```
(note that 0x80 instead of 0x08 gives a bit more results...)

### Usage

```
Usage of listperf:
  -a string
    	access key
  -b string
    	bucket for search
  -cpus int
    	Number of CPUs to use. Defaults to number of processors. (default 8)
  -e string
    	endpoint for bucket
  -p string
    	prefix for search
  -r string
    	region for search
  -s string
    	secret key
  -seq
    	Run queries sequentially when true
```

### To parallel or not to parallel

If the `-seq true` flag is given the queries are done sequentially 

```
$ time ./listperf -p "data/images/08" -b "multimedia-commons" -r "us-west-2" -e "https://s3-us-west-2.amazonaws.com"
Number of replies: 25743

real	0m2.335s
user	0m1.243s
sys	0m0.120s

$ time ./listperf -p "data/images/08" -b "multimedia-commons" -r "us-west-2" -e "https://s3-us-west-2.amazonaws.com" -seq true
Number of replies: 25743

real	0m17.078s
user	0m1.194s
sys	0m0.132s
```

Not a bad speed up for minor code changes...