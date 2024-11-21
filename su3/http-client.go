package su3

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

func InitClient() *http.Client {
	dc := func(ctx context.Context, network, addr string) (net.Conn, error) {
		fmt.Printf("preparing a new connection to %s %s\n", network, addr)

		dialer := net.Dialer{
			Timeout: time.Minute,
			ControlContext: func(ctx context.Context, network, addr string, c syscall.RawConn) error {
				fmt.Printf("created connection to %s %s\n", network, addr)

				return nil
			},
		}

		return dialer.DialContext(ctx, network, addr)
	}

	c := http.Client{
		Transport: &http.Transport{
			DialContext:           dc,
			MaxIdleConnsPerHost:   3,
			MaxIdleConns:          12,
			MaxConnsPerHost:       3,
			IdleConnTimeout:       time.Minute,
			ResponseHeaderTimeout: time.Second,
			WriteBufferSize:       8 << 10,
			ReadBufferSize:        8 << 10,
		},
		CheckRedirect: func(next *http.Request, history []*http.Request) error {
			if len(history) > 3 {
				return http.ErrUseLastResponse
			}

			_, _ = fmt.Fprintf(
				os.Stdout, "next request: %s\n", next.URL.String(),
			)

			_, _ = fmt.Fprintf(
				os.Stdout, "list of requests done before %s\n:", next.URL.String(),
			)

			for i := range history {
				_, _ = fmt.Fprintf(
					os.Stdout, "request #%d - %s\n", i+1, history[i].URL.String(),
				)
			}

			return nil
		},
		Timeout: time.Second * 30,
	}

	return &c
}
