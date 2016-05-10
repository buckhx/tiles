package tiles

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Server is a simple server that indexes data into tiles
// GET / - lists indexes
// POST /index "index={string}&lat={float}&lon={float}&value={string}"- Adds a value to the index at lat and lon
// GET /index?index={string}&x={int}&y={int}&z={int} - Aggregates data and returns an array for all data under tile at x/y/z
// GET /tile?lat={float}&lon={float}&z={int} - Returns the tile x/y/z as json, also accepts POST data w/ those params
type Server struct {
	port    string
	indexes *indexer
	mux     *http.ServeMux
}

// NewServer returns a *Server whose router has been bound to a tile indexer
func NewServer(port string) *Server {
	indexes := newIndexer()
	mux := http.NewServeMux()
	mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		var (
			err  error
			v    interface{}
			name string
		)
		if err = r.ParseForm(); err == nil {
			if name = r.Form.Get("name"); name == "" {
				err = fmt.Errorf("'name' required")
			}

		}
		switch {
		case err != nil:
			// some kind of eror, skip methods
		case r.Method == http.MethodGet:
			x, ex := strconv.Atoi(r.Form.Get("x"))
			y, ey := strconv.Atoi(r.Form.Get("y"))
			z, ez := strconv.Atoi(r.Form.Get("z"))
			switch {
			case ex != nil:
				err = fmt.Errorf("'x' required as integer: %s", r.Form.Get("x"))
			case ey != nil:
				err = fmt.Errorf("'y' required as integer: %s", r.Form.Get("y"))
			case ez != nil:
				err = fmt.Errorf("'z' required as integer: %s", r.Form.Get("z"))
			//TODO more validation
			default:
				t := Tile{X: x, Y: y, Z: z}
				v = indexes.get(name).Values(t)
			}
		case r.Method == http.MethodPost:
			lat, ey := strconv.ParseFloat(r.Form.Get("lat"), 64)
			lon, ex := strconv.ParseFloat(r.Form.Get("lon"), 64)
			val := r.Form.Get("value")
			switch {
			case ex != nil:
				err = fmt.Errorf("'lat' required as float: %s", r.Form.Get("lat"))
			case ey != nil:
				err = fmt.Errorf("'lon' required as float: %s", r.Form.Get("lon"))
			case val == "":
				err = fmt.Errorf("'value' required as string: %s", r.Form.Get("value"))
			default:
				t := FromCoordinate(lat, lon, ZMax)
				indexes.get(name).Add(t, val)
				v = fmt.Sprintf("Indexed %s into Tile %+v", val, t)
			}
		}
		writeJSON(w, v, err)
	})
	mux.HandleFunc("/tile", func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
			v   interface{}
		)
		err = r.ParseForm()
		lat, ey := strconv.ParseFloat(r.Form.Get("lat"), 64)
		lon, ex := strconv.ParseFloat(r.Form.Get("lon"), 64)
		z, ez := strconv.Atoi(r.Form.Get("z"))
		switch {
		case err != nil:
			// noop
		case ex != nil:
			err = fmt.Errorf("'lon' required as float: %s", r.Form.Get("lon"))
		case ey != nil:
			err = fmt.Errorf("'lat' required as float: %s", r.Form.Get("lat"))
		case ez != nil:
			err = fmt.Errorf("'z' required as integer: %s", r.Form.Get("z"))
		default:
			t := FromCoordinate(lat, lon, z)
			v = map[string]interface{}{"x": t.X, "y": t.Y, "z": t.Z}
		}
		writeJSON(w, v, err)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		v := indexes.keys()
		writeJSON(w, v, nil)
	})
	return &Server{
		port:    port,
		indexes: indexes,
		mux:     mux,
	}
}

// Mux returns the router bound to the tileindexes
func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

// ListenAndServe starts the server
func (s *Server) ListenAndServe() error {
	log.Print("buckhx/tiles listening at ", s.port)
	defer log.Print("done listening")
	return http.ListenAndServe(s.port, s.mux)
}

// ListenAndServeTLS starts the server with TLS transport
func (s *Server) ListenAndServeTLS(cert, key string) error {
	log.Print("buckhx/tiles listening over TLS at ", s.port)
	defer log.Print("done listening")
	return http.ListenAndServeTLS(s.port, cert, key, s.mux)
}

// indexer is a thread safe tile index container
type indexer struct {
	indexes map[string]TileIndex
	sync.RWMutex
}

func newIndexer() *indexer {
	return &indexer{indexes: make(map[string]TileIndex)}
}

// creates a new index if key doesn't exist. Key should be url-safe
func (i *indexer) get(key string) TileIndex {
	i.RLock()
	idx := i.indexes[key]
	i.RUnlock()
	if idx == nil {
		idx = NewTileIndex()
		i.Lock()
		i.indexes[key] = idx
		i.Unlock()
	}
	return idx
}

func (i *indexer) keys() []string {
	i.RLock()
	defer i.RUnlock()
	keys := make([]string, len(i.indexes))
	c := 0
	for k := range i.indexes {
		keys[c] = k
		c++
	}
	return keys
}
