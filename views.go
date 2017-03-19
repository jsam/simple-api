package main

import (
	"net/http"
	"encoding/json"
	"math/big"
)

// APIView is abstraction for view components.
type APIView func(ctx Context, w http.ResponseWriter, r *http.Request) (int, error)

// rootView is handler for global root path /
func rootView(ctx Context, w http.ResponseWriter, r *http.Request) (int, error) {

	jsonResp, err := json.Marshal(map[string]interface{}{
		"rpm": ctx.getRpm(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResp)

	return http.StatusOK, nil
}

// statusView is view handler on path /status and it's used to see what is avarage response time on this endpoint.
func statusView(ctx Context, w http.ResponseWriter, r *http.Request) (int, error) {

	route_bucket, _ := ctx.cache.Get(r.URL.Path)
	bucket := route_bucket.(map[string]interface{})

	resp := map[string]interface{}{
		"stats": bucket,
		"rpm": ctx.getRpm(),
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonResp)
	return http.StatusOK, nil
}

// statusView is view handler on path /status and it's used to see what is avarage response time on this endpoint.
func computeView(ctx Context, w http.ResponseWriter, r *http.Request) (int, error) {

	// INFO: Do some compute heavy stuff to measure performance.
	num := new(big.Int)
	num.Binomial(100000, 1000)

	route_bucket, _ := ctx.cache.Get(r.URL.Path)
	bucket := route_bucket.(map[string]interface{})

	resp := map[string]interface{}{
		"stats": bucket,
		"rpm": ctx.getRpm(),
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonResp)
	return http.StatusOK, nil
}


