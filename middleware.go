package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"time"
	"sync/atomic"
	"os"
)

// Context is holder of arbitrary data which might appear inside requests.
type Context struct {
	rpm uint64 // Global request per minute counter
	cache *Cache
}

// flush is Context method to flush the current API state to file.
func (c *Context) flush() error {
	serializer := map[string]interface{}{
		"last_active_at": time.Now(), "last_started_at": server.started_at,
		"cache_state": c.cache.state, "rpm": c.rpm,
	}
	dump, err := json.Marshal(serializer)

	if err != nil {
		return err
	}
	return ioutil.WriteFile(APIStateFilename, dump, 0600)
}

// load is Context method used to source state from current file dump into memory.
func (c *Context) load() (int, error) {

	var stateMap map[string]interface{}
	var delay int = 0

	stateByte, ioErr := ioutil.ReadFile(APIStateFilename)
	parseErr := json.Unmarshal(stateByte, &stateMap)

	if ioErr != nil || parseErr != nil {
		return delay, fmt.Errorf(StateLoadError)
	}

	// Reset RPM if service was down more then 60 seconds.
	lastActive, parseErr := time.Parse(time.RFC3339, stateMap["last_active_at"].(string))
	if parseErr != nil {
		return delay, parseErr
	}

	// If service was down less then 60 seconds, calculate delay to sync up with last execution.
	lastStarted, parseErr := time.Parse(time.RFC3339, stateMap["last_started_at"].(string))
	if parseErr != nil {
		return delay, parseErr
	}

	downtime := int(time.Now().Unix() - lastActive.Unix())
	if downtime >= 60 {
		c.rpm = 0
		delay = 0
	} else {
		delay = int(time.Now().Unix() - lastStarted.Unix())
	}

	c.rpm = uint64(stateMap["rpm"].(float64))

	return delay, nil
}

// startBeat is go routine method which will start a context ticker which will reset rpm every minute.
// If last API restart happened less then 60 seconds, delay is amount of seconds go routine will subtract from reset interval
// to sync up with previous execution.
func (c *Context) startBeat(delay int) {
	resetInterval := 60

	if delay < 60 {
		resetInterval = resetInterval - delay
	}

	for {
		<-time.After(time.Duration(resetInterval) * time.Second)
		c.rpm = 0
	}
}

func (c *Context) getRpm() uint64 {
	return atomic.LoadUint64(&c.rpm)
}

// NewContext is constructor for Context object.
func NewContext(sourceOld bool) *Context {
	c := &Context{
		rpm: 0,
		cache: SingleCache(),
	}
	var delaySec int = 0
	if sourceOld {
		delay, err := c.load()
		if err != nil {
			os.Remove(APIStateFilename)
		} else {
			delaySec = delay
		}

	}
	go c.startBeat(delaySec)
	return c

}

// APIHandler implements http.Handler and wraps Context object around APIViews
type APIHandler struct {
	*Context
	Handler APIView
}


// trackStats is method to track execution times on per route basis.
func (handler APIHandler) trackStats(route string, execTime int64) {
	route_bucket, _ := handler.cache.Get(route)
	bucket := route_bucket.(map[string]interface{})
	var execTimes = []int64{}
	if bucket["exec_times_ns"] != nil {
		execTimes = bucket["exec_times_ns"].([]int64)
	}
	execTimes = append(execTimes, execTime)

	var total int64 = 0
	for _, val := range execTimes {
		total += val
	}
	response_time := total / int64(len(execTimes))

	// Remove oldest data point from the execution times if there are more than 5 data points.
	if len(execTimes) > 5 {
		execTimes = execTimes[1:]
	}

	bucket["exec_times_ns"] = execTimes
	bucket["response_time_ms"] = float64(float64(response_time) / float64(time.Millisecond))
	handler.cache.Set(route, bucket)
}


// AppHandler implements http.Handler.
func (handler APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// NOTE: Since our cache is thread-safe, we could track rpm with the cache like this:
	//	counter, _ := handler.cache.Get("rpm")
	//	handler.cache.Set("rpm", counter.(int)+1)
	//
	// NOTE: But rpm is unsigned integer and we have ability to use atomic operations on it.
	atomic.AddUint64(&handler.rpm, 1)

	// NOTE: We always want to return json as response.
	w.Header().Set("Content-Type", "application/json")

	// NOTE: We want to track avg execution time
	start := time.Now()
	if status, err := handler.Handler(*handler.Context, w, r); err != nil {
		switch status {
		case http.StatusInternalServerError:
			handler.Context.flush()  // In case of service error, make sure to save current state.
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			handler.Context.flush()  // In case of service error, make sure to save current state.
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else {
		elapsed := time.Since(start)
		handler.trackStats(r.URL.Path, int64(elapsed))
	}
}

// NewRoutes will return map of routes and their respective handlers.
func NewRoutes() map[string]APIView {
	return map[string]APIView {
		"/":        rootView,
		"/status":  statusView,
		"/compute": computeView,
	}
}

// APIRouter will construct and return http multiplexer / router.
func APIRouter(ctx *Context) *http.ServeMux {
	mux := http.NewServeMux()
	for route, handler := range NewRoutes() {

		// Set route bucket for tracking of route stats.
		ctx.cache.Set(fmt.Sprintf("%s", route), make(map[string]interface{}))

		handler := APIHandler{ctx, handler}
		mux.Handle(route, handler)
	}
	return mux
}
