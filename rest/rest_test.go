package rest_test

import (
	. "github.com/datianshi/rest-func/rest"

	"bytes"
	"fmt"
	"io"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("Rest", func() {
	var rest Rest

	BeforeEach(func() {
		rest = Rest{
			URL: "http://example.com",
		}
	})

	Context("With Default request", func() {
		connect := rest.Build()
		It("Request should be default as GET Method", func() {
			Ω(connect.Request.Method).Should(Equal("GET"))
		})
		It("Request Body shoud be default nil", func() {
			Ω(connect.Request.Body).Should(BeNil())
		})
	})

	Context("With Authentication Header", func() {
		connect := rest.Build().WithBasicAuth("username", "password")
		It("Request Authentication header should Not be nil", func() {
			Ω(connect.Request.Header.Get("Authorization")).ShouldNot(BeNil())
		})
		It("Request Authentication header should be set correctly", func() {
			Ω(connect.Request.Header.Get("Authorization")).Should(Equal("Basic dXNlcm5hbWU6cGFzc3dvcmQ="))
		})
	})

	Context("With Content Type", func() {
		connect := rest.Build().WithContentType("customType")
		It("Request Content Type Header should be set correctly", func() {
			Ω(connect.Request.Header.Get("Content-Type")).Should(Equal("customType"))
		})
	})

	Context("With Form Values", func() {
		connect := rest.Build().WithFormValue(url.Values{"username": {"abc"}, "password": {"abc"}})
		It("Request method should be POST", func() {
			Ω(connect.Request.Method).Should(Equal("POST"))
		})
		It("Request Content type should be application/x-www-form-urlencoded", func() {
			Ω(connect.Request.Header.Get("Content-Type")).Should(Equal("application/x-www-form-urlencoded"))
		})
		It("Request forms value should be correct", func() {
			writer := bytes.Buffer{}
			io.Copy(&writer, connect.Request.Body)
			fmt.Println(writer.String())
		})
	})

	Context("With MultiFormUpload", func() {
		file, _ := os.Open("fixtures/upload.txt")
		defer file.Close()
		connect := rest.Build().WithMultipartForm("filename", file)
		It("Request method should be POST", func() {
			Ω(connect.Request.Method).Should(Equal("POST"))
		})
	})

	Context("With http proxy", func(){
		var server *ghttp.Server
		var proxyCatchURL string
		var catchProxyURL http.HandlerFunc = func(response http.ResponseWriter, request *http.Request){
			proxyCatchURL = request.URL.String()
		}
		BeforeEach(func() {
			server = ghttp.NewServer()
			os.Setenv("http_proxy", server.URL())
			server.AppendHandlers(ghttp.CombineHandlers(
				catchProxyURL,
			))
		})
		It("Should request through the proxy", func(){
			google := Rest{
				URL: "http://google.com",
			}
			google.Build().SkipSslVerify(true).Connect()
			Ω(proxyCatchURL).Should(Equal("http://google.com/"))
		})
		AfterEach(func() {
			server.Close()
		})
	})


})
