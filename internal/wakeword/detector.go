package wakeword

import (
	"bytes"
	"encoding/binary"
	"os"

	porcupine "github.com/Picovoice/porcupine/binding/go/v3"
)

type Detector struct {
	porcupine *porcupine.Porcupine
	keywordID int32
}

func NewDetector(modelPath, keywordPath string, sensitivity float32) (*Detector, error) {
	p := porcupine.Porcupine{
		AccessKey: os.Getenv("PORCUPINE_ACCESS_KEY"),
		KeywordPaths: []string{"models/Helix_en_mac_v3_0_0.ppn"},
	}
	err := p.Init()
	if err != nil {
		return nil, err
	}

	return &Detector{
		porcupine:  &p,
		keywordID:  0,
	}, nil
}

func (d *Detector) Process(audio []byte) (bool, error) {
	var samples []int16
	reader := bytes.NewReader(audio)
	for reader.Len() > 0 {
		var sample int16
		if err := binary.Read(reader, binary.LittleEndian, &sample); err != nil {
			return false, err
		}
		samples = append(samples, sample)
	}

	result, err := d.porcupine.Process(samples)
	if err != nil {
		return false, err
	}

	return result >= 0, nil
}

func (d *Detector) Close() error {
	return d.porcupine.Delete()
}

