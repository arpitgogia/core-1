package engine

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"coralreefci/engine/frontend"
)

func Test_activationServerHandler(t *testing.T) {
	http.HandleFunc(destinationEnd, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(destinationPort, nil)

	testAS := new(ActivationServer)

	var secrets = []struct {
		code  int
		value string
	}{
		{code: http.StatusForbidden, value: ""},
		{code: http.StatusForbidden, value: "mos-eisley"},
		{code: http.StatusOK, value: frontend.BackendSecret},
	}

	for _, secret := range secrets {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/handler-test", nil)
		if err != nil {
			t.Errorf("failure generating activation service test request: %v", err)
		}
		req.Form = url.Values{}
		req.Form.Set("state", secret.value)
		req.Form.Set("repos", string(94))

		handler := http.HandlerFunc(testAS.activationServerHandler)
		handler.ServeHTTP(rec, req)

		if received := rec.Code; received != secret.code {
			t.Errorf("handler returning incorrect status code; received %v, expected %v", received, secret.code)
		}
		req.Form.Del("state")
		req.Form.Del("repos")
	}
}
