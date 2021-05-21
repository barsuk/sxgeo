package sxgeo

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Detect and set here your host endianness manually if do not want to use DetectEndian
var hbo binary.ByteOrder

func init() {
	SetEndian(LITTLE)
}

var id2iso = [255]string{
	"",
	"AP",
	"EU",
	"AD",
	"AE",
	"AF",
	"AG",
	"AI",
	"AL",
	"AM",
	"CW",
	"AO",
	"AQ",
	"AR",
	"AS",
	"AT",
	"AU",
	"AW",
	"AZ",
	"BA",
	"BB",
	"BD",
	"BE",
	"BF",
	"BG",
	"BH",
	"BI",
	"BJ",
	"BM",
	"BN",
	"BO",
	"BR",
	"BS",
	"BT",
	"BV",
	"BW",
	"BY",
	"BZ",
	"CA",
	"CC",
	"CD",
	"CF",
	"CG",
	"CH",
	"CI",
	"CK",
	"CL",
	"CM",
	"CN",
	"CO",
	"CR",
	"CU",
	"CV",
	"CX",
	"CY",
	"CZ",
	"DE",
	"DJ",
	"DK",
	"DM",
	"DO",
	"DZ",
	"EC",
	"EE",
	"EG",
	"EH",
	"ER",
	"ES",
	"ET",
	"FI",
	"FJ",
	"FK",
	"FM",
	"FO",
	"FR",
	"SX",
	"GA",
	"GB",
	"GD",
	"GE",
	"GF",
	"GH",
	"GI",
	"GL",
	"GM",
	"GN",
	"GP",
	"GQ",
	"GR",
	"GS",
	"GT",
	"GU",
	"GW",
	"GY",
	"HK",
	"HM",
	"HN",
	"HR",
	"HT",
	"HU",
	"ID",
	"IE",
	"IL",
	"IN",
	"IO",
	"IQ",
	"IR",
	"IS",
	"IT",
	"JM",
	"JO",
	"JP",
	"KE",
	"KG",
	"KH",
	"KI",
	"KM",
	"KN",
	"KP",
	"KR",
	"KW",
	"KY",
	"KZ",
	"LA",
	"LB",
	"LC",
	"LI",
	"LK",
	"LR",
	"LS",
	"LT",
	"LU",
	"LV",
	"LY",
	"MA",
	"MC",
	"MD",
	"MG",
	"MH",
	"MK",
	"ML",
	"MM",
	"MN",
	"MO",
	"MP",
	"MQ",
	"MR",
	"MS",
	"MT",
	"MU",
	"MV",
	"MW",
	"MX",
	"MY",
	"MZ",
	"NA",
	"NC",
	"NE",
	"NF",
	"NG",
	"NI",
	"NL",
	"NO",
	"NP",
	"NR",
	"NU",
	"NZ",
	"OM",
	"PA",
	"PE",
	"PF",
	"PG",
	"PH",
	"PK",
	"PL",
	"PM",
	"PN",
	"PR",
	"PS",
	"PT",
	"PW",
	"PY",
	"QA",
	"RE",
	"RO",
	"RU",
	"RW",
	"SA",
	"SB",
	"SC",
	"SD",
	"SE",
	"SG",
	"SH",
	"SI",
	"SJ",
	"SK",
	"SL",
	"SM",
	"SN",
	"SO",
	"SR",
	"ST",
	"SV",
	"SY",
	"SZ",
	"TC",
	"TD",
	"TF",
	"TG",
	"TH",
	"TJ",
	"TK",
	"TM",
	"TN",
	"TO",
	"TL",
	"TR",
	"TT",
	"TV",
	"TW",
	"TZ",
	"UA",
	"UG",
	"UM",
	"US",
	"UY",
	"UZ",
	"VA",
	"VC",
	"VE",
	"VG",
	"VI",
	"VN",
	"VU",
	"WF",
	"WS",
	"YE",
	"YT",
	"RS",
	"ZA",
	"ZM",
	"ME",
	"ZW",
	"A1",
	"XK",
	"O1",
	"AX",
	"GG",
	"IM",
	"JE",
	"BL",
	"MF",
	"BQ",
	"SS",
}

