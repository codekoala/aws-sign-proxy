package aws_sign_proxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"go.uber.org/zap"
)

type RequestSigner struct {
	log    *zap.Logger
	config Config
	signer *v4.Signer
}

func NewRequestSigner(log *zap.Logger, config Config, signer *v4.Signer) *RequestSigner {
	return &RequestSigner{
		log:    log,
		config: config,
		signer: signer,
	}
}

func (rs *RequestSigner) Proxy(w http.ResponseWriter, r *http.Request) {
	var (
		body *bytes.Reader
		keys []string
		err  error
	)

	if buf, err := ioutil.ReadAll(r.Body); err == nil {
		body = bytes.NewReader(buf)
	}

	url := r.URL
	url.Scheme = rs.config.TargetProto
	url.Host = rs.config.TargetHost
	rs.log.Info("proxying request",
		zap.String("client", r.RemoteAddr),
		zap.String("url", url.String()))
	req, err := http.NewRequest(r.Method, url.String(), body)

	rs.CopyOutboundHeaders(r, req)

	// sign the request
	switch req.Method {
	case "POST", "PUT":
		_, err = rs.signer.Sign(req, body, rs.config.Provider, rs.config.Region, time.Now())
	default:
		_, err = rs.signer.Sign(req, nil, rs.config.Provider, rs.config.Region, time.Now())
	}
	if err != nil {
		rs.log.Warn("failed to sign request", zap.Error(err))
	}

	for key := range req.Header {
		keys = append(keys, key)
	}
	rs.log.Info("proxied request headers", zap.Any("headers", keys))
	rs.log.Debug("signed request header", zap.String("Authorization", req.Header.Get("Authorization")))

	// issue the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		rs.log.Warn("failed to issue request", zap.Error(err))
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
		rs.log.Warn("failed to copy response body", zap.Error(err))
	}
}

func (rs *RequestSigner) IsBlockedHeader(name string) bool {
	name = strings.ToLower(name)
	for _, blocked := range rs.config.BlockHeaders {
		if strings.ToLower(blocked) == name {
			return true
		}
	}

	return false
}

func (rs *RequestSigner) CopyOutboundHeaders(from, to *http.Request) {
	// copy headers from the incoming request
	for key, values := range from.Header {
		if rs.IsBlockedHeader(key) {
			rs.log.Info("dropping blocked header", zap.String("header", key), zap.Any("value", values))
			continue
		}

		for _, value := range values {
			rs.log.Info("copying header", zap.String("header", key), zap.String("value", value))
			to.Header.Add(key, value)
		}
	}

	// add any extra headers required for the target service
	for key, value := range rs.config.ExtraHeaders {
		rs.log.Info("adding header", zap.String("header", key), zap.String("value", value))
		to.Header.Set(key, value)
	}
}
