package testdata

import (
	"bytes"
	"encoding/binary"
)

// CreateValidMP4WithResolution creates a minimal but valid MP4 file with specified resolution
func CreateValidMP4WithResolution(width, height uint32) []byte {
	buf := &bytes.Buffer{}
	
	// Write ftyp box
	writeFtypBox(buf)
	
	// Write moov box with trak
	writeMoovBox(buf, width, height)
	
	return buf.Bytes()
}

// CreateValidMP4 creates a minimal but valid MP4 file with 1920x1080 resolution
func CreateValidMP4() []byte {
	return CreateValidMP4WithResolution(1920, 1080)
}

// CreateInvalidMP4 creates an invalid MP4 file for testing error cases
func CreateInvalidMP4() []byte {
	return []byte("invalid mp4 content")
}

func writeFtypBox(buf *bytes.Buffer) {
	// ftyp box
	binary.Write(buf, binary.BigEndian, uint32(32)) // box size
	buf.WriteString("ftyp")                         // box type
	buf.WriteString("isom")                         // major brand
	binary.Write(buf, binary.BigEndian, uint32(512)) // minor version
	buf.WriteString("isom")                         // compatible brand 1
	buf.WriteString("iso2")                         // compatible brand 2
	buf.WriteString("avc1")                         // compatible brand 3
	buf.WriteString("mp41")                         // compatible brand 4
}

func writeMoovBox(buf *bytes.Buffer, width, height uint32) {
	moovStart := buf.Len()
	
	// Placeholder for moov box size
	binary.Write(buf, binary.BigEndian, uint32(0))
	buf.WriteString("moov")
	
	// Write mvhd box
	writeMvhdBox(buf)
	
	// Write trak box
	writeTrakBox(buf, width, height)
	
	// Update moov box size
	moovEnd := buf.Len()
	moovSize := uint32(moovEnd - moovStart)
	binary.BigEndian.PutUint32(buf.Bytes()[moovStart:moovStart+4], moovSize)
}

func writeMvhdBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(108)) // box size
	buf.WriteString("mvhd")                          // box type
	binary.Write(buf, binary.BigEndian, uint32(0))   // version + flags
	binary.Write(buf, binary.BigEndian, uint32(0))   // creation time
	binary.Write(buf, binary.BigEndian, uint32(0))   // modification time
	binary.Write(buf, binary.BigEndian, uint32(1000)) // timescale
	binary.Write(buf, binary.BigEndian, uint32(0))   // duration
	binary.Write(buf, binary.BigEndian, uint32(0x00010000)) // rate
	binary.Write(buf, binary.BigEndian, uint16(0x0100)) // volume
	binary.Write(buf, binary.BigEndian, uint16(0))     // reserved
	binary.Write(buf, binary.BigEndian, uint64(0))     // reserved
	// Matrix (36 bytes)
	for i := 0; i < 9; i++ {
		if i == 0 || i == 4 || i == 8 {
			binary.Write(buf, binary.BigEndian, uint32(0x00010000))
		} else {
			binary.Write(buf, binary.BigEndian, uint32(0))
		}
	}
	// Pre-defined (24 bytes)
	for i := 0; i < 6; i++ {
		binary.Write(buf, binary.BigEndian, uint32(0))
	}
	binary.Write(buf, binary.BigEndian, uint32(2)) // next track ID
}

func writeTrakBox(buf *bytes.Buffer, width, height uint32) {
	trakStart := buf.Len()
	
	// Placeholder for trak box size
	binary.Write(buf, binary.BigEndian, uint32(0))
	buf.WriteString("trak")
	
	// Write tkhd box
	writeTkhdBox(buf, width, height)
	
	// Write mdia box
	writeMdiaBox(buf, width, height)
	
	// Update trak box size
	trakEnd := buf.Len()
	trakSize := uint32(trakEnd - trakStart)
	binary.BigEndian.PutUint32(buf.Bytes()[trakStart:trakStart+4], trakSize)
}