type Info struct {
	Ver         byte   // C
	Time        uint32 // N
	Type        byte
	Charset     byte
	BIdxLen     byte
	MIdxLen     uint16 // n
	Range       uint16
	DbItems     uint32
	IdLen       byte
	MaxRegion   uint16
	MaxCity     uint16
	RegionSize  uint32
	CitySize    uint32
	MaxCountry  uint16
	CountrySize uint32
	PackSize    uint16
}

// String распечатает информацию о БД
func (i *Info) String() string {
	t := time.Unix(int64(i.Time), 0)
	info := fmt.Sprintf("Version %v, Time %s, Type %v, Charset %v", i.Ver, t, i.Type, i.Charset)
	return info
}

type Meta struct {
	BlockLen byte
	BIdxStr  []byte
	MIdxStr  []byte

	Pack        [][]byte
	PackHeaders [][]Header

	DbBegin      int64
	BIdxArr      []uint32
	MIdxArr      [][]byte
	RegionsBegin int64
	CitiesBegin  int64
}

var (
	DB      []byte
	Regions []byte
	Cities  []byte
)

type Full struct {
	City    *City    `json:"city"`
	Country *Country `json:"country"`
	Region  *Region  `json:"region"`
}

type City struct {
	ID         int     `json:"id"`
	NameRu     string  `json:"name_ru"`
	NameEn     string  `json:"name_en"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
	RegionSeek float64 `json:"region_seek"`
}

type Country struct {
	ID     uint8  `json:"id"`
	ISO    string `json:"iso"`
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
}

type Region struct {
	ID     int    `json:"id"`
	ISO    string `json:"iso"`
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
}

var I Info
var M Meta

// ToMemory Reads the whole DB to the memory without temp file
func ToMemory(f *bytes.Reader) (bool, error) {
	h := make([]byte, 3)
	_, err := f.Read(h)
	if err != nil {
		return false, fmt.Errorf("cannot read header: %v", err)
	}
	if err == io.EOF {
		return false, fmt.Errorf("достигнут конец файла. Так быстро? %v", err)
	}

	if string(h) != "SxG" {
		return false, fmt.Errorf("this is not a SxGeo database! %v", err)
	}

	if err := binary.Read(f, binary.BigEndian, &I); err != nil {
		return false, fmt.Errorf("cannot unpack info: %v", err)
	}

	M.BlockLen = 3 + I.IdLen

	packRaw := make([]byte, I.PackSize)
	_, err = f.Read(packRaw)
	if err != nil {
		return false, fmt.Errorf("cannot read pack: %v", err)
	}
	packHeaders := bytes.Split(packRaw, []byte("\000"))
	hs := make([][]Header, 3)
	for i := 0; i < len(packHeaders); i++ {
		ph, err := unpackHeaders(packHeaders[i])
		if err != nil {
			return false, fmt.Errorf("cannot unpack headers: %v", err)
		}
		hs[i] = ph
	}
	//os.Exit(0)
	//M.Pack = packHeaders
	M.PackHeaders = hs

	bIdxStr := make([]byte, int(I.BIdxLen)*4)
	_, err = f.Read(bIdxStr)
	if err != nil {
		return false, fmt.Errorf("cannot read b-index: %v", err)
	}
	M.BIdxStr = bIdxStr

	// m_idx_str
	mIdxStr := make([]byte, I.MIdxLen*4)
	_, err = f.Read(mIdxStr)
	if err != nil {
		return false, fmt.Errorf("cannot read m-index: %v", err)
	}
	M.MIdxStr = mIdxStr

	bIdxArr := make([]uint32, len(M.BIdxStr)/4)
	if err := binary.Read(bytes.NewReader(M.BIdxStr), binary.BigEndian, &bIdxArr); err != nil {
		return false, fmt.Errorf("cannot unpack b-idx-array: %v", err)
	}
	M.BIdxArr = bIdxArr

	var mIdxArr [][]byte
	r := bytes.NewReader(mIdxStr)
	for i := 0; i < len(mIdxStr); {
		word := make([]byte, 4)
		_, err := r.Read(word)
		if err != nil {
			return false, fmt.Errorf("ты не угадал с числом байтов: %v", err)
		}
		mIdxArr = append(mIdxArr, word)
		i += 4
	}
	M.MIdxArr = mIdxArr

	dbBegin, err := f.Seek(0, 1)
	if err != nil {
		return false, fmt.Errorf("cannot seek to offset: %v", err)
	}
	M.DbBegin = dbBegin

	db := make([]byte, int(I.DbItems)*int(M.BlockLen))
	_, err = f.Read(db)
	if err != nil {
		return false, fmt.Errorf("cannot read db to the memory: %v", err)
	}
	DB = db

	regions := make([]byte, int(I.RegionSize))
	_, err = f.Read(regions)
	if err != nil {
		return false, fmt.Errorf("cannot read regions to the memory: %v", err)
	}
	Regions = regions

	cities := make([]byte, int(I.CitySize))
	_, err = f.Read(cities)
	if err != nil {
		return false, fmt.Errorf("cannot read cities to the memory: %v", err)
	}
	Cities = cities

	M.RegionsBegin = M.DbBegin + int64(I.DbItems)*int64(M.BlockLen)
	M.CitiesBegin = M.RegionsBegin + int64(I.RegionSize)

	return true, nil
}


// ReadDBToMemory Reads the whole DB to the memory
func ReadDBToMemory(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("cannot open file: %v", err)
	}
	defer f.Close()

	h := make([]byte, 3)
	_, err = f.Read(h)
	if err != nil {
		return false, fmt.Errorf("cannot read header: %v", err)
	}
	if err == io.EOF {
		return false, fmt.Errorf("достигнут конец файла. Так быстро? %v", err)
	}

	if string(h) != "SxG" {
		return false, fmt.Errorf("this is not a SxGeo database! %v", err)
	}

	if err := binary.Read(f, binary.BigEndian, &I); err != nil {
		return false, fmt.Errorf("cannot unpack info: %v", err)
	}

	M.BlockLen = 3 + I.IdLen

	packRaw := make([]byte, I.PackSize)
	_, err = f.Read(packRaw)
	if err != nil {
		return false, fmt.Errorf("cannot read pack: %v", err)
	}
	packHeaders := bytes.Split(packRaw, []byte("\000"))
	hs := make([][]Header, 3)
	for i := 0; i < len(packHeaders); i++ {
		ph, err := unpackHeaders(packHeaders[i])
		if err != nil {
			return false, fmt.Errorf("cannot unpack headers: %v", err)
		}
		hs[i] = ph
	}
	//os.Exit(0)
	//M.Pack = packHeaders
	M.PackHeaders = hs

	bIdxStr := make([]byte, int(I.BIdxLen)*4)
	_, err = f.Read(bIdxStr)
	if err != nil {
		return false, fmt.Errorf("cannot read b-index: %v", err)
	}
	M.BIdxStr = bIdxStr

	// m_idx_str
	mIdxStr := make([]byte, I.MIdxLen*4)
	_, err = f.Read(mIdxStr)
	if err != nil {
		return false, fmt.Errorf("cannot read m-index: %v", err)
	}
	M.MIdxStr = mIdxStr

	bIdxArr := make([]uint32, len(M.BIdxStr)/4)
	if err := binary.Read(bytes.NewReader(M.BIdxStr), binary.BigEndian, &bIdxArr); err != nil {
		return false, fmt.Errorf("cannot unpack b-idx-array: %v", err)
	}
	M.BIdxArr = bIdxArr

	var mIdxArr [][]byte
	r := bytes.NewReader(mIdxStr)
	for i := 0; i < len(mIdxStr); {
		word := make([]byte, 4)
		_, err := r.Read(word)
		if err != nil {
			return false, fmt.Errorf("ты не угадал с числом байтов: %v", err)
		}
		mIdxArr = append(mIdxArr, word)
		i += 4
	}
	M.MIdxArr = mIdxArr

	dbBegin, err := f.Seek(0, 1)
	if err != nil {
		return false, fmt.Errorf("cannot seek to offset: %v", err)
	}
	M.DbBegin = dbBegin

	db := make([]byte, int(I.DbItems)*int(M.BlockLen))
	_, err = f.Read(db)
	if err != nil {
		return false, fmt.Errorf("cannot read db to the memory: %v", err)
	}
	DB = db

	regions := make([]byte, int(I.RegionSize))
	_, err = f.Read(regions)
	if err != nil {
		return false, fmt.Errorf("cannot read regions to the memory: %v", err)
	}
	Regions = regions

	cities := make([]byte, int(I.CitySize))
	_, err = f.Read(cities)
	if err != nil {
		return false, fmt.Errorf("cannot read cities to the memory: %v", err)
	}
	Cities = cities

	M.RegionsBegin = M.DbBegin + int64(I.DbItems)*int64(M.BlockLen)
	M.CitiesBegin = M.RegionsBegin + int64(I.RegionSize)

	return true, nil
}

// Header распакованный заголовок города, региона или страны
type Header struct {
	Name     string
	Selector string
	Length   int
	PowExp   float64
}

// распаковывает заранее заголовки города, региона или страны
func unpackHeaders(p []byte) ([]Header, error) {
	var headers []Header
	packHeaders := bytes.Split(p, []byte("/"))
	for _, ph := range packHeaders {
		x := bytes.Split(ph, []byte(":"))
		typeChars, name := x[0], string(x[1])
		selector := string(typeChars[0]) // $type0
		h := Header{
			Name:     name,
			Selector: selector,
		}

		if selector == "N" || selector == "n" || selector == "c" {
			length, err := strconv.Atoi(string(typeChars[1]))
			if err != nil {
				return nil, fmt.Errorf("cannot parse type char for n: %v", err)
			}
			powexp := math.Pow(10, float64(length))
			h.Length = length
			h.PowExp = powexp
		}
		headers = append(headers, h)
	}

	return headers, nil
}

func GetCityFull(ip string) (*Full, error) {
	seek, err := Seek(ip)
	if err != nil {
		return nil, err
	}

	if seek < 1 {
		return nil, fmt.Errorf("unknown error with seek")
	}

	parsedCity, err := parseFullCity(seek)
	if err != nil {
		return nil, fmt.Errorf("error with full city")
	}

	return parsedCity, nil
}

const (
	PACKTYPECOUNTRY = 0
	PACKTYPEREGION  = 1
	PACKTYPECITY    = 2
)

func parseFullCity(seek uint32) (*Full, error) {
	full := new(Full)
	var regionSeek uint32
	if seek < I.CountrySize {
		country, err := readData(seek, I.MaxCountry, PACKTYPECOUNTRY)
		if err != nil {
			return nil, fmt.Errorf("cannot read country data")
		}

		full.City = &City{}
		full.Region = &Region{}
		full.Country = &Country{
			ID:     country["id"].(uint8),
			ISO:    fmt.Sprintf("%s", country["iso"]),
			NameRu: fmt.Sprintf("%s", country["name_ru"]),
			NameEn: fmt.Sprintf("%s", country["name_en"]),
		}
	} else {
		city, err := readData(seek, I.MaxCity, PACKTYPECITY)
		if err != nil {
			return nil, fmt.Errorf("cannot read country data")
		}
		regionSeek = city["region_seek"].(uint32)

		//full.Country = &Country{
		//	ID:  city["country_id"].(uint8),
		//	ISO: id2iso[city["country_id"].(uint8)],
		//}

		full.City = &City{
			ID:        int((city["id"]).(uint32)),
			NameRu:    fmt.Sprintf("%s", city["name_ru"]),
			NameEn:    fmt.Sprintf("%s", city["name_en"]),
			Latitude:  city["lat"].(float64),
			Longitude: city["lon"].(float64),
		}

		region, err := readData(regionSeek, I.MaxRegion, PACKTYPEREGION)
		if err != nil {
			return nil, fmt.Errorf("cannot read region data")
		}

		full.Region = &Region{
			ID:     int((region["id"]).(uint32)),
			ISO:    fmt.Sprintf("%s", region["iso"]),
			NameRu: fmt.Sprintf("%s", region["name_ru"]),
			NameEn: fmt.Sprintf("%s", region["name_en"]),
		}

		country, err := readData(uint32(region["country_seek"].(uint16)), I.MaxCountry, PACKTYPECOUNTRY)
		if err != nil {
			return nil, fmt.Errorf("failed to read country")
		}
		full.Country = &Country{
			ID:     country["id"].(uint8),
			ISO:    fmt.Sprintf("%s", country["iso"]),
			NameRu: fmt.Sprintf("%s", country["name_ru"]),
			NameEn: fmt.Sprintf("%s", country["name_en"]),
		}
	}

	return full, nil
}

func readData(seek uint32, max uint16, packType int) (map[string]interface{}, error) {
	var raw []byte
	if (seek > 0) && (max > 0) {
		if packType == 1 {
			raw = Regions[seek : seek+uint32(max)]
		} else {
			raw = Cities[seek : seek+uint32(max)]
		}
	}

	// странная заглушка на случай, если вдруг
	if len(raw) < 1 {
		unpacked := make(map[string]interface{})
		for _, ph := range M.PackHeaders[packType] {
			if ph.Selector == "b" || ph.Selector == "c" {
				unpacked[ph.Name] = nil
			} else {
				unpacked[ph.Name] = nil
			}
		}
		return unpacked, nil
	}

	unpacked, err := unpack(packType, raw)
	if err != nil {
		return nil, fmt.Errorf("cannot unpack city or region")
	}

	return unpacked, nil
}

func unpack(packType int, item []byte) (map[string]interface{}, error) {
	unpacked := make(map[string]interface{})
	position := 0
	for _, ph := range M.PackHeaders[packType] {
		name := ph.Name
		selector := ph.Selector // $type0

		var length int
		switch selector {
		case "t", "T":
			length = 1
		case "s", "n", "S":
			length = 2
		case "m", "M":
			length = 3
		case "d":
			length = 8
		case "c":
			length = ph.Length
			//panic("c!")
		case "b":
			length = bytes.Index(item[position:], []byte{0})
			if length < 0 {
				panic("здесь должон быть нулевой байт!")
			}
			//panic("b!")
		default:
			length = 4
			//panic("default!")
		}

		val := item[position : position+length]
		var result interface{}
		switch selector {
		case "t":
			var v int8
			if err := binary.Read(bytes.NewReader(val), binary.BigEndian, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			panic("t")
		case "T":
			var v byte
			if err := binary.Read(bytes.NewReader(val), binary.BigEndian, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			//panic("T")
		case "s":
			var v int16
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			panic("s")
		case "S":
			var v uint16
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			//panic("S")
		case "m":
			//unpack('l', $val . (ord($val{2}) >> 7 ? "\xff" : "\0"));
			var v rune
			var sel byte
			if val[2]>>7 != 0 {
				sel = 0xFF
			} else {
				sel = 0x00
			}

			buf := bytes.NewBuffer([]byte{})
			buf.Write(val)
			buf.Write([]byte{sel})
			if err := binary.Read(buf, hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			panic("m")
		case "M":
			var v uint32
			buf := bytes.NewBuffer([]byte{})
			buf.Write(val)
			buf.Write([]byte{0})
			if err := binary.Read(strings.NewReader(buf.String()), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			//panic("M")
		case "i":
			var v rune
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			panic("i")
		case "I":
			var v uint32
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			panic("I")
		case "f":
			//f	float (машинно-зависимые размер и представление)
			var v float32
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			panic("f")
		case "d":
			//d	double (машинно-зависимые размер и представление)
			var v float64
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = v
			panic("d")
		case "n":
			var v int16
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = float64(v) / ph.PowExp
			//panic("n")
		case "N":
			var v rune
			if err := binary.Read(bytes.NewReader(val), hbo, &v); err != nil {
				return nil, fmt.Errorf("cannot unpack selected: %v", err)
			}
			result = float64(v) / ph.PowExp
			//panic("N")
		case "c":
			v := bytes.TrimRight(val, " ")
			result = v
			//panic("c")
		case "b":
			result = val
			length++
			//panic("b")
		}
		position += length
		unpacked[name] = result
		//fmt.Printf("%[1]s, %[1]T\n", result)
	}

	return unpacked, nil
}

// Seek seeks an IP index in the DB
// It's a reflect of PHP get_num($ip)
func Seek(ip string) (uint32, error) {
	IP := net.ParseIP(ip)
	if IP == nil {
		return 0, fmt.Errorf("wrong IP format")
	}
	if IP.IsLoopback() || IP.IsMulticast() || IP.IsUnspecified() {
		return 0, fmt.Errorf("IP is loopback or multicast or unspecified")
	}

	ip4 := IP.To4()
	if ip4 == nil {
		return 0, fmt.Errorf("IP format is not IPv4")
	}
	if ip4[0] > I.BIdxLen-1 {
		return 0, fmt.Errorf("IP is out of DB diapason")
	}

	ipN := []byte(ip4) // $ipn = pack('N', $ipn);

	blocksMin := M.BIdxArr[ip4[0]-1]
	blocksMax := M.BIdxArr[ip4[0]]

	var min, max uint32
	if blocksMax-blocksMin > uint32(I.Range) {
		// Ищем блок в основном индексе
		start := blocksMin / uint32(I.Range)
		end := blocksMax/uint32(I.Range) - 1
		part := searchIdx(ipN, start, end)

		// Нашли номер блока, в котором нужно искать IP, теперь находим нужный блок в БД
		if part > 0 {
			min = part * uint32(I.Range)
		}
		if part > uint32(I.MIdxLen) {
			max = I.DbItems
		} else {
			max = (part + 1) * uint32(I.Range)
		}

		// Нужно проверить чтобы блок не выходил за пределы блока первого байта
		if min < blocksMin {
			min = blocksMin
		}

		if max > blocksMax {
			max = blocksMax
		}
	} else {
		min = blocksMin
		max = blocksMax
	}

	// Находим нужный диапазон в БД
	dec, err := searchDb(ipN, min, max)
	if err != nil {
		return 0, fmt.Errorf("cannot find in DB")
	}

	return dec, nil
}

func searchDb(ipN []byte, min uint32, max uint32) (uint32, error) {
	if max-min > 1 {
		bcd := ipN[1:]
		for max-min > 8 {
			offset := (min + max) >> 1
			start := int(offset) * int(M.BlockLen)
			end := start + 3
			dbSubs := DB[start:end]
			if string(bcd) > string(dbSubs) {
				min = offset
			} else {
				max = offset
			}
		}

		//while ( $ipn >= substr($str, $min * $this->block_len, 3) && ++$min < $max ) { };
		for string(bcd) >= string(DB[int(min)*int(M.BlockLen):int(min)*int(M.BlockLen)+3]) &&
			func(p *uint32) uint32 {
				*p += 1
				return *p
			}(&min) < max {
		}
	} else {
		min++
	}

	start := int(min)*int(M.BlockLen) - int(I.IdLen)

	bin := DB[start : start+int(I.IdLen)]

	hx := make([]byte, hex.EncodedLen(len(bin)))
	_ = hex.Encode(hx, bin)

	s, err := strconv.ParseInt(fmt.Sprintf("%s", hx), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert hex to dec")
	}

	return uint32(s), nil
}

func searchIdx(ipN []byte, min uint32, max uint32) uint32 {
	mx := max
	mn := min
	for (mx - mn) > 8 {
		offset := (mn + mx) >> 1
		if string(ipN) > string(M.MIdxArr[offset]) {
			mn = offset
		} else {
			mx = offset
		}
	}

	for string(ipN) > string(M.MIdxArr[mn]) &&
		func(p *uint32) uint32 {
			old := *p
			*p += 1
			return old
		}(&mn) < mx {
	}

	return mn
}
