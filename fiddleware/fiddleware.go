package fiddleware

import (
	"context"
	"github.com/bernsblack/fiddleware/util"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	//RequestIdContextKey is a UUID that is used to identify the request
	RequestIdContextKey   = "requestIdContextKey"
	HandlerNameContextKey = "handlerName"
	HandlerPathContextKey = "handlerPath"
	LiteralPathContextKey = "literalPath" // literal path is the same as the handler path
)

type LoggingResponseWriter struct {
	ctx         context.Context
	logger      util.Logger
	recorder    Recorder
	innerWriter http.ResponseWriter
}

func (f *Fiddleware) wrapResponseWriter(ctx context.Context, w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{ctx, f.logger, f.recorder, w}
}

func (s *LoggingResponseWriter) Header() http.Header {
	return s.innerWriter.Header()
}

func (s *LoggingResponseWriter) WriteWithError(bytes []byte, errors ...error) (int, error) {
	handlerName := s.ctx.Value(HandlerNameContextKey).(string)
	handlerPath := s.ctx.Value(HandlerPathContextKey).(string)
	literalPath := s.ctx.Value(LiteralPathContextKey).(string)
	requestId := s.ctx.Value(RequestIdContextKey).(string)

	s.recorder.Record(requestId, bytes) // records the response dump because hte request dump has already been set - consider having explicit request dump recorder methods

	if len(errors) > 0 {
		for _, err := range errors {
			s.logger.Errorf("[%s] [%s] error for '%s' => %+v", handlerName, handlerPath, literalPath, err)
		}
	}
	return s.Write(bytes)
}

// Write is a wrapper for the response data that gets written back to the client. The data is intercepted
func (s *LoggingResponseWriter) Write(bytes []byte) (int, error) {
	handlerName := s.ctx.Value(HandlerNameContextKey).(string)
	handlerPath := s.ctx.Value(HandlerPathContextKey).(string)
	literalPath := s.ctx.Value(LiteralPathContextKey).(string)
	s.logger.Infof(
		`
============================================RESPONSE=START==============================================================
[%s] [%s] response for '%s' => %s
============================================RESPONSE=STOP===============================================================
`, handlerPath, literalPath, handlerName, string(bytes))
	return s.innerWriter.Write(bytes)
}

func (s *LoggingResponseWriter) WriteHeader(statusCode int) {
	s.innerWriter.WriteHeader(statusCode)
}

type Fiddleware struct {
	logger      util.Logger
	recorder    Recorder
	innerRouter *mux.Router

	// TODO: REFACTOR TO USE map[string]Settings instead of map[string]struct{} or rather a sync.Map
	logSettings map[string]struct{} // map of paths to log - modify this to select log settings per path - for example logging the requests, time, responses
}

// Wrap will wrap the router to allow for logging and recording of requests
func Wrap(router *mux.Router, logger util.Logger, recorder Recorder, logSettings map[string]struct{}) *Fiddleware {
	return &Fiddleware{
		innerRouter: router,
		logger:      logger,
		recorder:    recorder,
		logSettings: logSettings,
	}
}

// AddLoggingToHandler is a wrapper function that adds logging to handler func. One can turn on certain parts of the logging, e.g. response, request timeouts
func (f *Fiddleware) AddLoggingToHandler(ctx context.Context, handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handlerName := util.GetFunctionName(handler)
		handlerPath := ctx.Value(HandlerPathContextKey).(string)
		requestId := ctx.Value(RequestIdContextKey).(string)
		literalPath := request.URL.Path

		f.logger.Infof(`
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
------------------------------------------------------------------------------------------------------------------------
routerPath:[%s] actualPath:[%s] handlerName:[%s]
------------------------------------------------------------------------------------------------------------------------
`, handlerPath, request.URL.Path, handlerName)
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			f.logger.Infof(`
------------------------------------------------------------------------------------------------------------------------
routerPath:[%s] actualPath:[%s] handlerName:[%s] took %s
------------------------------------------------------------------------------------------------------------------------
<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
`, handlerPath, request.URL.Path, handlerName, duration)
		}()

		// do the logging here
		requestDumpBytes, err := httputil.DumpRequest(request, true)
		if err != nil {
			f.logger.Infof("Error dumping response => %+v", err)
		}

		f.recorder.Record(requestId, requestDumpBytes)

		f.logger.Infof(
			`
============================================REQUEST=START===============================================================
[%s] [%s] request for '%s' => %s
============================================REQUEST=STOP================================================================
`, handlerPath, literalPath, handlerName, requestDumpBytes)

		ctx = context.WithValue(ctx, HandlerNameContextKey, handlerName)
		ctx = context.WithValue(ctx, LiteralPathContextKey, request.URL.Path)
		wrappedWriter := f.wrapResponseWriter(ctx, writer)
		handler(wrappedWriter, request)
	}
}
func (f *Fiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.innerRouter.ServeHTTP(w, r)
}

// HandleFunc is a wrapper function that adds logging to handler func. One can turn on certain parts of the logging, e.g. response, request timeouts
func (f *Fiddleware) HandleFunc(pattern string, handler http.HandlerFunc) *mux.Route {
	// if the pattern is in the map, then log wrap it in the log middleware
	// do the wrapping on the registering of the handlers so that when serving we don't have to check if the pattern is in the map everytime.
	if _, ok := f.logSettings[pattern]; ok {
		ctx := context.WithValue(context.Background(), HandlerPathContextKey, pattern)
		ctx = context.WithValue(ctx, RequestIdContextKey, uuid.New().String())
		wrappedHandler := f.AddLoggingToHandler(ctx, handler)
		return f.innerRouter.HandleFunc(pattern, wrappedHandler)
	} else {
		return f.innerRouter.HandleFunc(pattern, handler)
	}
}
