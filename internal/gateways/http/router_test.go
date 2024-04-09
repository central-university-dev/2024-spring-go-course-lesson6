package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Все неизвестные пути должны возвращать http.StatusNotFound.
func TestUnknownRoute(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"GET", "GET", http.StatusNotFound},
		{"POST", "POST", http.StatusNotFound},
		{"PUT", "PUT", http.StatusNotFound},
		{"DELETE", "DELETE", http.StatusNotFound},
		{"HEAD", "HEAD", http.StatusNotFound},
		{"OPTIONS", "OPTIONS", http.StatusNotFound},
		{"PATCH", "PATCH", http.StatusNotFound},
		{"CONNECT", "CONNECT", http.StatusNotFound},
		{"TRACE", "TRACE", http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/unknown", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want, w.Code, "Получили в ответ не тот код")
		})
	}
}

// Тесты /users
func TestUsersRoutes(t *testing.T) {
	t.Run("POST_users", func(t *testing.T) {
		t.Run("valid_request_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"name": "Пользователь 1"
			}`
			req, _ := http.NewRequest("POST", "/users", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Получили в ответ не тот код")
			assert.True(t, json.Valid(w.Body.Bytes()), "В ответе не json")
		})

		t.Run("request_body_has_unsupported_format_415", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `<User>
				<Id>1</Id>
				<Name>Пользователь 1</Name>
			</User>`
			req, _ := http.NewRequest("POST", "/users", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/xml")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnsupportedMediaType, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_has_syntax_error_400", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{ невалидный json }`
			req, _ := http.NewRequest("POST", "/users", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_is_valid_but_it_has_invalid_data_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"id": -1,
				"name": "Пользователь -1"
			}`
			req, _ := http.NewRequest("POST", "/users", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("OPTIONS_users_204", func(t *testing.T) {
		router := gin.Default()
		setupRouter(router)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/users", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code, "Получили в ответ не тот код")
		allowed := strings.Split(w.Header().Get("Allow"), ",")
		assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
		assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
	})

	// Другие методы не поддерживаем.
	t.Run("OTHER_users_405", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  int
		}{
			{"GET", "GET", 405},
			{"PUT", "PUT", 405},
			{"DELETE", "DELETE", 405},
			{"HEAD", "HEAD", 405},
			{"PATCH", "PATCH", 405},
			{"CONNECT", "CONNECT", 405},
			{"TRACE", "TRACE", 405},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := gin.Default()
				setupRouter(router)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.input, "/users", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.want, w.Code, "Получили в ответ не тот код")
				allowed := strings.Split(w.Header().Get("Allow"), ",")
				assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
				assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
			})
		}
	})
}

// Тесты /sensors
func TestSensorsRoutes(t *testing.T) {
	t.Run("GET_sensors", func(t *testing.T) {
		t.Run("success_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Получили в ответ не тот код")
			assert.True(t, json.Valid(w.Body.Bytes()), "В ответе не json")
		})

		t.Run("requested_unsupported_body_format_406", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/sensors", nil)
			req.Header.Add("Accept", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 406, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("HEAD_sensors", func(t *testing.T) {
		t.Run("success_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Получили в ответ не тот код")
			assert.NotEmpty(t, w.Header().Get("Content-Length"), "Content-Length не задан")
		})

		t.Run("requested_unsupported_body_format_406", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/sensors", nil)
			req.Header.Add("Accept", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 406, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("POST_sensors", func(t *testing.T) {
		t.Run("valid_request_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"serial_number": 1234567890,
				"type": "cc",
				"description": "Датчик температуры",
				"is_active": true,
			}`
			req, _ := http.NewRequest("POST", "/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Получили в ответ не тот код")
			assert.True(t, json.Valid(w.Body.Bytes()), "В ответе не json")
		})

		t.Run("request_body_has_unsupported_format_415", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `<Sensor>
				<Id>1</Id>
				<SerialNumber>1234567890</SerialNumber>
				<Type>cc</Type>
				<CurrentState>1</CurrentState>
				<Description>Датчик температуры</Description>
				<IsActive>true</IsActive>
				<RegisteredAt>2018-01-01T00:00:00Z</RegisteredAt>
				<LastActivity>2018-01-01T00:00:00Z</LastActivity>
			</Sensor>`
			req, _ := http.NewRequest("POST", "/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 415, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_has_syntax_error_400", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{ невалидный json }`
			req, _ := http.NewRequest("POST", "/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 400, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_is_valid_but_it_has_invalid_data_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"id": -1,
				"serial_number": 1234567890,
				"type": "cc",
				"current_state": 1,
				"description": "Датчик температуры",
				"is_active": true,
				"registered_at": "2018-01-01T00:00:00Z",
				"last_activity": "2018-01-01T00:00:00Z"
			}`
			req, _ := http.NewRequest("POST", "/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 422, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("OPTIONS_sensors_204", func(t *testing.T) {
		router := gin.Default()
		setupRouter(router)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/sensors", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code, "Получили в ответ не тот код")
		allowed := strings.Split(w.Header().Get("Allow"), ",")
		assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
		assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
		assert.Contains(t, allowed, "GET", "В разрешённых методах нет GET")
		assert.Contains(t, allowed, "HEAD", "В разрешённых методах нет HEAD")
	})

	t.Run("GET_sensors_sensor_id", func(t *testing.T) {
		t.Run("sensor_exists_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/sensors/1", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Получили в ответ не тот код")
			assert.True(t, json.Valid(w.Body.Bytes()), "В ответе не json")
		})

		t.Run("id_has_invalid_format_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/sensors/1", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 422, w.Code, "Получили в ответ не тот код")
		})

		t.Run("requested_unsupported_body_format_406", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/sensors/1", nil)
			req.Header.Add("Accept", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 406, w.Code, "Получили в ответ не тот код")
		})

		t.Run("sensor_doesnt_exist_404", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/sensors/2", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("HEAD_sensors_sensor_id", func(t *testing.T) {
		t.Run("sensor_exists_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/sensors/1", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Получили в ответ не тот код")
			assert.NotEmpty(t, w.Header().Get("Content-Length"), "Content-Length не задан")
		})

		t.Run("id_has_invalid_format_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/sensors/1", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 422, w.Code, "Получили в ответ не тот код")
		})

		t.Run("requested_unsupported_body_format_406", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/sensors/1", nil)
			req.Header.Add("Accept", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 406, w.Code, "Получили в ответ не тот код")
		})

		t.Run("sensor_doesnt_exist_404", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/sensors/2", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("OPTIONS_sensors_sensor_id_204", func(t *testing.T) {
		router := gin.Default()
		setupRouter(router)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/sensors/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code, "Получили в ответ не тот код")
		allowed := strings.Split(w.Header().Get("Allow"), ",")
		assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
		assert.Contains(t, allowed, "GET", "В разрешённых методах нет GET")
		assert.Contains(t, allowed, "HEAD", "В разрешённых методах нет HEAD")
	})

	// Другие методы не поддерживаем.
	t.Run("OTHER_users_405", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  int
		}{
			{"GET", "GET", 405},
			{"PUT", "PUT", 405},
			{"DELETE", "DELETE", 405},
			{"HEAD", "HEAD", 405},
			{"PATCH", "PATCH", 405},
			{"CONNECT", "CONNECT", 405},
			{"TRACE", "TRACE", 405},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := gin.Default()
				setupRouter(router)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.input, "/users", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.want, w.Code, "Получили в ответ не тот код")
				allowed := strings.Split(w.Header().Get("Allow"), ",")
				assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
				assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
			})
		}
	})
}

