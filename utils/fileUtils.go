package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type BoxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

func GetMP4Duration(reader io.ReaderAt) (lengthOfTime uint32) {
	var info = make([]byte, 0x10)
	var boxHeader BoxHeader
	var offset int64 = 0

	runs := 0
	for { // TODO
		if runs == 10 {
			return
		}
		_, err := reader.ReadAt(info, offset)
		if err != nil {
			log.Println(err)
			return
		}
		boxHeader = getHeaderBoxInfo(info)
		fourccType := getFourccType(boxHeader)
		if fourccType == "moov" {
			break
		}
		if fourccType == "mdat" {
			if boxHeader.Size == 1 {
				offset += int64(boxHeader.Size64)
				continue
			}
		}
		offset += int64(boxHeader.Size)
		runs += 1
	}

	moovStartBytes := make([]byte, 0x100)
	_, err := reader.ReadAt(moovStartBytes, offset)
	if err != nil {
		fmt.Println("err:")
		fmt.Println(err)
		return
	}

	timeScaleOffset := 0x1C
	durationOffest := 0x20
	timeScale := binary.BigEndian.Uint32(moovStartBytes[timeScaleOffset : timeScaleOffset+4])
	Duration := binary.BigEndian.Uint32(moovStartBytes[durationOffest : durationOffest+4])
	lengthOfTime = Duration / timeScale
	return lengthOfTime

}

func getHeaderBoxInfo(data []byte) (boxHeader BoxHeader) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &boxHeader)
	return
}

func getFourccType(boxHeader BoxHeader) (fourccType string) {
	fourccType = string(boxHeader.FourccType[:])
	return
}
