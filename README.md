# rest-func

## A library to make golang restful calls easier and readable

## Example

```
import "github.com/datianshi/rest-func/rest"

r := &rest.Rest{
  URL: fmt.Sprintf("%s%s", p.OpsManUrl, uploadUrl),
}
var response *http.Response
response,err = r.Build().
  WithHttpMethod(rest.POST).
  SkipSslVerify(p.SkipSsl).
  WithHttpHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
  WithMultipartForm("form Name", file).
  Connect()
```

## Build Methods

```
WithHttpMethod(method httpMethod)
WithHttpHeader(key string, value string)
WithContentType(value string)
WithBasicAuth(user string, password string)
WithMultipartForm(paramName string, file *os.File)
WithHttpBody(body io.ReadCloser)
WithFormValue(values url.Values)
SkipSslVerify(skip bool)
```
