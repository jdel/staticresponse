package staticresponse_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jdel/staticresponse"
)

func TestStaticResponse(t *testing.T) {
	testCases := []struct {
		desc          string
		cfg           *staticresponse.Config
		expectedError bool
	}{
		{
			desc: "default config",
			cfg:  staticresponse.CreateConfig(),
		},
		{
			desc: "custom config teapot multi-value header",
			cfg: &staticresponse.Config{
				StatusCode: 418,
				Body:       "I'm a teapot",
				Headers:    http.Header{"Cake": {"One", "Two"}},
			},
		},
		{
			desc: "custom config json multiple headers",
			cfg: &staticresponse.Config{
				StatusCode: 200,
				Body:       "{\"statusCode\": \"OK\"}",
				Headers:    http.Header{"One": {"1"}, "Two": {"2"}},
			},
		},
		{
			desc: "custom config 200 no body",
			cfg:  &staticresponse.Config{StatusCode: 200},
		},
		{
			desc:          "custom config 1000 error",
			cfg:           &staticresponse.Config{StatusCode: 1000, Body: "NOK"},
			expectedError: true,
		},
	}
	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			// t.Parallel()

			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			handler, handlerErr := staticresponse.New(ctx, next, test.cfg, "staticresponse")
			if !test.expectedError {
				recorder := httptest.NewRecorder()

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
				if err != nil {
					t.Fatal(err)
				}

				handler.ServeHTTP(recorder, req)

				if recorder.Result().StatusCode != test.cfg.StatusCode {
					t.Errorf("invalid response code: %d (expected %d)", recorder.Result().StatusCode, test.cfg.StatusCode)
				}

				if b, err := io.ReadAll(recorder.Result().Body); err == nil && string(b) != test.cfg.Body {
					t.Errorf("invalid response body: %s (expected %s)", string(b), test.cfg.Body)
				}

				if !reflect.DeepEqual(recorder.Result().Header, test.cfg.Headers) {
					t.Errorf("headers mismatch: %v (expected %v)", test.cfg.Headers, recorder.Result().Header)
				}
			} else if handlerErr == nil {
				t.Fatal("error expected")
			}
		})
	}
}
