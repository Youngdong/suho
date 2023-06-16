package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const HERO_OFFSET = 0x02 // 2
const HERO_LENGTH = 22
const HERO_MAX_INDEX = 255

type HeroInfo struct {
	HeroMap map[int]Hero
}

func newHeroInfo() *HeroInfo {
	return &HeroInfo{HeroMap: make(map[int]Hero)}
}

func (h *HeroInfo) String() string {
	var b bytes.Buffer
	for k, v := range h.HeroMap {
		b.WriteString(fmt.Sprintf("%3d [%s]\n", k, v))
	}

	return b.String()
}

func (h *HeroInfo) Read(saveIndex int, f *os.File) {

	heroOffset := HEADER + ((saveIndex - 1) * SAVE_OFFSET) + HERO_OFFSET
	fmt.Printf("Hero Offset : %x, %d\n", heroOffset, heroOffset)

	seekOffset, err := f.Seek(int64(heroOffset), 0)
	if err != nil {
		fmt.Println("file seek error")
		panic(err)
	}
	fmt.Printf("offset : %d\n", seekOffset)

	for i := 1; i <= HERO_MAX_INDEX; i++ {
		byteArray := make([]byte, HERO_LENGTH)
		readSize, err := f.Read(byteArray)
		if len(byteArray) != readSize || err != nil {
			fmt.Println("file read error")
			panic(err)
		}

		h.HeroMap[i] = newHero(i, byteArray)
	}
}

type Hero struct {
	Name     string
	Property HeroProperty
}

func newHero(index int, byteArray []byte) Hero {
	return Hero{Name: HeroName[index], Property: newHeroProperty(byteArray)}
}

func (h *Hero) String() string {
	return fmt.Sprintf("%6s [%s]", h.Name, h.Property)
}

type HeroProperty struct {
	Pre1 uint8
	Pre2 uint8

	Age      uint8
	Ruled    uint8
	Location uint8
	Body     uint8
	BodyMax  uint8

	Integrity uint8
	Mercy     uint8
	Courage   uint8
	Str       uint8
	Dex       uint8

	Wis     uint8
	StrExp  uint8
	DexExp  uint8
	WisExp  uint8
	Royalty uint8

	Etc2    uint8
	Popular uint16
	Men     uint8
	Bitmap  uint8
}

func (h HeroProperty) String() string {
	return fmt.Sprintf("%2x, %2x, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %3d, %4d, %3d, %2x",
		h.Pre1, h.Pre2,
		h.Age, h.Ruled, h.Location, h.Body, h.BodyMax, h.Str, h.Dex, h.Wis,
		h.Integrity, h.Mercy, h.Courage, h.StrExp, h.DexExp, h.WisExp, h.Royalty,
		h.Etc2, h.Popular, h.Men, h.Bitmap)
}

func (h HeroProperty) isHero() bool {
	return h.Bitmap&0xA0 == 0x80
}

func (h HeroProperty) isPeople() bool {
	return h.Bitmap&0xA0 == 0xA0
}

func newHeroProperty(byteArray []byte) HeroProperty {
	var heroProperty HeroProperty
	buf := bytes.NewReader(byteArray)
	err := binary.Read(buf, binary.LittleEndian, &heroProperty)

	if err != nil {
		fmt.Printf("Hero Serialize Error %s\n", byteArray)
	}
	heroProperty.Location = heroProperty.Location + 1
	return heroProperty
}

