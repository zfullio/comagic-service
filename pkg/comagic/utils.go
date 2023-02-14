package comagic

import (
	"fmt"
	"time"
)

func generateID(method Method) (id string) {
	id = fmt.Sprintf("%s_%s", method, time.Now().Format("2006-01-02 15:04"))
	return
}

func (c *Client) buildLink() (link string) {
	link = fmt.Sprintf("https://%v/v%s", c.host, c.Version)
	return
}
