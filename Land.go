package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const LAND_OFFSET = 0x173A // 5946
const LAND_LENGTH = 24
const LAND_MAX_INDEX = 49

type LandInfo struct {
	PropertyMap map[int]LandProperty
}

func newLandInfo() LandInfo {
	var landInfo LandInfo
	landInfo.PropertyMap = make(map[int]LandProperty)

	return landInfo
}

func (l LandInfo) String() string {
	var b bytes.Buffer
	for k, v := range l.PropertyMap {
		b.WriteString(fmt.Sprintf("Land:[%2d][%s]\n", k, v))
	}

	return b.String()
}

func (l *LandInfo) Read(saveIndex int, f *os.File) {
	landOffset := HEADER + ((saveIndex - 1) * SAVE_OFFSET) + LAND_OFFSET
	fmt.Printf("Land Offset : %x, %d\n", landOffset, landOffset)

	seekOffset, err := f.Seek(int64(landOffset), 0)
	if err != nil {
		fmt.Println("file seek error")
		panic(err)
	}
	fmt.Printf("offset : %d\n", seekOffset)

	for i := 1; i <= LAND_MAX_INDEX; i++ {
		byteArray := make([]byte, LAND_LENGTH)
		readSize, err := f.Read(byteArray)
		if len(byteArray) != readSize || err != nil {
			fmt.Println("file read error")
			panic(err)
		}

		l.PropertyMap[i] = newLandProperty(byteArray)
	}
}

func (l *LandInfo) printNeighbor(landNo int) {
	fmt.Printf("Current Land : %d\n", landNo)

	PrintHeader()
	for _, neighborLandNo := range NeighborMap[landNo] {
		fmt.Printf("Land:[%2d][%s]\n", neighborLandNo, l.PropertyMap[neighborLandNo])
	}
}

func PrintHeader() {
	fmt.Printf("LAND     [%5s %5s %5s %5s %5s %4s %4s %3s %3s %3s %3s %3s %3s %3s %3s %5s]\n",
		"HERO", "PEOPL", "ETC",
		"MONEY", "FOOD", "STEE", "FUR", "RAT",
		"FLO", "LAN", "WEA", "SUP", "ARM",
		"SKI", "ETC", "OWNER")
}

type LandProperty struct {
	Hero   uint16
	Poeple uint16
	Etc4   uint16

	Money uint16
	Food  uint16
	Steel uint16
	Fur   uint16
	Rate  uint8

	Flood   uint8
	Land    uint8
	Wealth  uint8
	Support uint8
	Arms    uint8

	Skill uint8
	Etc1  uint8
	Ruler uint16
}

func (p LandProperty) String() string {
	return fmt.Sprintf("%4x, %4x, %4x, %4d, %4d, %3d, %3d, %2d, %2d, %2d, %2d, %2d, %2d, %2d, %2x, %4x",
		p.Hero, p.Poeple, p.Etc4,
		p.Money, p.Food, p.Steel, p.Fur, p.Rate,
		p.Flood, p.Land, p.Wealth, p.Support, p.Arms,
		p.Skill, p.Etc1, p.Ruler)
}

func newLandProperty(byteArray []byte) LandProperty {
	landProperty := LandProperty{}
	buf := bytes.NewReader(byteArray)
	err := binary.Read(buf, binary.LittleEndian, &landProperty)
	if err != nil {
		fmt.Printf("LandProperty Serialize Error %s\n", byteArray)
	}

	return landProperty
}
