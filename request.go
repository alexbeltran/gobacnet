package gobacnet

import "log"

func (c *Client) sendRequest() error {
	id, err := c.tsm.GetFree()
	if err != nil {
		return err
	}
	// get my address?

	log.Printf("id:%d", id)
	return nil
}