func writeTkhdBox(buf *bytes.Buffer, width, height uint32) {
	binary.Write(buf, binary.BigEndian, uint32(92)) // box size
	buf.WriteString("tkhd")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(7))  // version + flags (track enabled)
	binary.Write(buf, binary.BigEndian, uint32(0))  // creation time
	binary.Write(buf, binary.BigEndian, uint32(0))  // modification time
	binary.Write(buf, binary.BigEndian, uint32(1))  // track ID
	binary.Write(buf, binary.BigEndian, uint32(0))  // reserved
	binary.Write(buf, binary.BigEndian, uint32(0))  // duration
	binary.Write(buf, binary.BigEndian, uint64(0))  // reserved
	binary.Write(buf, binary.BigEndian, uint16(0))  // layer
	binary.Write(buf, binary.BigEndian, uint16(0))  // alternate group
	binary.Write(buf, binary.BigEndian, uint16(0))  // volume
	binary.Write(buf, binary.BigEndian, uint16(0))  // reserved
	// Matrix (36 bytes)
	for i := 0; i < 9; i++ {
		if i == 0 || i == 4 || i == 8 {
			binary.Write(buf, binary.BigEndian, uint32(0x00010000))
		} else {
			binary.Write(buf, binary.BigEndian, uint32(0))
		}
	}
	binary.Write(buf, binary.BigEndian, uint32(width<<16))  // width
	binary.Write(buf, binary.BigEndian, uint32(height<<16)) // height
}

func writeMdiaBox(buf *bytes.Buffer, width, height uint32) {
	mdiaStart := buf.Len()
	
	// Placeholder for mdia box size
	binary.Write(buf, binary.BigEndian, uint32(0))
	buf.WriteString("mdia")
	
	// Write mdhd box
	writeMdhdBox(buf)
	
	// Write hdlr box
	writeHdlrBox(buf)
	
	// Write minf box
	writeMinfBox(buf, width, height)
	
	// Update mdia box size
	mdiaEnd := buf.Len()
	mdiaSize := uint32(mdiaEnd - mdiaStart)
	binary.BigEndian.PutUint32(buf.Bytes()[mdiaStart:mdiaStart+4], mdiaSize)
}

func writeMdhdBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(32)) // box size
	buf.WriteString("mdhd")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(0))  // version + flags
	binary.Write(buf, binary.BigEndian, uint32(0))  // creation time
	binary.Write(buf, binary.BigEndian, uint32(0))  // modification time
	binary.Write(buf, binary.BigEndian, uint32(1000)) // timescale
	binary.Write(buf, binary.BigEndian, uint32(0))  // duration
	binary.Write(buf, binary.BigEndian, uint16(0x55c4)) // language (und)
	binary.Write(buf, binary.BigEndian, uint16(0))  // pre-defined
}

func writeHdlrBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(33)) // box size
	buf.WriteString("hdlr")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(0))  // version + flags
	binary.Write(buf, binary.BigEndian, uint32(0))  // pre-defined
	buf.WriteString("vide")                         // handler type (video)
	binary.Write(buf, binary.BigEndian, uint32(0))  // reserved
	binary.Write(buf, binary.BigEndian, uint32(0))  // reserved
	binary.Write(buf, binary.BigEndian, uint32(0))  // reserved
	buf.WriteByte(0)                                // name (empty string)
}

func writeMinfBox(buf *bytes.Buffer, width, height uint32) {
	minfStart := buf.Len()
	
	// Placeholder for minf box size
	binary.Write(buf, binary.BigEndian, uint32(0))
	buf.WriteString("minf")
	
	// Write vmhd box
	writeVmhdBox(buf)
	
	// Write dinf box
	writeDinfBox(buf)
	
	// Write stbl box
	writeStblBox(buf, width, height)
	
	// Update minf box size
	minfEnd := buf.Len()
	minfSize := uint32(minfEnd - minfStart)
	binary.BigEndian.PutUint32(buf.Bytes()[minfStart:minfStart+4], minfSize)
}

func writeVmhdBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(20)) // box size
	buf.WriteString("vmhd")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(1))  // version + flags
	binary.Write(buf, binary.BigEndian, uint16(0))  // graphics mode
	binary.Write(buf, binary.BigEndian, uint16(0))  // opcolor R
	binary.Write(buf, binary.BigEndian, uint16(0))  // opcolor G
	binary.Write(buf, binary.BigEndian, uint16(0))  // opcolor B
}

func writeDinfBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(36)) // box size
	buf.WriteString("dinf")                         // box type
	
	// Write dref box
	binary.Write(buf, binary.BigEndian, uint32(28)) // dref box size
	buf.WriteString("dref")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(0))  // version + flags
	binary.Write(buf, binary.BigEndian, uint32(1))  // entry count
	
	// Write url box
	binary.Write(buf, binary.BigEndian, uint32(12)) // url box size
	buf.WriteString("url ")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(1))  // version + flags (self-contained)
}