// Тесты /users/{user_id}/sensors
func TestUsersSensorsRoutes(t *testing.T) {
	t.Run("GET_users_user_id_sensors", func(t *testing.T) {
		t.Run("user_exists_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/users/1/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Получили в ответ не тот код")
			assert.True(t, json.Valid(w.Body.Bytes()), "В ответе не json")
		})

		t.Run("id_has_invalid_format_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/users/abc/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 422, w.Code, "Получили в ответ не тот код")
		})

		t.Run("requested_unsupported_body_format_406", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/users/1/sensors", nil)
			req.Header.Add("Accept", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 406, w.Code, "Получили в ответ не тот код")
		})

		t.Run("user_doesnt_exist_404", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/users/2/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("HEAD_users_user_id_sensors", func(t *testing.T) {
		t.Run("user_exists_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/users/1/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Получили в ответ не тот код")
			assert.NotEmpty(t, w.Header().Get("Content-Length"), "Content-Length не задан")
		})

		t.Run("id_has_invalid_format_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/users/abc/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 422, w.Code, "Получили в ответ не тот код")
		})

		t.Run("requested_unsupported_body_format_406", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/users/1/sensors", nil)
			req.Header.Add("Accept", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 406, w.Code, "Получили в ответ не тот код")
		})

		t.Run("user_doesnt_exist_404", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("HEAD", "/users/2/sensors", nil)
			req.Header.Add("Accept", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("POST_users_user_id_sensors", func(t *testing.T) {
		t.Run("valid_request_body_and_user_exists_200", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"sensor_id": 1,
			}`
			req, _ := http.NewRequest("POST", "/users/1/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, 201, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_has_unsupported_format_415", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `<SensorToUserBinding>
				<SensorId>1</SensorId>
			</SensorToUserBinding>`
			req, _ := http.NewRequest("POST", "/users/1/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/xml")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, 415, w.Code, "Получили в ответ не тот код")
		})

		t.Run("invalid_request_body_400", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{ невалидный json }`
			req, _ := http.NewRequest("POST", "/users/1/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, 400, w.Code, "Получили в ответ не тот код")
		})

		t.Run("valid_request_body_but_user_doesnt_exist_404", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"sensor_id": 1,
			}`
			req, _ := http.NewRequest("POST", "/users/2/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_is_valid_but_it_has_invalid_data_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"sensor_id": -1,
			}`
			req, _ := http.NewRequest("POST", "/users/1/sensors", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "Basic YWRtaW46UGFzc3cwcmQ=")
			router.ServeHTTP(w, req)

			assert.Equal(t, 422, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("OPTIONS_users_user_id_sensors_204", func(t *testing.T) {
		router := gin.Default()
		setupRouter(router)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/users/1/sensors", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code, "Получили в ответ не тот код")
		allowed := strings.Split(w.Header().Get("Allow"), ",")
		assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
		assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
		assert.Contains(t, allowed, "HEAD", "В разрешённых методах нет HEAD")
		assert.Contains(t, allowed, "GET", "В разрешённых методах нет GET")
	})

	// Другие методы не поддерживаем.
	t.Run("OTHER_users_user_id_sensors_405", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  int
		}{
			{"PUT", "PUT", 405},
			{"DELETE", "DELETE", 405},
			{"PATCH", "PATCH", 405},
			{"CONNECT", "CONNECT", 405},
			{"TRACE", "TRACE", 405},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := gin.Default()
				setupRouter(router)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.input, "/users", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.want, w.Code, "Получили в ответ не тот код")
				allowed := strings.Split(w.Header().Get("Allow"), ",")
				assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
				assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
				assert.Contains(t, allowed, "HEAD", "В разрешённых методах нет HEAD")
				assert.Contains(t, allowed, "GET", "В разрешённых методах нет GET")
			})
		}
	})
}