var HeroName = map[int]string{
	1:   "고구",
	2:   "노지심",
	3:   "사진",
	4:   "송강",
	5:   "임충",
	6:   "무송",
	7:   "양지",
	8:   "이규",
	9:   "이응",
	10:  "이준",
	11:  "조개",
	12:  "주무",
	13:  "연순",
	14:  "배선",
	15:  "추연",
	16:  "번서",
	17:  "왕윤",
	18:  "최도성",
	19:  "증도",
	20:  "축용",
	21:  "이충",
	22:  "화영",
	23:  "관승",
	24:  "진명",
	25:  "왕진",
	26:  "경영",
	27:  "허관충",
	28:  "소가수",
	29:  "변상",
	30:  "난정옥",
	31:  "장청",
	32:  "호성",
	33:  "석보",
	34:  "자진",
	35:  "오용",
	36:  "노준의",
	37:  "연청",
	38:  "숙원경",
	39:  "동평",
	40:  "등원각",
	41:  "방걸",
	42:  "호삼낭",
	43:  "고염",
	44:  "손입",
	45:  "엽청",
	46:  "색초",
	47:  "축표",
	48:  "주동",
	49:  "황신",
	50:  "나진인",
	51:  "장숙야",
	52:  "초정",
	53:  "공손승",
	54:  "교열",
	55:  "전호",
	56:  "위정국",
	57:  "저형",
	58:  "선정규",
	59:  "사문공",
	60:  "구악",
	61:  "방경",
	62:  "여사낭",
	63:  "손안",
	64:  "이사사",
	65:  "뇌횡",
	66:  "방납",
	67:  "이운",
	68:  "유문충",
	69:  "장개",
	70:  "장횡",
	71:  "왕환",
	72:  "서영",
	73:  "호현",
	74:  "호정작",
	75:  "축호",
	76:  "체득장",
	77:  "양웅",
	78:  "호준",
	79:  "두미",
	80:  "황안",
	81:  "동관",
	82:  "산사기",
	83:  "주근",
	84:  "방만춘",
	85:  "상청",
	86:  "백흠",
	87:  "석수",
	88:  "한도",
	89:  "대종",
	90:  "원낭",
	91:  "저능",
	92:  "곽성",
	93:  "여방",
	94:  "중양",
	95:  "설찬",
	96:  "문환장",
	97:  "왕영",
	98:  "증괴",
	99:  "장순",
	100: "한존보",
	101: "문달",
	102: "선찬",
	103: "증승",
	104: "주귀",
	105: "비보",
	106: "추윤",
	107: "정득손",
	108: "고대수",
	109: "소정",
	110: "양전",
	111: "주통",
	112: "유당",
	113: "목홍",
	114: "설영",
	115: "체경",
	116: "황문엽",
	117: "동증",
	118: "가여경",
	119: "구붕",
	120: "범전",
	121: "동위",
	122: "장청",
	123: "임원",
	124: "시문빈",
	125: "양춘",
	126: "증삭",
	127: "원소칠",
	128: "방기",
	129: "마인",
	130: "공양",
	131: "손신",
	132: "원흥",
	133: "공왕",
	134: "미생",
	135: "전진붕",
	136: "오응성",
	137: "전표",
	138: "구정",
	139: "두흥",
	140: "오치",
	141: "하현통",
	142: "왕경",
	143: "당빈",
	144: "사정",
	145: "장경",
	146: "왕인",
	147: "능진",
	148: "동맹",
	149: "전표",
	150: "후건",
	151: "계삼은",
	152: "목춘",
	153: "등비",
	154: "송청",
	155: "당세영",
	156: "공빈",
	157: "우방희",
	158: "안도전",
	159: "유몽룡",
	160: "정지서",
	161: "장상",
	162: "침수",
	163: "공명",
	164: "악화",
	165: "서방",
	166: "원소이",
	167: "정천수",
	168: "증밀",
	169: "오복",
	170: "욱보사",
	171: "체택",
	172: "진달",
	173: "조정",
	174: "시은",
	175: "항원진",
	176: "양임",
	177: "원소오",
	178: "학사문",
	179: "하길",
	180: "마영",
	181: "탕융",
	182: "석용",
	183: "황포단",
	184: "포욱",
	185: "도종왕",
	186: "당세웅",
	187: "성본",
	188: "소양",
	189: "해진",
	190: "이곤",
	191: "유민",
	192: "구혈",
	193: "맹강",
	194: "곽세광",
	195: "장위",
	196: "항윤",
	197: "우옥린",
	198: "어득원",
	199: "김대견",
	200: "방모",
	201: "두천",
	202: "주부",
	203: "해보",
	204: "오이",
	205: "단경주",
	206: "왕의",
	207: "엽춘",
	208: "경공",
	209: "하도",
	210: "무대",
	211: "성희",
	212: "정승조",
	213: "체경",
	214: "이입",
	215: "포문영",
	216: "배여해",
	217: "손이낭",
	218: "시천",
	219: "체복",
	220: "포도을",
	221: "냉공",
	222: "장세개",
	223: "운종무",
	224: "송만",
	225: "장몽방",
	226: "단삼낭",
	227: "안사영",
	228: "염세걸",
	229: "방학도",
	230: "장충",
	231: "온극양",
	232: "왕정륙",
	233: "이달",
	234: "이회",
	235: "장문원",
	236: "문중용",
	237: "위충",
	238: "방천정",
	239: "여심",
	240: "이길",
	241: "서문경",
	242: "구소을",
	243: "마만리",
	244: "하청",
	245: "백승",
	246: "황문병",
	247: "장보",
	248: "등용",
	249: "이귀",
	250: "유고",
	251: "모중의",
	252: "반금련",
	253: "이서란",
	254: "백수영",
	255: "반공운",
}
