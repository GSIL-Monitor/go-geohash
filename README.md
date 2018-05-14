# go-geohash

This is golang version geohash with bounding box searches refers to https://github.com/davetroy/geohash-js


## BenchMark Result

Test with MacOS High Sierra, 2.3 GHz Intel Core i5, 8 GB 2133 MHz LPDDR3

````
BenchmarkEncodeGeoHashPrecision12-4   	 3000000	       577 ns/op
BenchmarkEncodeGeoHashPrecision6-4    	 5000000	       347 ns/op
BenchmarkDecodeGeoHash-4              	 5000000	       330 ns/op
BenchmarkGetAdjacentGridGeoHash1-4    	 2000000	       919 ns/op
BenchmarkGetAdjacentGridGeoHash2-4    	 2000000	       939 ns/op
````
