package tiles

// CoordinateToTile take float lat/lons and a zoom and return a tile
func CoordinateToTile(lat, lon float64, zoom int) (Tile, TilePixel) {
	coord := ClippedCoords(lat, lon)
	pixel := coord.ToPixel(zoom)
	return pixel.ToTile()
}
