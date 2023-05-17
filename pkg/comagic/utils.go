package comagic

import (
	"fmt"
	"time"
)

func generateID(method Method) string {
	return fmt.Sprintf("%s_%s", method, time.Now().Format("2006-01-02 15:04"))
}

func (c *Client) buildLink() string {
	return fmt.Sprintf("https://%v/v%s", c.host, c.Version)
}
