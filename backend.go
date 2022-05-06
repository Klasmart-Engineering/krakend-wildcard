package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/luraproject/lura/v2/logging"
)

const pluginName = "wildcard"
const headerName = "X-Krakend-Wildcard"
const logPrefix = "[PLUGIN:WILDCARD]"

type registerer string

var logger logging.Logger = nil

func (r registerer) RegisterLogger(v interface{}) {
	l, ok := v.(logging.Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(logPrefix, HandlerRegisterer, "plugin loaded!!!")
}

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer(pluginName)

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, config map[string]interface{}, handler http.Handler) (http.Handler, error) {

	if !configContainsPlugin(config) {
		return nil, fmt.Errorf("%s plugin was not named in configuration", pluginName)
	}

	endpoints := config[pluginName].(map[string]interface{})["endpoints"].(map[string]interface{})

	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no target endpoints configured")
	}
	targetEndpoints := parseEndpoints(endpoints)

	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		splitPath := strings.Split(req.URL.Path, "/")
		if len(splitPath) == 0 {
			handler.ServeHTTP(w, req)
			return
		}
		pathToCheck := splitPath[1]
		if pathToCheck == "__wildcard" {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		targetUrl, ok := targetEndpoints[pathToCheck]

		if !ok {
			handler.ServeHTTP(w, req)
			return
		}
		var targetPath string
		if len(splitPath) < 2 {
			targetPath = "/"
		} else {
			var builder strings.Builder
			builder.WriteString("/")
			builder.WriteString(strings.Join(splitPath[2:], "/"))
			targetPath = builder.String()
		}

		req.Header.Set(headerName, targetPath)
		req.URL.Path = targetUrl
		logger.Debug(logPrefix, "routing traffic to", req.URL.Path)
		handler.ServeHTTP(w, req)
	}), nil
}

func init() {}

func parseEndpoints(endpoints map[string]interface{}) map[string]string {
	targetEndpoints := make(map[string]string)

	for wc, endpoint := range endpoints {
		ep := endpoint.([]interface{})
		for _, e := range ep {
			if strings.HasPrefix(e.(string), "/") {
				targetEndpoints[e.(string)[1:]] = wc
			} else {
				targetEndpoints[e.(string)] = wc
			}

		}
	}

	return targetEndpoints
}

func configContainsPlugin(extra map[string]interface{}) bool {
	s := reflect.ValueOf(extra["name"])
	if s.Kind() == reflect.Slice {
		xs := extra["name"].([]interface{})
		for _, n := range xs {
			if n == pluginName {
				return true
			}
		}
	} else if s.Kind() == reflect.String && extra["name"].(string) == pluginName {
		return true
	}
	return false
}
