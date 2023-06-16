package main

import (
	"bytes"
	"errors"
	"fmt"
)

const HERO_CONST_VALUE = 46

type TotalInfo struct {
	landTotalMap map[int]LandTotal
}

func newTotalInfo(li LandInfo, hi HeroInfo) TotalInfo {
	var ti TotalInfo
	ti.landTotalMap = make(map[int]LandTotal)

	heroMap := make(map[int][]Hero)
	peopleMap := make(map[int][]Hero)
	for _, hero := range hi.HeroMap {
		if hero.Property.isHero() {
			heroMap[int(hero.Property.Location)] = append(heroMap[int(hero.Property.Location)], hero)
		} else if hero.Property.isPeople() {
			peopleMap[int(hero.Property.Location)] = append(peopleMap[int(hero.Property.Location)], hero)
		} else {
			//fmt.Printf("Not Hero, People : [%s]\n", hero)
		}
	}

	for index, landProperty := range li.PropertyMap {
		rulerIdx, err := getRulerIdx(landProperty.Ruler)

		var ruler, owner Hero

		if err == nil {
			ruler = hi.HeroMap[rulerIdx]
			owner = hi.HeroMap[int(ruler.Property.Ruled+1)]
		}

		landMen := 0
		for _, hero := range heroMap[index] {
			landMen += int(hero.Property.Men)
		}

		ti.landTotalMap[index] = LandTotal{LandNo: index, Land: li.PropertyMap[index], Men: landMen,
			Ruler: ruler, Owner: owner, Heros: heroMap[index], People: peopleMap[index]}
	}

	return ti
}

func (ti TotalInfo) printNeighbor(landNo int) {
	fmt.Printf("Current Land : %d\n", landNo)

	fmt.Println(getLandTotalHeader())
	for _, neighborLandNo := range NeighborMap[landNo] {
		fmt.Println(ti.landTotalMap[neighborLandNo])
		//ti.landTotalMap[neighborLandNo].printHero()
		//ti.landTotalMap[neighborLandNo].printPeople()
	}
}

type LandTotal struct {
	LandNo int
	Land   LandProperty
	Ruler  Hero
	Owner  Hero
	Men    int
	Heros  []Hero
	People []Hero
}

func (lt LandTotal) String() string {
	var heroBuffer bytes.Buffer
	for _, hero := range lt.Heros {
		heroBuffer.WriteString(fmt.Sprintf("%s ", hero.Name))
	}

	var peopleBuffer bytes.Buffer
	for _, hero := range lt.People {
		peopleBuffer.WriteString(fmt.Sprintf("%s ", hero.Name))
	}

	return fmt.Sprintf("[%2d][%3s][%3s][%4d,%4d,%4d,%3d,%3d,%2d,%2d,%2d]\nh:[%s] p:[%s]",
		lt.LandNo, lt.Owner.Name, lt.Ruler.Name, lt.Land.Money, lt.Land.Food, lt.Men,
		lt.Land.Arms, lt.Land.Skill, lt.Land.Land, lt.Land.Wealth, lt.Land.Support, heroBuffer.String(), peopleBuffer.String())
}

func (lt LandTotal) printHero() {
	fmt.Println("Heros...")
	for _, hero := range lt.Heros {
		fmt.Println(hero)
	}
}

func (lt LandTotal) printPeople() {
	fmt.Println("People...")
	for _, hero := range lt.People {
		fmt.Println(hero)
	}
}

func getLandTotalHeader() string {
	return fmt.Sprintf("[NO][%6s][%6s][MONEYFOOD  MEN,ARM,SKI,LA,WE,SU]", "OWNER", "RULER")
}

func getRulerIdx(rulerValue uint16) (int, error) {
	if int(rulerValue) > 0 {
		return (int(rulerValue) - HERO_CONST_VALUE) / HERO_LENGTH, nil
	} else {
		return 0, errors.New("No Ruler")
	}
}
