package rest_test

import (
	. "github.com/datianshi/rest-func/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/url"
	"bytes"
	"io"
	"fmt"
	"os"
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
})
