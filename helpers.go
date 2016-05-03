package tiles

// CoordinateToTile take float lat/lons and a zoom and return a tile
func CoordinateToTile(lat, lon float64, zoom uint) (Tile, TilePixel) {
	coord := ClippedCoords(lat, lon)
	pixel := coord.ToPixel(zoom)
	return pixel.ToTile()
}