func writeStblBox(buf *bytes.Buffer, width, height uint32) {
	stblStart := buf.Len()
	
	// Placeholder for stbl box size
	binary.Write(buf, binary.BigEndian, uint32(0))
	buf.WriteString("stbl")
	
	// Write stsd box
	writeStsdBox(buf, width, height)
	
	// Write stts box
	writeSttsBox(buf)
	
	// Write stsc box
	writeStscBox(buf)
	
	// Write stsz box
	writeStszBox(buf)
	
	// Write stco box
	writeStcoBox(buf)
	
	// Update stbl box size
	stblEnd := buf.Len()
	stblSize := uint32(stblEnd - stblStart)
	binary.BigEndian.PutUint32(buf.Bytes()[stblStart:stblStart+4], stblSize)
}

func writeStsdBox(buf *bytes.Buffer, width, height uint32) {
	stsdStart := buf.Len()
	
	// Placeholder for stsd box size
	binary.Write(buf, binary.BigEndian, uint32(0))
	buf.WriteString("stsd")
	binary.Write(buf, binary.BigEndian, uint32(0)) // version + flags
	binary.Write(buf, binary.BigEndian, uint32(1)) // entry count
	
	// Write avc1 sample entry
	avc1Start := buf.Len()
	binary.Write(buf, binary.BigEndian, uint32(0)) // placeholder for avc1 size
	buf.WriteString("avc1")                        // sample entry type
	
	// Reserved (6 bytes)
	for i := 0; i < 6; i++ {
		buf.WriteByte(0)
	}
	binary.Write(buf, binary.BigEndian, uint16(1)) // data reference index
	
	// Video sample entry fields
	binary.Write(buf, binary.BigEndian, uint16(0)) // pre-defined
	binary.Write(buf, binary.BigEndian, uint16(0)) // reserved
	binary.Write(buf, binary.BigEndian, uint32(0)) // pre-defined
	binary.Write(buf, binary.BigEndian, uint32(0)) // pre-defined
	binary.Write(buf, binary.BigEndian, uint32(0)) // pre-defined
	binary.Write(buf, binary.BigEndian, uint16(width))  // width
	binary.Write(buf, binary.BigEndian, uint16(height)) // height
	binary.Write(buf, binary.BigEndian, uint32(0x00480000)) // horizontal resolution
	binary.Write(buf, binary.BigEndian, uint32(0x00480000)) // vertical resolution
	binary.Write(buf, binary.BigEndian, uint32(0)) // reserved
	binary.Write(buf, binary.BigEndian, uint16(1)) // frame count
	
	// Compressor name (32 bytes)
	for i := 0; i < 32; i++ {
		buf.WriteByte(0)
	}
	
	binary.Write(buf, binary.BigEndian, uint16(24)) // depth
	binary.Write(buf, binary.BigEndian, uint16(0xFFFF)) // pre-defined
	
	// Update avc1 size
	avc1End := buf.Len()
	avc1Size := uint32(avc1End - avc1Start)
	binary.BigEndian.PutUint32(buf.Bytes()[avc1Start:avc1Start+4], avc1Size)
	
	// Update stsd size
	stsdEnd := buf.Len()
	stsdSize := uint32(stsdEnd - stsdStart)
	binary.BigEndian.PutUint32(buf.Bytes()[stsdStart:stsdStart+4], stsdSize)
}

func writeSttsBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(16)) // box size
	buf.WriteString("stts")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(0))  // version + flags
	binary.Write(buf, binary.BigEndian, uint32(0))  // entry count
}

func writeStscBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(16)) // box size
	buf.WriteString("stsc")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(0))  // version + flags
	binary.Write(buf, binary.BigEndian, uint32(0))  // entry count
}

func writeStszBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(20)) // box size
	buf.WriteString("stsz")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(0))  // version + flags
	binary.Write(buf, binary.BigEndian, uint32(0))  // sample size
	binary.Write(buf, binary.BigEndian, uint32(0))  // sample count
}

func writeStcoBox(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, uint32(16)) // box size
	buf.WriteString("stco")                         // box type
	binary.Write(buf, binary.BigEndian, uint32(0))  // version + flags
	binary.Write(buf, binary.BigEndian, uint32(0))  // entry count
}