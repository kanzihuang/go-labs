package simplerpc

import (
	"encoding/json"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

func NewClient(network, address string) (*Client, error) {
	conn, err := net.DialTimeout(network, address, time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Call(req *ProxyReq) *ProxyResp {
	encoder := json.NewEncoder(c.conn)
	err := encoder.Encode(req)
	if err != nil {
		return &ProxyResp{
			Data: nil,
			Err:  NewError(err),
		}
	}
	decoder := json.NewDecoder(c.conn)
	resp := &ProxyResp{}
	err = decoder.Decode(resp)
	if err != nil {
		return &ProxyResp{
			Data: nil,
			Err:  NewError(err),
		}
	}
	return resp
}