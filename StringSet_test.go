package stringset

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringSet(t *testing.T) {
	stringSet := NewStringSet()
	testSetImpl(t, stringSet)
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
	for count > 0 {
		subtrahend := min(count, 1024)
		var buffer [1024]string
		var i int
		for i = 0; i < subtrahend; i++ {
			buffer[i] = "addr:" + strconv.FormatInt(rand.Int63n(NumItems*2), 16)
		}
		b.StartTimer()
		for i = 0; i < subtrahend; i++ {
			set.Add(buffer[i])
		}
		b.StopTimer()
		count -= subtrahend
	}
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}
