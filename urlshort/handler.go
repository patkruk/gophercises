package urlshort

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
	URL  string `yaml:"url,omitempty" json:"url,omitempty"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse the yaml
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yamlBytes, &pathUrls)
	if err != nil {
		return nil, err
	}

	// convert yaml into map
	pathToUrls := convertSliceOfPathUrlsToMap(pathUrls)

	// return MapHandler using the map
	return MapHandler(pathToUrls, fallback), nil
}

// JSONHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse the json
	var pathUrls []pathUrl
	err := json.Unmarshal(jsonBytes, &pathUrls)
	if err != nil {
		return nil, err
	}

	// convert json into map
	pathToUrls := convertSliceOfPathUrlsToMap(pathUrls)

	// return MapHandler using the map
	return MapHandler(pathToUrls, fallback), nil
}

// convertSliceOfPathUrlsToMap converts slices of pathUrl into a map
func convertSliceOfPathUrlsToMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)

	for _, pu := range pathUrls {
		pathToUrls[pu.Path] = pu.URL
	}

	return pathToUrls
}
