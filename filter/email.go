package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/labstack/echo/v4"
)

func BlockByEmail(c echo.Context) bool {
	req := c.Request()
	uri := req.RequestURI

	if uri == "/api/v8/manzana/registration" {
		requestData := make(map[string]interface{})
		b, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(b))
		json.Unmarshal(b, &requestData)

		return isGmail(requestData["EmailAddress"].(string))
	}

	return true
}

func isGmail(email string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@gmail.com$`, email)
	if matched {
		fmt.Println(email)
		return true
	}

	return false
}
