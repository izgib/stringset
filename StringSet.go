package stringset

import (
	"hash/maphash"
)

const (
	defaultCapacityBits = 4
	defaultCapacity     = 1 << defaultCapacityBits
	resizeThreshold     = 0.75
)

type StringSet interface {
	Add(string) bool
	In(string) bool
	Delete(string) bool
}

func NewStringSet() StringSet {
	set := &stringSet{hash: maphash.Hash{}, capacity: defaultCapacity, B: defaultCapacityBits}
	set.resizeSet()
	return set
}

type stringSet struct {
	hash     maphash.Hash
	capacity uint32
	count    uint32
	buckets  [][]*string
	B        uint8
}

func (s *stringSet) Add(str string) bool {
	hash := s.getHash(str)
	m := hash & bucketMask(s.B)
	bucket := s.buckets[m]
	for _, el := range bucket {
		if str == *el {
			return false
		}
	}
	if float32((s.count+1)/s.capacity) >= resizeThreshold {
		s.B++
		s.resizeSet()
		m = hash & bucketMask(s.B)
	}
	s.buckets[m] = append(s.buckets[m], &str)
	s.count++
	return true
}

func (s *stringSet) resizeSet() {
	s.capacity = bucketShift(s.B)
	buckets := make([][]*string, s.capacity)
	for i := uint32(0); i < s.capacity; i++ {
		buckets[i] = make([]*string, 0, 1)
	}
	for _, bucket := range s.buckets {
		for _, str := range bucket {
			m := s.getHash(*str) & bucketMask(s.B)
			buckets[m] = append(buckets[m], str)
		}
	}
	s.buckets = buckets
}

func (s *stringSet) getHash(str string) uint32 {
	s.hash.Reset()
	s.hash.WriteString(str)
	return uint32(s.hash.Sum64())
}

func (s *stringSet) In(str string) bool {
	m := s.getHash(str) & bucketMask(s.B)
	for _, el := range s.buckets[m] {
		if str == *el {
			return true
		}
	}
	return false
}

func (s *stringSet) Delete(str string) bool {
	m := s.getHash(str) & bucketMask(s.B)
	bucket := s.buckets[m]
	for i, el := range bucket {
		if str == *el {
			s.buckets[i] = removeEl(bucket, i)
			s.count--
			if float32(s.count/bucketShift(s.B-1)) < resizeThreshold {
				s.B--
				s.resizeSet()
			}
			return true
		}
	}
	return false
}

// bucketShift returns 1<<b, optimized for code generation.
func bucketShift(b uint8) uint32 {
	// Masking the shift amount allows overflow checks to be elided.
	return uint32(1) << (b & (32 - 1))
}

// bucketMask returns 1<<b - 1, optimized for code generation.
func bucketMask(b uint8) uint32 {
	return bucketShift(b) - 1
}

func removeEl(slice []*string, ind int) []*string {
	for i := ind + 1; i < len(slice); i++ {
		slice[i-1] = slice[i]
	}
	return slice[:len(slice)-1]
}