// Тесты /events
func TestEventsRoutes(t *testing.T) {
	t.Run("POST_events", func(t *testing.T) {
		t.Run("valid_request_201", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"sensor_id": 1,
				"sensor_serial_number": "1234567890",
				"timestamp": "2024-04-08T11:24:29.747Z",
				"payload": 10
			}`
			req, _ := http.NewRequest("POST", "/events", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 201, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_has_unsupported_format_415", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `<SensorEvent>
				<SensorId>1</SensorId>
				<SensorSerialNumber>1234567890</SensorSerialNumber>
				<Timestamp>2024-04-08T11:24:29.747Z</Timestamp>
				<Payload>10</Payload>
			</SensorEvent>`
			req, _ := http.NewRequest("POST", "/events", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/xml")
			router.ServeHTTP(w, req)

			assert.Equal(t, 415, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_has_syntax_error_400", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{ невалидный json }`
			req, _ := http.NewRequest("POST", "/events", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 400, w.Code, "Получили в ответ не тот код")
		})

		t.Run("request_body_is_valid_but_it_has_invalid_data_422", func(t *testing.T) {
			router := gin.Default()
			setupRouter(router)

			w := httptest.NewRecorder()

			body := `{
				"sensor_id": -1,
				"sensor_serial_number": "1234567890",
				"timestamp": "2024-04-08T11:24:29.747Z",
				"payload": 10
			}`
			req, _ := http.NewRequest("POST", "/events", bytes.NewReader([]byte(body)))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, 422, w.Code, "Получили в ответ не тот код")
		})
	})

	t.Run("OPTIONS_events_204", func(t *testing.T) {
		router := gin.Default()
		setupRouter(router)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/events", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code, "Получили в ответ не тот код")
		allowed := strings.Split(w.Header().Get("Allow"), ",")
		assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
		assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
	})

	// Другие методы не поддерживаем.
	t.Run("OTHER_users_405", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			want  int
		}{
			{"GET", "GET", 405},
			{"PUT", "PUT", 405},
			{"DELETE", "DELETE", 405},
			{"HEAD", "HEAD", 405},
			{"PATCH", "PATCH", 405},
			{"CONNECT", "CONNECT", 405},
			{"TRACE", "TRACE", 405},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				router := gin.Default()
				setupRouter(router)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest(tt.input, "/users", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.want, w.Code, "Получили в ответ не тот код")
				allowed := strings.Split(w.Header().Get("Allow"), ",")
				assert.Contains(t, allowed, "OPTIONS", "В разрешённых методах нет OPTIONS")
				assert.Contains(t, allowed, "POST", "В разрешённых методах нет POST")
			})
		}
	})
}
