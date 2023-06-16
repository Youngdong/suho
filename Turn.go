package main

import (
	"bytes"
	"fmt"
	"os"
)

const TURN_OFFSET = 0x15F8 // 5624

type TurnInfo struct {
	Index   int
	TurnMap map[int]int
}

func NewTurnInfo() TurnInfo {
	return TurnInfo{Index: 0, TurnMap: make(map[int]int)}
}

func (t *TurnInfo) String() string {
	var byteBuffer bytes.Buffer
	for k, v := range t.TurnMap {
		byteBuffer.WriteString(fmt.Sprintf("Idx[%d]:Land[%d]\n", k, v))
	}

	return fmt.Sprintf("TurnInfo Index:[%d]\nMap:[%s]", t.Index, byteBuffer.String())
}

func (t *TurnInfo) GetCurrentLand() int {
	return t.TurnMap[t.Index]
}

func (t *TurnInfo) Read(saveIndex int, file *os.File) {
	turnOffset := HEADER + ((saveIndex - 1) * SAVE_OFFSET) + TURN_OFFSET
	fmt.Printf("Turn Offset : %x, %d\n", turnOffset, turnOffset)

	seekOffset, err := file.Seek(int64(turnOffset), 0)
	if err != nil {
		fmt.Println("file seek error")
		panic(err)
	}
	fmt.Printf("offset : %x\n", seekOffset)

	byteArray := make([]byte, 51)
	readSize, err := file.Read(byteArray)
	if len(byteArray) != readSize || err != nil {
		fmt.Println("file read error")
		panic(err)
	}

	t.Index = int(byteArray[0]) + 1
	for i := 2; i < 51; i++ {
		t.TurnMap[i-1] = int(byteArray[i]) + 1
	}
}
