package audio

import (
	"bytes"
	"encoding/binary"

	"github.com/gordonklaus/portaudio"
)

type player struct {
	stream *portaudio.Stream
	buffer []int16
}

func NewPlayer(sampleRate int, framesPerBuffer int) (*player, error) {
	buffer := make([]int16, framesPerBuffer)

	stream, err := portaudio.OpenDefaultStream(0, 1, float64(sampleRate), framesPerBuffer, buffer)
	if err != nil {
		return nil, err
	}

	if err := stream.Start(); err != nil {
		stream.Close()
		return nil, err
	}

	return &player{
		stream: stream,
		buffer: buffer,
	}, nil
}

func (p *player) Write(data []byte) error {
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.LittleEndian, p.buffer); err != nil {
		return err
	}

	return p.stream.Write()
}

func (p *player) Close() error {
	return p.stream.Close()
}