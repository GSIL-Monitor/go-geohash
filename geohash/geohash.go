//  GeoHash的encode/decode方法，以及高效查询附近8个格子GeoHash的方法
//  实现参考了 https://en.wikipedia.org/wiki/Geohash 以及 https://github.com/davetroy/geohash-js
//  根据GeoHash原理，GeoHash每一个字符，都对应32个格子中的某一个固定的格子。

//  如果字符在奇数(ODD)位，对应的划分是 经度、纬度、经度、纬度、经度。那么每个字符分布的位置如下:
//  +---------------+------+
//  |b  c  f  g  u  v  y  z|
//  |8  9  d  e  s  t  w  x|
//  |2  3  6  7  k  m  q  r|
//  |0  1  4  5  h  j  n  p|
//  +----------------------+

//  类似的，如果字符在偶数(EVEN)位, 对应的划分是 纬度、经度、纬度、经度、纬度
//  +------------+
//  | p  r  x  z |
//  | n  q  w  y |
//  | j  m  t  v |
//  | h  k  s  u |
//  | 5  7  e  g |
//  | 4  6  d  f |
//  | 1  3  9  c |
//  | 0  2  8  b |
//  +------------+

//  查询一个格子附近的8个格子的GeoHash字符串，实际上就是查询该格子对应GeoHash的最后一个字符在
//  上面两表之一中(根据最后一个字符是奇数位还是偶数位),其周围的8个字符即可。
//  如果恰好字符在边界处，那么需要回退到原GeoHash的上一个字符，查询该字符同方向上的下一个字符...

package geohash

import (
	"fmt"
	"bytes"
	"git.xiaojukeji.com/soda-server/recsys-broker/utils"
	"git.xiaojukeji.com/soda-framework/go-log"
	"strings"
)

type BaseDirection int
const (
	RIGHT BaseDirection = 0
	LEFT BaseDirection = 1
	TOP BaseDirection = 2
	BOTTOM BaseDirection = 3
)

const (
	EVEN = 0
	ODD = 1
)

var (
	Base32 = []byte{'0','1','2','3','4','5','6','7','8','9','b','c','d','e','f','g',
	'h','j','k','m','n','p','q','r','s','t','u','v','w','x','y','z'}
	Bits = []int{16, 8, 4, 2, 1}
	CharBitLength = len(Bits)
	NEIGHBORS = make([][][]byte, 2)
	BORDERS = make([][][]byte, 2)
)

func init() {
	// 存储Base32中每个字符的相邻的上下左右字符。这里实际存储的是每个字符相反方向的字符
	NEIGHBORS[EVEN] =  make([][]byte, 4)
	NEIGHBORS[EVEN][RIGHT] =  []byte{'b','c','0','1','f','g','4','5','2','3','8','9',
	'6','7','d','e','u','v','h','j','y','z','n','p','k','m','s','t','q','r','w','x'}
	NEIGHBORS[EVEN][LEFT] =   []byte{'2','3','8','9','6','7','d','e','b','c','0','1',
	'f','g','4','5','k','m','s','t','q','r','w','x','u','v','h','j','y','z','n','p'}
	NEIGHBORS[EVEN][TOP] =    []byte{'p','0','r','2','1','4','3','6','x','8','z','b',
	'9','d','c','f','5','h','7','k','j','n','m','q','e','s','g','u','t','w','v','y'}
	NEIGHBORS[EVEN][BOTTOM] = []byte{'1','4','3','6','5','h','7','k','9','d','c','f',
	'e','s','g','u','j','n','m','q','p','0','r','2','t','w','v','y','x','8','z','b'}

	NEIGHBORS[ODD] =  make([][]byte, 4)
	NEIGHBORS[ODD][RIGHT] = NEIGHBORS[EVEN][TOP]
	NEIGHBORS[ODD][LEFT] = NEIGHBORS[EVEN][BOTTOM]
	NEIGHBORS[ODD][TOP] = NEIGHBORS[EVEN][RIGHT]
	NEIGHBORS[ODD][BOTTOM] = NEIGHBORS[EVEN][LEFT]

	BORDERS[EVEN] = make([][]byte, 4)
	BORDERS[EVEN][RIGHT] = []byte{'b','c','f','g','u','v','y','z'}
	BORDERS[EVEN][LEFT] = []byte{'0','1','4','5','h','j','n','p'}
	BORDERS[EVEN][TOP] = []byte{'p','r','x','z'}
	BORDERS[EVEN][BOTTOM] = []byte{'0','2','8','b'}

	BORDERS[ODD] = make([][]byte, 4)
	BORDERS[ODD][RIGHT] = BORDERS[EVEN][TOP]
	BORDERS[ODD][LEFT] = BORDERS[EVEN][BOTTOM]
	BORDERS[ODD][TOP] = BORDERS[EVEN][RIGHT]
	BORDERS[ODD][BOTTOM] = BORDERS[EVEN][LEFT]
}

