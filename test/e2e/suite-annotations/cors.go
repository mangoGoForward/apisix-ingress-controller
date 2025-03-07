// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package annotations

import (
	"fmt"
	"net/http"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-ingress-controller/test/e2e/scaffold"
)

var _ = ginkgo.Describe("suite-annotations: cors annotations", func() {
	s := scaffold.NewDefaultScaffold()

	ginkgo.It("enable in ingress networking/v1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/enable-cors: "true"
    k8s.apisix.apache.org/cors-allow-origin: https://foo.com,https://bar.com
    k8s.apisix.apache.org/cors-allow-headers: x-foo-1,x-foo-2
    k8s.apisix.apache.org/cors-allow-methods: GET,POST,PUT
  name: ingress-v1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          service:
            name: %s
            port:
              number: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		resp := s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()

		resp = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://baz.com").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()

		resp = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://foo.com").Expect()
		resp.Status(http.StatusOK)
		resp.Header("Access-Control-Allow-Origin").Equal("https://foo.com")
		resp.Header("Access-Control-Allow-Methods").Equal("GET,POST,PUT")
		resp.Header("Access-Control-Allow-Headers").Equal("x-foo-1,x-foo-2")
	})

	ginkgo.It("disable in ingress networking/v1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/enable-cors: "false"
    k8s.apisix.apache.org/cors-allow-origin: https://foo.com,https://bar.com
    k8s.apisix.apache.org/cors-allow-headers: x-foo-1,x-foo-2
    k8s.apisix.apache.org/cors-allow-methods: GET,POST,PUT
  name: ingress-v1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          service:
            name: %s
            port:
              number: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		resp := s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://foo.com").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()
	})

	ginkgo.It("enable in ingress networking/v1beta1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/enable-cors: "true"
    k8s.apisix.apache.org/cors-allow-origin: https://foo.com,https://bar.com
    k8s.apisix.apache.org/cors-allow-headers: x-foo-1,x-foo-2
    k8s.apisix.apache.org/cors-allow-methods: GET,POST,PUT
  name: ingress-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		resp := s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()

		resp = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://baz.com").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()

		resp = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://foo.com").Expect()
		resp.Status(http.StatusOK)
		resp.Header("Access-Control-Allow-Origin").Equal("https://foo.com")
		resp.Header("Access-Control-Allow-Methods").Equal("GET,POST,PUT")
		resp.Header("Access-Control-Allow-Headers").Equal("x-foo-1,x-foo-2")
	})

	ginkgo.It("disable in ingress networking/v1beta1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/enable-cors: "false"
    k8s.apisix.apache.org/cors-allow-origin: https://foo.com,https://bar.com
    k8s.apisix.apache.org/cors-allow-headers: x-foo-1,x-foo-2
    k8s.apisix.apache.org/cors-allow-methods: GET,POST,PUT
  name: ingress-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		resp := s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()

		resp = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://foo.com").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()
	})

	ginkgo.It("enable in ingress extensions/v1beta1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/enable-cors: "true"
    k8s.apisix.apache.org/cors-allow-origin: https://foo.com,https://bar.com
    k8s.apisix.apache.org/cors-allow-headers: x-foo-1,x-foo-2
    k8s.apisix.apache.org/cors-allow-methods: GET,POST,PUT
  name: ingress-extensions-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		resp := s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()

		resp = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://baz.com").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()

		resp = s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://foo.com").Expect()
		resp.Status(http.StatusOK)
		resp.Header("Access-Control-Allow-Origin").Equal("https://foo.com")
		resp.Header("Access-Control-Allow-Methods").Equal("GET,POST,PUT")
		resp.Header("Access-Control-Allow-Headers").Equal("x-foo-1,x-foo-2")
	})

	ginkgo.It("disable in ingress extensions/v1beta1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/enable-cors: "false"
    k8s.apisix.apache.org/cors-allow-origin: https://foo.com,https://bar.com
    k8s.apisix.apache.org/cors-allow-headers: x-foo-1,x-foo-2
    k8s.apisix.apache.org/cors-allow-methods: GET,POST,PUT
  name: ingress-extensions-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /ip
        pathType: Exact
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		resp := s.NewAPISIXClient().GET("/ip").WithHeader("Host", "httpbin.org").WithHeader("Origin", "https://foo.com").Expect()
		resp.Status(http.StatusOK)
		// As httpbin itself adds this header, we don't check it here.
		// resp.Header("Access-Control-Allow-Origin").Empty()
		resp.Header("Access-Control-Allow-Methods").Empty()
		resp.Header("Access-Control-Allow-Headers").Empty()
	})
})
