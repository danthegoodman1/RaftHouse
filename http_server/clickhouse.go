package http_server

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s *HTTPServer) HTTPInfo(c *CustomContext) error {
	str := strings.Builder{}
	str.WriteString("HTTP Request:\n")
	str.WriteString("\n")
	str.WriteString("\tProto: " + c.Request().Proto)
	str.WriteString("\n")
	str.WriteString("\tMethod: " + c.Request().Method)
	str.WriteString("\n")
	str.WriteString("\tHeaders:\n")
	for key, vals := range c.Request().Header {
		for _, val := range vals {
			str.WriteString(fmt.Sprintf("\t\t%s: %s\n", key, val))
		}
	}
	str.WriteString("\n")
	str.WriteString("\tQuery:\n")
	for key, vals := range c.Request().URL.Query() {
		for _, val := range vals {
			if replayHeader := c.Request().Header.Get("X-Replayed"); key == "replay_lim" && replayHeader < val {
				fmt.Println("replaying with lim "+val, replayHeader)
				c.Response().Header().Set("x-replay", "http://control-plane:8080") // add self replay
			}
			str.WriteString(fmt.Sprintf("\t\t%s: %s\n", key, val))
		}
	}

	str.WriteString("\tBody:\n")
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("error in io.ReadAll(c.Request().Body: %w", err)
	}
	str.WriteString(fmt.Sprintf("\t\t%s\n", string(body)))

	fmt.Println(str.String())

	return c.NoContent(http.StatusNoContent)
}