// GeoHash的编码函数，输入 longitude（经度）, latitude（纬度）以及精度，返回一个字符串
func EncodeGeoHash(latitude, longitude float64, precision int) (string, error) {
	if precision <=0 || ! utils.CoordinateCheck(latitude, longitude){
		log.Errorf("encodeGeoHash|input para lat %v lng %v precision %v invalid",
			latitude, longitude, precision)
		return "", fmt.Errorf("input para lat %v lng %v precision %v invalid", latitude, longitude, precision)
	}

	var IndexEven = true
	var mid float64
	latStart, latEnd := -90.0, 90.0
	lngStart, lngEnd := -180.0, 180.0
	var buffer bytes.Buffer
	buffer.Grow(precision)

	var ch = 0
	var bitIndex = 0
	var length = 0

	for length < precision {
		if IndexEven {
			mid = (lngStart + lngEnd) / 2
			if longitude > mid {
				ch |= Bits[bitIndex]
				lngStart = mid
			} else {
				lngEnd = mid
			}
		} else {
			mid =  (latStart + latEnd) / 2
			if latitude > mid {
				ch |= Bits[bitIndex]
				latStart = mid
			} else {
				latEnd = mid
			}
		}
		IndexEven = !IndexEven
		if bitIndex < 4 {
			bitIndex++
		} else {
			buffer.WriteByte(Base32[ch])
			bitIndex = 0
			ch = 0
			length++
		}
	}

	return buffer.String(), nil
}

// GeoHash的解码函数，给定一个GeoHash字符串，返回latitude 纬度，longitude经度信息
func DecodeGeoHash(geoHash string) (lat1, lng1 []float64) {
	var IndexEven = true
	var lat, lng []float64
	lat = append(lat, -90.0)
	lat = append(lat, 90.0)

	lng = append(lng, -180.0)
	lng = append(lng, 180.0)

	latErr, lngErr := 90.0, 180.0

	length := len(geoHash)
	for i:=0;i<length;i++ {
		ch := geoHash[i]
		idx := byteSliceIndexOf(Base32, ch)
		for j:=0; j < CharBitLength; j++ {
			mask := Bits[j]
			if IndexEven {
				lngErr /= 2
				refineInterval(lng, idx, mask)
			} else {
				latErr /= 2
				refineInterval(lat, idx, mask)
			}
			IndexEven = !IndexEven
		}
	}
	lat = append(lat, (lat[0] + lat[1]) / 2)
	lat = append(lat, latErr)
	lng= append(lng, (lng[0] + lng[1]) / 2)
	lng = append(lng, lngErr)
	return lat, lng
}

// 二分法，缩小格子区间
func refineInterval(interval []float64, idx, mask int) {

	if (idx & mask) > 0 {
		interval[0] = (interval[0] + interval[1]) / 2
	} else {
		interval[1] = (interval[0] + interval[1]) / 2
	}
}


func byteSliceIndexOf(byteSlice []byte, k byte) int {
	for idx, v := range byteSlice {
		if v == k {
			return idx
		}
	}

	return -1
}

func calculateAdjacent(geoHash string, direction BaseDirection) string {
	length := len(geoHash)
	if length == 0  {
		return ""
	}
	var IndexType = EVEN
	if length % 2 == 1 {
		IndexType = ODD
	}

	lowerHash := strings.ToLower(geoHash)
	lastChr := lowerHash[length - 1]

	var base = lowerHash[:length-1]
	var lastChrOnBorder = false
	if  byteSliceIndexOf(BORDERS[IndexType][direction], lastChr) != -1 {
		base = calculateAdjacent(base, direction)
		lastChrOnBorder = true
	}

	if lastChrOnBorder && base == "" {
		return ""
	}

	idx := byteSliceIndexOf(NEIGHBORS[IndexType][direction], lastChr)
	if idx == -1 {
		log.Errorf("calculateAdjacent bad idx. geoHash %v direction %v IndexType %v lastChr %v",
			geoHash, direction, IndexType, lastChr)
		return ""
	}
	return base + string(Base32[idx])
}

// 获取GeoHash附近的8个格子的GeoHash，从最上面开始顺时针顺序
func GetAdjacentGridGeoHash(geoHash string) ([]string) {
	leftGrid := calculateAdjacent(geoHash, LEFT)
	rightGrid := calculateAdjacent(geoHash, RIGHT)
	topGrid := calculateAdjacent(geoHash, TOP)
	bottomGrid := calculateAdjacent(geoHash, BOTTOM)

	topLeftGrid := calculateAdjacent(topGrid, LEFT)
	topRightGrid := calculateAdjacent(topGrid, RIGHT)
	bottomLeftGrid := calculateAdjacent(bottomGrid, LEFT)
	bottomRightGrid := calculateAdjacent(bottomGrid, RIGHT)

	return []string{topGrid, topRightGrid, rightGrid, bottomRightGrid,
	bottomGrid, bottomLeftGrid, leftGrid, topLeftGrid}
}
