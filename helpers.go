package tiles

func CoordinateToTile(lat, lon float64, zoom uint) (Tile, TilePixel) {
	coord := ClippedCoords(lat, lon)
	pixel := coord.ToPixel(zoom)
	return pixel.ToTile()
}
