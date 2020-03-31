package stringset

import (
	"fmt"
	"hash/maphash"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringSet(t *testing.T) {
	stringSet := NewStringSet()
	testSetImpl(t, stringSet)

}

func Test_stringSet_resizeSet(t *testing.T) {

	stringSet := &set{hash: maphash.Hash{}, capacity: defaultCapacity, B: defaultCapacityBits}
	stringSet.resizeSet()
	for i := 0; i < defaultCapacity*loadFactorNum/loadFactorDen; i++ {
		assert.True(t, stringSet.Add(fmt.Sprint("val :", i)))
	}
	assert.True(t, cap(stringSet.buckets) == defaultCapacity<<1)
}

func Test_mapSet(t *testing.T) {
	mapSet := NewMapSet()
	testSetImpl(t, mapSet)
}

func testSetImpl(t *testing.T, set StringSet) {
	str := "lol 12"
	wrongStr := "1445vig"
	assert.True(t, set.Add(str))
	assert.True(t, set.In(str))
	assert.True(t, !set.Add(str))
	assert.True(t, set.Delete(str))
	assert.True(t, !set.In(str))

	assert.True(t, !set.In(wrongStr))
	assert.True(t, !set.Delete(wrongStr))
}

const NumItems = 1 << 15

func BenchmarkStringSet(b *testing.B) {
	b.StopTimer()
	stringSet := &set{hash: maphash.Hash{}, capacity: defaultCapacity, B: defaultCapacityBits}
	stringSet.resizeSet()

	count := NumItems
	coinCount := 0
	elBeforeResize := 0
	var statLogged bool
	for count > 0 {
		subtrahend := min(count, 1024)
		var buffer [1024]string
		var i int
		for i = 0; i < subtrahend; i++ {
			buffer[i] = "addr:" + strconv.FormatInt(rand.Int63n(NumItems*32), 16)
		}

		i = 0
		for {
			//b.Logf("count: %d, i: %d", count, i)
			if elBeforeResize == 0 {
				elBeforeResize = int(bucketShift(stringSet.B) * loadFactorNum / loadFactorDen)
				statLogged = false
			} else {
				if i != 0 {
					break
				}
			}
			b.StartTimer()
			for ; i < subtrahend && elBeforeResize > 0; i++ {
				if stringSet.Add(buffer[i]) {
					elBeforeResize--
				} else {
					coinCount++
				}
				if !statLogged && elBeforeResize == 1 {
					b.StopTimer()
					collisionInfo(b, int(bucketShift(stringSet.B)*loadFactorNum/loadFactorDen), stringSet.buckets)
					statLogged = true
					b.StartTimer()
				}
			}
		}
		b.StopTimer()
		count -= subtrahend
	}
	b.Logf("coincidences: %d", coinCount)
}

func collisionInfo(b *testing.B, count int, buckets [][]*string) {
	histMap := make(map[int]int, 2)
	used := 0
	for _, b := range buckets {
		count := len(b)
		if count == 0 {
			continue
		}
		histMap[count] = histMap[count] + 1
		used++
	}
	i := 0
	serEl := 1
	hist := strings.Builder{}
	printSep := false
	for serEl <= len(histMap) {
		col, ok := histMap[i]
		if ok {
			if printSep {
				hist.WriteString(", ")
			} else {
				printSep = true
			}
			hist.WriteString(fmt.Sprintf("%d: %.4f", i, float32(col)/float32(used)))
			serEl++
		}
		i++
	}
	b.Logf("strings: %d, %s", count, hist.String())
}

func BenchmarkNewMapSet(b *testing.B) {
	b.StopTimer()
	mapSet := NewMapSet()
	testSet(b, mapSet)
}

func testSet(b *testing.B, set StringSet) {
	count := NumItems
	coinCount := 0
	for count > 0 {
		subtrahend := min(count, 1024)
		var buffer [1024]string
		var i int
		for i = 0; i < subtrahend; i++ {
			buffer[i] = "addr:" + strconv.FormatInt(rand.Int63n(NumItems*32), 16)
		}
		b.StartTimer()
		for i = 0; i < subtrahend; i++ {
			if !set.Add(buffer[i]) {
				coinCount++
			}
		}
		b.StopTimer()
		count -= subtrahend
	}
	b.Logf("coincidences: %d", coinCount)
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

// max returns the maximum value in the input slice. If the slice is empty, max will panic.
func max(s []float64) float64 {
	return s[maxIdx(s)]
}

// maxIdx returns the index of the maximum value in the input slice. If several
// entries have the maximum value, the first such index is returned. If the slice
// is empty, maxIdx will panic.
func maxIdx(s []float64) int {
	if len(s) == 0 {
		panic("floats: zero slice length")
	}
	max := math.NaN()
	var ind int
	for i, v := range s {
		if math.IsNaN(v) {
			continue
		}
		if v > max || math.IsNaN(max) {
			max = v
			ind = i
		}
	}
	return ind
}
