package build_coordinator

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

var (
	SERVICE_CREATION_JSON = map[string]interface{}{
		"data": []map[string]interface{}{
			map[string]interface{}{
				"type": "service",
				"attributes": map[string]interface{}{
					"description": "test string",
					"methods":     []string{http.MethodGet, http.MethodPost},
					"parameters": []map[string]interface{}{
						map[string]interface{}{
							"type":        "string",
							"description": "test string",
							"constraints": []map[string]interface{}{
								map[string]interface{}{
									"key":   "maxLength",
									"value": 10,
								},
								map[string]interface{}{
									"key":   "minLength",
									"value": 5,
								},
							},
							"required": false,
						},
						map[string]interface{}{
							"type":        "integer",
							"description": "test integer",
							"constraints": []map[string]interface{}{
								map[string]interface{}{
									"key":   "maximum",
									"value": 1000,
								},
								map[string]interface{}{
									"key":   "minimum",
									"value": 500,
								},
							},
							"required": true,
						},
						map[string]interface{}{
							"type":        "array",
							"description": "test array",
							"constraints": []map[string]interface{}{
								map[string]interface{}{
									"key":   "maxItems",
									"value": 1000,
								},
								map[string]interface{}{
									"key":   "uniqueItems",
									"value": true,
								},
							},
							"required": true,
						},
					},
					"registers":        []interface{}{},
					"subscribers":      []interface{}{},
					"mode":             "LBModeRoundRobin",
					"dependencies":     []interface{}{},
					"version":          "1.0",
					"platform_version": "golang1.8.1",
				},
			},
		},
	}

	SERVICE_REGISTRY_JSON = map[string]interface{}{
		"data": []map[string]interface{}{
			map[string]interface{}{
				"type": "registry",
				"attributes": map[string]interface{}{
					"ip":   "localhost",
					"port": 12345,
					"api":  "/test",
				},
			},
		},
	}

	SERVICE_SUBSCRIPTION_JSON = map[string]interface{}{
		"data": []map[string]interface{}{
			map[string]interface{}{
				"type":       "subscription",
				"attributes": map[string]interface{}{},
			},
		},
	}
)

func TestServicesCRUDPositiveScenerios(t *testing.T) {
	var (
		// create httpexpect instance
		e = httpexpect.New(t, "http://localhost:8080")
	)

	// service creation
	response := e.POST("/v1/services").
		WithJSON(SERVICE_CREATION_JSON).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Object().
		Value("data").
		Array().
		Element(0).
		Object()

	// services query
	e.GET("/v1/services/{srvid}").
		WithPath(
			"srvid",
			response.
				Value("id").
				Raw(),
		).
		Expect().
		Status(http.StatusOK).
		JSON()

	// services query
	e.GET("/v1/services").
		Expect().
		Status(http.StatusOK).
		JSON()

	// service alter
	e.PUT("/v1/services/{srvid}").
		WithPath(
			"srvid",
			response.
				Value("id").
				Raw(),
		).
		WithJSON(
			GetServiceAlterJSON(
				response.
					Value("id").
					Raw(),
			),
		).
		Expect().
		Status(http.StatusOK).
		JSON()

	// service reister
	e.POST("/v1/services/{srvid}/registry").
		WithPath(
			"srvid",
			response.
				Value("id").
				Raw(),
		).
		WithJSON(SERVICE_REGISTRY_JSON).
		Expect().
		Status(http.StatusOK).
		JSON()

	// service subscribe
	response_subscription := e.POST("/v1/services/{srvid}/subscription").
		WithPath(
			"srvid",
			response.
				Value("id").
				Raw(),
		).
		WithJSON(SERVICE_SUBSCRIPTION_JSON).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("data").
		Array().
		Element(0).
		Object()

	// service deregister
	e.DELETE("/v1/services/{srvid}/registry/{ip}/{port}").
		WithPath("srvid", response.
			Value("id").
			Raw(),
		).
		WithPath(
			"ip",
			"localhost",
		).
		WithPath(
			"port",
			12345,
		).
		Expect().
		Status(http.StatusOK).
		JSON()

	e.DELETE("/v1/services/{srvid}/subscription/{token}").
		WithPath("srvid", response.
			Value("id").
			Raw(),
		).
		WithPath(
			"token",
			response_subscription.
				Value("id").
				Raw(),
		).
		Expect().
		JSON()

	// service delete
	e.DELETE("/v1/services/{srvid}").
		WithPath(
			"srvid",
			response.
				Value("id").
				Raw(),
		).
		Expect().
		Status(http.StatusOK).
		JSON()
}

func TestServicesCRUDNegativeScenerios(t *testing.T) {
	var (
		// create httpexpect instance
		e = httpexpect.New(t, "http://localhost:8080")
	)

	// service creation
	response := e.POST("/v1/services").
		WithJSON(SERVICE_CREATION_JSON).
		Expect().
		Status(http.StatusCreated).
		JSON().
		Object().
		Value("data").
		Array().
		Element(0).
		Object()

	// services query
	e.GET("/v1/services/{srvid}").
		WithPath(
			"srvid",
			"error_string",
		).
		Expect().
		Status(http.StatusNotFound).
		JSON()

	// service alter
	e.PUT("/v1/services/{srvid}").
		WithPath(
			"srvid",
			"error_string",
		).
		WithJSON(
			GetServiceAlterJSON(
				"error_string",
			),
		).
		Expect().
		Status(http.StatusNotFound).
		JSON()

	// service reister
	e.POST("/v1/services/{srvid}/registry").
		WithPath(
			"srvid",
			response.
				Value("id").
				Raw(),
		).
		WithJSON(SERVICE_REGISTRY_JSON).
		Expect().
		Status(http.StatusOK).
		JSON()

	// service subscribe
	e.POST("/v1/services/{srvid}/subscription").
		WithPath(
			"srvid",
			response.
				Value("id").
				Raw(),
		).
		WithJSON(SERVICE_SUBSCRIPTION_JSON).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Value("data").
		Array().
		Element(0).
		Object()

	// service deregister
	e.DELETE("/v1/services/{srvid}/registry/{ip}/{port}").
		WithPath("srvid", response.
			Value("id").
			Raw(),
		).
		WithPath(
			"ip",
			"error_string",
		).
		WithPath(
			"port",
			12345,
		).
		Expect().
		Status(http.StatusConflict).
		JSON()

	e.DELETE("/v1/services/{srvid}/subscription/{token}").
		WithPath("srvid", response.
			Value("id").
			Raw(),
		).
		WithPath(
			"token",
			"error_string",
		).
		Expect().
		Status(http.StatusConflict).
		JSON()

	// service delete
	e.DELETE("/v1/services/{srvid}").
		WithPath(
			"srvid",
			"error_string",
		).
		Expect().
		Status(http.StatusNotFound).
		JSON()
}

func GetServiceAlterJSON(id interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"type": "service",
			"attributes": map[string]interface{}{
				"service_id":       id,
				"description":      "test string",
				"methods":          []string{http.MethodGet},
				"parameters":       []map[string]interface{}{},
				"registers":        []interface{}{},
				"subscribers":      []interface{}{},
				"mode":             "LBModeRoundRobin",
				"dependencies":     []interface{}{},
				"version":          "1.0",
				"platform_version": "golang1.8.1",
			},
		},
	}
}
