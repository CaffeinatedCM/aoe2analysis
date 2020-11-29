package parser

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"strings"
)

func readChunk(reader io.Reader, chunkLen int) ([]byte, error){
	buf := make([]byte, chunkLen)
	read, err := reader.Read(buf)
	if err != nil || read != chunkLen {
		return nil, err
	} else if read != chunkLen {
		return nil, errors.New("read wrong number of bytes, expected " + (string)(chunkLen) + " got " + (string)(read))
	}

	return buf, nil
}

func readCString(reader io.Reader) (*string, error) {
	stringBuilder := strings.Builder{}

	buf := make([]byte, 1)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		if buf[0] == 0x00 {
			break
		}
		stringBuilder.WriteByte(buf[0])
	}

	result := stringBuilder.String()
	return &result, nil
}

func readHeader(reader io.Reader) (*Header, error) {
	header := Header{}
	headerLenBuf, err := readChunk(reader, 4)
	if err != nil {
		return nil, errors.New("could not read header length " + err.Error())
	}
	header.length = (int32)(binary.LittleEndian.Uint32(headerLenBuf))

	compressedHeaderBytes, err := readChunk(reader, (int)(header.length) - 8)
	if err != nil {
		return nil, errors.New("could not read compressed header: " + err.Error())
	}

	headerReader := flate.NewReader(bytes.NewReader(compressedHeaderBytes))
	defer func() {
		err := headerReader.Close()
		if err != nil {
			log.Fatalf("could not close header reader")
		}
	}()

	version, err := readCString(headerReader)
	if err != nil {
		return nil, errors.New("failed to read version: " + err.Error())
	}
	header.Version = *version

	return &header, nil
}

func Parse(reader io.Reader) (*RecordedGame, error) {
	game := RecordedGame{}

	header, err := readHeader(reader)
	if err != nil {
		return nil, errors.New("could not read header: " + err.Error())
	}
	game.Header = *header

	return &game, nil
}
