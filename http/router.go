// =========================================================================
// Copyright ©  2019 AppsByRam authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"net/http"

	"github.com/exceller-io/pkg/metrics"
	"github.com/gorilla/mux"
)

func newRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	//add health route
	router.
		Methods("GET").
		Path("/health").
		Name("health").
		Handler(healthHandler())

	//add metrics route
	router.
		Methods("GET").
		Path("/metrics").
		Name("metrics").
		Handler(metrics.PrometheusHandler())
	return router
}

//HealthReport represents a health report
type HealthReport struct {
	Status string `json:"status" yaml:"status"`
}

func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		report := HealthReport{
			Status: "UP",
		}
		p := NewPayload()
		p.WriteResponse(ContentTypeJSON, http.StatusOK, report, w)
	}
}
