// Package internal
package internal

import "io"

type Client interface {
	FetchFile(url string) (io.ReadCloser, error)
}
