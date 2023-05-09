package collection_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/collection"
)

func TestMain(m *testing.M) {
	logx.Disable()
	m.Run()
}

func BenchmarkRawSet(b *testing.B) {
	m := make(map[any]struct{})
	for i := 0; i < b.N; i++ {
		m[i] = struct{}{}
		_ = m[i]
	}
}

func BenchmarkSet(b *testing.B) {
	s := collection.NewSet[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
		_ = s.Contains(i)
	}
}

func TestAdd(t *testing.T) {
	// given
	set := collection.NewSet[int]()
	values := []int{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddInt(t *testing.T) {
	// given
	set := collection.NewSet[int]()
	values := []int{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	keys := set.Keys()
	sort.Ints(keys)
	assert.EqualValues(t, values, keys)
}

func TestAddInt64(t *testing.T) {
	// given
	set := collection.NewSet[int64]()
	values := []int64{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(int64(1)) && set.Contains(int64(2)) && set.Contains(int64(3)))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddUint(t *testing.T) {
	// given
	set := collection.NewSet[uint]()
	values := []uint{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(uint(1)) && set.Contains(uint(2)) && set.Contains(uint(3)))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddUint64(t *testing.T) {
	// given
	set := collection.NewSet[uint64]()
	values := []uint64{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(uint64(1)) && set.Contains(uint64(2)) && set.Contains(uint64(3)))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddStr(t *testing.T) {
	// given
	set := collection.NewSet[string]()
	values := []string{"1", "2", "3"}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains("1") && set.Contains("2") && set.Contains("3"))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestContainsWithoutElements(t *testing.T) {
	// given
	set := collection.NewSet[int]()

	// then
	assert.False(t, set.Contains(1))
}

func TestRemove(t *testing.T) {
	// given
	set := collection.NewSet[int]()
	set.Add([]int{1, 2, 3}...)

	// when
	set.Remove(2)

	// then
	assert.True(t, set.Contains(1) && !set.Contains(2) && set.Contains(3))
}

func TestCount(t *testing.T) {
	// given
	set := collection.NewSet[int]()
	set.Add([]int{1, 2, 3}...)

	// then
	assert.Equal(t, set.Count(), 3)
}
