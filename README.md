#tiles [![](https://godoc.org/github.com/buckhx/tiles?status.svg)](https://godoc.org/github.com/buckhx/tiles) [![Build Status](https://travis-ci.org/buckhx/tiles.svg?branch=master)](https://travis-ci.org/buckhx/tiles)
A map tiling library written in pure Go with no external dependencies outside of the golang stdlib. Converts between WGS84 coordinates to Slippy Map Tiles and Quadkeys. Includes a TileIndex that can be used to aggregate data into parent tiles.

    go get github.com/buckhx/tiles

The godoc has good documentation and a couple examples, but here are a few as well

##### Coordinate/Tile conversion
```
// Empire State Build coordinates and zoom level 18
z := 18
lat := 40.7484
lon := -73.9857

// Basic usage
t1 := tiles.FromCoordinate(lat, lon, z)

// If you want a little more granularity
x := tiles.Coordinate{Lat: lat, Lon: lon}
p := x.ToPixel(z)
t2, _ := p.ToTile()
// t1 == t2
```

##### TileIndex
The TileIndex allows for data to be indexed by tiles and aggregated up to their parents when requested
```
idx := tiles.NewTileIndex()
esb := tiles.FromCoordinate(40.7484, -73.9857, 18)
sol := tiles.FromCoordinate(40.6892, -74.0445, 18)
bbn := tiles.FromCoordinate(51.5007, -0.1246, 18)
idx.Add(esb, "EmpireStateBuilding")
idx.Add(sol, "StatueOfLiberty")
idx.Add(bbn, "BigBen")
nyc := tiles.Tile{X: 75, Y: 96, Z: 8}
den := tiles.Tile{X: 106, Y: 194, Z: 9}
fmt.Println("ESB Tile: ", idx.Values(esb))
fmt.Println("SOL Tile: ", idx.Values(sol))
fmt.Println("NYC Tile: ", idx.Values(nyc))    //contains esb & sol values!
fmt.Println("DENVER Tile: ", idx.Values(den)) //contains no values!
```

##### Benchmarks

Here are some microbenchmarks for converting a location at zoom level 18. 
There's nothing really to compare them to, but should give a sense of op time on a 2.3 GHz core i7 MBP. 
```
$ go test -bench=. -benchmem -benchtime 10s
BenchmarkTileFromCoordinate-8	100000000	       165 ns/op	       0 B/op	       0 allocs/op
BenchmarkTileFromQuadkey-8   	300000000	        43.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkQuadkeyFromTile-8   	200000000	        95.4 ns/op	      32 B/op	       1 allocs/op
```
