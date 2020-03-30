package stringset

import (
	"fmt"
	"hash/maphash"
	"math/rand"
	"strconv"
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
	stringSet := NewStringSet()
	testSet(b, stringSet)
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
