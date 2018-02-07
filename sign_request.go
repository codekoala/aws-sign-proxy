package aws_sign_proxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"go.uber.org/zap"
)

var log *zap.Logger

func SignRequest(config Config, signer *v4.Signer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			body *bytes.Reader
			keys []string
			err  error
		)

		if buf, err := ioutil.ReadAll(r.Body); err == nil {
			body = bytes.NewReader(buf)
		}

		url := r.URL
		url.Scheme = config.TargetProto
		url.Host = config.TargetHost
		log.Info("proxying request",
			zap.String("client", r.RemoteAddr),
			zap.String("url", url.String()))
		req, err := http.NewRequest(r.Method, url.String(), body)

		// copy headers from the incoming request
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// add any extra headers required for the target service
		for key, value := range config.ExtraHeaders {
			req.Header.Set(key, value)
		}

		// sign the request
		switch req.Method {
		case "POST", "PUT":
			_, err = signer.Sign(req, body, config.Provider, config.Region, time.Now())
		default:
			_, err = signer.Sign(req, nil, config.Provider, config.Region, time.Now())
		}
		if err != nil {
			log.Warn("failed to sign request", zap.Error(err))
		}

		for key := range req.Header {
			keys = append(keys, key)
		}
		log.Info("proxied request headers", zap.Any("headers", keys))
		log.Debug("signed request header", zap.String("Authorization", req.Header.Get("Authorization")))

		// issue the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Warn("failed to issue request", zap.Error(err))
		}
		defer resp.Body.Close()

		// copy the status code from the response
		w.WriteHeader(resp.StatusCode)

		// copy headers from the response
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// copy the body from the response
		if _, err = io.Copy(w, resp.Body); err != nil {
			log.Warn("failed to copy response body", zap.Error(err))
		}
	}
}
