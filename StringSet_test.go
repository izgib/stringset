package stringset

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringSet(t *testing.T) {
	stringSet := NewStringSet()
	str := "lol 12"
	wrongStr := "1445vig"
	assert.True(t, stringSet.Add(str))
	assert.True(t, stringSet.In(str))
	assert.True(t, !stringSet.Add(str))
	assert.True(t, stringSet.Delete(str))
	assert.True(t, !stringSet.In(str))

	assert.True(t, !stringSet.In(wrongStr))
	assert.True(t, !stringSet.Delete(wrongStr))
}

func Test_lul(t *testing.T) {
	slice := make([]*string, 0, 10)
	for i, str := range slice {
		fmt.Printf("%d: %s\n", i, *str)
	}
	str := "hello"
	slice = append(slice, &str)
	for i, str := range slice {
		fmt.Printf("%d: %s\n", i, *str)
	}
	println(slice)

	slice = removeEl(slice, 0)
	for i, str := range slice {
		fmt.Printf("%d: %s\n", i, *str)
	}
	print(slice)
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
