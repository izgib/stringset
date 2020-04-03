package stringset

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const NumItems = 1 << 15

var testExamples [NumItems]string

func init() {
	for i := 0; i < NumItems; i++ {
		testExamples[i] = "addr:" + strconv.FormatInt(rand.Int63n(NumItems*32), 16)
	}
}

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

func BenchmarkStringSetWithCollisionInfo(b *testing.B) {
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

func BenchmarkStringSet(b *testing.B) {
	stringSet := NewStringSet()
	testSet(b, stringSet)
}

func BenchmarkNewMapSet(b *testing.B) {
	mapSet := NewMapSet()
	testSet(b, mapSet)
}

func testSet(b *testing.B, set StringSet) {
	for i := 0; i < NumItems; i++ {
		set.Add(testExamples[i])
	}

	/*coinCount := 0
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
	b.Logf("coincidences: %d", coinCount)*/
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}
