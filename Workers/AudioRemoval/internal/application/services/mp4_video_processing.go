package services

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type MP4VideoProcessingService struct{}

func NewMP4VideoProcessingService() *MP4VideoProcessingService {
	return &MP4VideoProcessingService{}
}

func (s *MP4VideoProcessingService) RemoveAudio(inputData []byte) ([]byte, error) {
	reader := bytes.NewReader(inputData)
	var output bytes.Buffer

	for {
		var boxSize uint32
		var boxType [4]byte

		if err := binary.Read(reader, binary.BigEndian, &boxSize); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read box size: %w", err)
		}

		if _, err := reader.Read(boxType[:]); err != nil {
			return nil, fmt.Errorf("failed to read box type: %w", err)
		}

		boxTypeStr := string(boxType[:])
		contentSize := boxSize - 8
		content := make([]byte, contentSize)
		if _, err := reader.Read(content); err != nil {
			return nil, fmt.Errorf("failed to read box content: %w", err)
		}

		// Procesar según tipo de box
		switch boxTypeStr {
		case "trak":
			if s.isVideoTrack(content) {
				s.writeBox(&output, boxSize, boxType, content)
			}
		case "moov":
			processedMoov := s.processMoovBox(content)
			newSize := uint32(len(processedMoov) + 8)
			s.writeBox(&output, newSize, boxType, processedMoov)
		default:
			s.writeBox(&output, boxSize, boxType, content)
		}
	}

	return output.Bytes(), nil
}

func (s *MP4VideoProcessingService) writeBox(output *bytes.Buffer, size uint32, boxType [4]byte, content []byte) {
	binary.Write(output, binary.BigEndian, size)
	output.Write(boxType[:])
	output.Write(content)
}

func (s *MP4VideoProcessingService) processMoovBox(moovData []byte) []byte {
	reader := bytes.NewReader(moovData)
	var output bytes.Buffer

	for {
		var boxSize uint32
		var boxType [4]byte

		if err := binary.Read(reader, binary.BigEndian, &boxSize); err != nil {
			break
		}

		if _, err := reader.Read(boxType[:]); err != nil {
			break
		}

		contentSize := boxSize - 8
		content := make([]byte, contentSize)
		reader.Read(content)

		if string(boxType[:]) == "trak" {
			if s.isVideoTrack(content) {
				s.writeBox(&output, boxSize, boxType, content)
			}
		} else {
			s.writeBox(&output, boxSize, boxType, content)
		}
	}

	return output.Bytes()
}

func (s *MP4VideoProcessingService) isVideoTrack(trackData []byte) bool {
	reader := bytes.NewReader(trackData)

	for {
		var boxSize uint32
		var boxType [4]byte

		if err := binary.Read(reader, binary.BigEndian, &boxSize); err != nil {
			break
		}

		if _, err := reader.Read(boxType[:]); err != nil {
			break
		}

		if string(boxType[:]) == "mdia" {
			mdiaContent := make([]byte, boxSize-8)
			reader.Read(mdiaContent)
			return s.isVideoHandler(mdiaContent)
		}

		reader.Seek(int64(boxSize-8), io.SeekCurrent)
	}

	return false
}

func (s *MP4VideoProcessingService) isVideoHandler(mdiaData []byte) bool {
	reader := bytes.NewReader(mdiaData)

	for {
		var boxSize uint32
		var boxType [4]byte

		if err := binary.Read(reader, binary.BigEndian, &boxSize); err != nil {
			break
		}

		if _, err := reader.Read(boxType[:]); err != nil {
			break
		}

		if string(boxType[:]) == "hdlr" {
			hdlrContent := make([]byte, boxSize-8)
			reader.Read(hdlrContent)

			if len(hdlrContent) >= 12 {
				handlerType := string(hdlrContent[8:12])
				return handlerType == "vide"
			}
		}

		reader.Seek(int64(boxSize-8), io.SeekCurrent)
	}

	return false
}
