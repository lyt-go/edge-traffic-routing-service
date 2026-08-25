package model

import "bytes"

func ParseBackendFrame(frame []byte) []byte {
	return bytes.TrimSpace(frame)
}
