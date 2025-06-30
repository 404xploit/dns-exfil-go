package main

import (
	"encoding/base32"
	"strings"
)

// EncodeBase32 encodes a byte slice to a base32 string.
func EncodeBase32(data []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(data)
}

// DecodeBase32 decodes a base32 string to a byte slice.
func DecodeBase32(s string) ([]byte, error) {
	return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s)
}

// ChunkData divides data into chunks of a specified size.
func ChunkData(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}


