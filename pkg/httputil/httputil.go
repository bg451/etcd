// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// borrowed from golang/net/context/ctxhttp/cancelreq.go

// Package httputil provides HTTP utility functions.
package httputil

import (
	"io"
	"io/ioutil"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
)

func RequestCanceler(rt http.RoundTripper, req *http.Request) func() {
	ch := make(chan struct{})
	req.Cancel = ch

	return func() {
		close(ch)
	}
}

// GracefulClose drains http.Response.Body until it hits EOF
// and closes it. This prevents TCP/TLS connections from closing,
// therefore available for reuse.
func GracefulClose(resp *http.Response) {
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

func JoinOrCreateSpanFromHeader(op string, h http.Header) opentracing.Span {
	t := opentracing.GlobalTracer()
	sp, err := t.Join(op, opentracing.TextMap, opentracing.HTTPHeaderTextMapCarrier(h))
	if err != nil {
		return t.StartSpan(op)
	}
	return sp
}
