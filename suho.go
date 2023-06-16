package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

const HEADER = 0x0C        // 12
const SAVE_OFFSET = 0x1BD3 // 7123

var NeighborMap = map[int][]int{
	1:  {13, 12, 2},
	2:  {1, 12, 3},
	3:  {2, 12, 11, 4},
	4:  {3, 11, 10, 6, 5},
	5:  {4, 6},
	6:  {5, 4, 10, 9, 8, 7},
	7:  {6, 8},
	8:  {6, 9, 22, 7},
	9:  {10, 21, 22, 8, 6},
	10: {4, 11, 20, 21, 9, 6},
	11: {3, 12, 19, 23, 20, 10, 4},
	12: {1, 13, 18, 19, 11, 3, 2},
	13: {14, 17, 18, 12, 1},
	14: {15, 17, 13},
	15: {16, 17, 14},
	16: {26, 17, 15},
	17: {14, 15, 16, 26, 25, 18, 13},
	18: {13, 17, 25, 24, 19, 12},
	19: {12, 18, 24, 23, 11},
	20: {11, 23, 30, 21, 10},
	21: {10, 20, 31, 35, 31, 22, 9},
	22: {9, 21, 31, 32, 8},
	23: {19, 24, 28, 29, 30, 20, 11},
	24: {18, 25, 27, 28, 23, 19},
	25: {17, 26, 27, 24, 18},
	26: {17, 16, 40, 27, 25},
	27: {25, 26, 40, 39, 28, 24},
	28: {27, 39, 38, 29, 23, 24},
	29: {23, 28, 38, 37, 36, 30},
	30: {23, 29, 36, 35, 21, 20},
	31: {21, 35, 47, 49, 34, 32, 22},
	32: {22, 31, 34, 33},
	33: {32, 34},
	34: {31, 39, 33, 32},
	35: {30, 36, 45, 47, 31, 21},
	36: {29, 37, 44, 45, 35, 30},
	37: {38, 43, 44, 36, 29},
	38: {28, 39, 41, 43, 37, 29},
	39: {27, 40, 41, 38, 28},
	40: {39, 27, 26},
	41: {42, 43, 38, 39},
	42: {43, 41},
	43: {41, 42, 44, 37, 38},
	44: {38, 43, 46, 45, 36},
	45: {36, 44, 46, 48, 47, 35},
	46: {44, 48, 46},
	47: {35, 45, 48, 49, 31},
	48: {47, 45, 46, 49},
	49: {47, 48, 34, 31},
}

func main() {
	saveFile := flag.String("saveFile", ".\\SAVEDATA", "Save File")
	saveIndex := flag.Int("saveIndex", 1, "Save Index")

	flag.Parse()
	fmt.Println(*saveFile, *saveIndex)

	readSaveFile(*saveIndex, *saveFile)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Watcher Initialize Error")
		panic(err)
	}
	fmt.Println("watcher initialized...")

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("watcher event : ", event)

				if event.Has(fsnotify.Write) {
					fmt.Println("watcher modified file : ", event.Name)
					readSaveFile(*saveIndex, *saveFile)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("watcher error : ", err)
			}
		}
	}()

	err = watcher.Add(*saveFile)
	if err != nil {
		fmt.Println("watcher add error")

	}
	fmt.Println("watcher started...")

	<-make(chan struct{})
}

func readSaveFile(saveIndex int, saveFile string) {
	f, err := os.Open(saveFile)
	if err != nil {
		fmt.Println("file open error")
		panic(err)
	}
	defer f.Close()

	landInfo := newLandInfo()
	landInfo.Read(saveIndex, f)
	PrintHeader()
	fmt.Println(landInfo)

	turnInfo := NewTurnInfo()
	turnInfo.Read(saveIndex, f)
	//fmt.Println(turnInfo)
	//fmt.Printf("Current Land : [%d]\n", turnInfo.GetCurrentLand())

	heroInfo := newHeroInfo()
	heroInfo.Read(saveIndex, f)
	//fmt.Println(heroInfo)

	//landInfo.printNeighbor(turnInfo.GetCurrentLand())

	totalInfo := newTotalInfo(landInfo, *heroInfo)
	totalInfo.printNeighbor(turnInfo.GetCurrentLand())
}
