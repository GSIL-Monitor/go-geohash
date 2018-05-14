package geohash

import (
	"testing"
	"time"
	"math/rand"
)

func TestGeoHash_EncodeGeoHash(t *testing.T) {
	lat, lng := 0.0, 0.0
	precision := int(10)
	var err error
	var value string

	// 正常功能
	lat = 39.928167
	lng = 116.389550
	precision = 4
	value, err = EncodeGeoHash(lat, lng, precision)
	if err != nil {
		t.Fatalf("EncodeGeoHash error %v", err)
	}
	if value != "wx4g" {
		t.Fatalf("EncodeGeoHash value %v not equal with wx4g", value)
	}

	// 边界条件
	lat = 91.0
	lng = 0.0
	precision = 1
	_, err = EncodeGeoHash(lat, lng, precision)
	if err == nil {
		t.Fatalf("para check error. lat is invalid")
	}

	// 边界条件
	lat = 23.4
	lng = 181.2
	precision = 1
	_, err = EncodeGeoHash(lat, lng, precision)
	if err == nil {
		t.Fatalf("para check error. lng is invalid")
	}

	// 边界条件
	lat = 23.4
	lng = 121.2
	precision = -1
	_, err = EncodeGeoHash(lat, lng, precision)
	if err == nil {
		t.Fatalf("para check error. precision is invalid")
	}

}

func TestGeoHash_DecodeGeoHash(t *testing.T) {
	// 正常功能
	geoHash := "ezs42"
	lat, lng := DecodeGeoHash(geoHash)

	if len(lat) != 4 || len(lng) != 4 {
		t.Fatalf("bad len of lat or lng")
	}

	latMin := lat[0]
	latMax := lat[1]
	latErr := lat[3]

	if latMin >=42.584 || latMin <= 42.582 {
		t.Fatalf("bad latMin %v", latMin)
	}

	if latMax >=42.628 || latMax <= 42.626 {
		t.Fatalf("bad latMax %v", latMax)
	}

	if latErr >=0.023 || latErr <= 0.021 {
		t.Fatalf("bad latErr %v", latErr)
	}

	lngMin := lng[0]
	lngMax := lng[1]
	lngErr := lng[3]

	if  lngMin >=-5.624 || lngMin <= -5.626 {
		t.Fatalf("bad lngMin %v", lngMin)
	}

	if lngMax >=-5.580 || lngMax <= -5.583 {
		t.Fatalf("bad lngMax %v", lngMax)
	}

	if lngErr >=0.0220 || lngErr <= 0.0219 {
		t.Fatalf("bad lngErr %v", lngErr)
	}

}

func TestGeoHash_calculateAdjacent(t *testing.T) {
	//正常功能测试
	geoHash := "wx4g"
	top := calculateAdjacent(geoHash, TOP)
	if top != "wx4u" {
		t.Fatalf("geoHash %v bad top %v", geoHash, top)
	}

	topLeft := calculateAdjacent(top, LEFT)
	if topLeft != "wx4s" {
		t.Fatalf("geoHash %v bad topLeft %v", geoHash, topLeft)
	}

	topRight := calculateAdjacent(top, RIGHT)
	if topRight != "wx5h" {
		t.Fatalf("geoHash %v bad topRight %v", geoHash, topRight)
	}

	left := calculateAdjacent(geoHash, LEFT)
	if left != "wx4e" {
		t.Fatalf("geoHash %v bad left %v", geoHash, left)
	}

	right := calculateAdjacent(geoHash, RIGHT)
	if right != "wx55" {
		t.Fatalf("geoHash %v bad right %v", geoHash, right)
	}

	bottom := calculateAdjacent(geoHash, BOTTOM)
	if bottom != "wx4f" {
		t.Fatalf("geoHash %v bad bottom %v", geoHash, bottom)
	}

	bottomLeft := calculateAdjacent(bottom, LEFT)
	if bottomLeft != "wx4d" {
		t.Fatalf("geoHash %v bad bottomLeft %v", geoHash, bottomLeft)
	}

	bottomRight := calculateAdjacent(bottom, RIGHT)
	if bottomRight != "wx54" {
		t.Fatalf("geoHash %v bad bottomRight %v", geoHash, bottomRight)
	}
}

func TestGetAdjacentGridGeoHash(t *testing.T) {
	geoHash := "zvxupc"
	resultSlice := GetAdjacentGridGeoHash(geoHash)
	if len(resultSlice) != 8 {
		t.Fatalf("len(resultSlice) %v != 8", len(resultSlice))
	}

	if resultSlice[2] != "" {
		t.Fatalf("right %v should be blank", resultSlice[2])
	}

	if resultSlice[1] != "" {
		t.Fatalf("top right %v should be blank", resultSlice[1])
	}

	if resultSlice[3] != "" {
		t.Fatalf("bottom right %v should be blank", resultSlice[3])
	}

}

func BenchmarkEncodeGeoHashPrecision12(b *testing.B) {
	unixNano := time.Now().UnixNano()
	s := rand.NewSource(unixNano)
	r := rand.New(s)
	precision := int(12)
	for i := 0; i < b.N; i++ {
		lat := r.Float64() * 180 - 90.0
		lng := r.Float64()* 360 - 180.0
		EncodeGeoHash(lat, lng, precision)
	}
}


func BenchmarkEncodeGeoHashPrecision6(b *testing.B) {
	unixNano := time.Now().UnixNano()
	s := rand.NewSource(unixNano)
	r := rand.New(s)
	precision := int(6)
	for i := 0; i < b.N; i++ {
		lat := r.Float64() * 180 - 90.0
		lng := r.Float64()* 360 - 180.0
		EncodeGeoHash(lat, lng, precision)
	}
}

func BenchmarkDecodeGeoHash(b *testing.B) {
	geoHash := "wtw37qtx"
	for i := 0; i < b.N; i++ {
		DecodeGeoHash(geoHash)
	}
}

func BenchmarkGetAdjacentGridGeoHash1(b *testing.B) {

	geoHash := "wx4g"
	for i := 0; i < b.N; i++ {
		GetAdjacentGridGeoHash(geoHash)
	}
}

func BenchmarkGetAdjacentGridGeoHash2(b *testing.B) {
	geoHash := "wtw37qtx"
	for i := 0; i < b.N; i++ {
		GetAdjacentGridGeoHash(geoHash)
	}
}
