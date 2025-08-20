package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock using Testify's mock framework (not utilized in tests but provided as per instructions)
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(key string, value any) {
	m.Called(key, value)
}

func (m *MockCache) Get(key string) (any, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockCache) Del(key string) any {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockCache) RPush(key string, data []any) int {
	args := m.Called(key, data)
	return args.Int(0)
}

func (m *MockCache) LRange(key string, start, end int) []any {
	args := m.Called(key, start, end)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]any)
}

func TestSetAndGetStoresAndRetrievesValue(t *testing.T) {
	c := New()
	c.Set("foo", 123)

	val, ok := c.Get("foo")
	require.True(t, ok)
	assert.Equal(t, 123, val)
}

func TestRPushAppendsAndLRangeReturnsSubslice(t *testing.T) {
	c := New()

	n := c.RPush("list", []any{1, 2, 3})
	require.Equal(t, 3, n)

	n = c.RPush("list", []any{4, 5})
	require.Equal(t, 5, n)

	sub := c.LRange("list", 1, 3) // returns indices [1:3) => 2,3
	assert.Equal(t, []any{2, 3}, sub)
}

func TestDelReturnsOldValueAndRemovesKey(t *testing.T) {
	c := New()
	c.Set("a", "x")

	old := c.Del("a")
	assert.Equal(t, "x", old)

	_, ok := c.Get("a")
	assert.False(t, ok)
}

func TestLRangeReturnsEmptyArrayWhenStartGreaterThanEnd(t *testing.T) {
	c := New()
	c.RPush("list", []any{1, 2, 3})

	got := c.LRange("list", 2, 1)
	assert.Equal(t, []any{}, got)
}

func TestLRangeClampsNegativeStartAndOversizedEnd(t *testing.T) {
	c := New()
	c.RPush("list", []any{10, 20, 30, 40})

	got := c.LRange("list", -5, 999) // clamps to v[0:len-1] => v[0:3] => 10,20,30
	require.NotNil(t, got)
	assert.Equal(t, []any{10, 20, 30}, got)
}

func TestLRangeReturnsEmptyArrayForEmptyList(t *testing.T) {
	c := New()

	got := c.LRange("nonexistent", 0, 5)
	assert.Equal(t, []any{}, got)
}

func TestLRangeReturnsEntireListForStartAndEndEqualToListLength(t *testing.T) {
	c := New()
	c.RPush("list", []any{"strawberry", "apple", "blueberry", "grape", "orange"})
	got := c.LRange("list", 0, 3)
	assert.Equal(t, []any{"strawberry", "apple", "blueberry", "grape"}, got)
	got = c.LRange("list", 3, 4)
	assert.Equal(t, []any{"grape", "orange"}, got)
}

func TestGetReturnsFalseForMissingKey(t *testing.T) {
	c := New()

	val, ok := c.Get("missing")
	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestSetOverwritesExistingKeyReturnsUpdatedValue(t *testing.T) {
	c := New()

	c.Set("k", 1)
	c.Set("k", 2)

	val, ok := c.Get("k")
	require.True(t, ok)
	assert.Equal(t, 2, val)
}

func TestDelMissingKeyReturnsNilNoPanic(t *testing.T) {
	c := New()
	c.Set("other", "val")

	var ret any
	assert.NotPanics(t, func() {
		ret = c.Del("missing")
	})
	assert.Nil(t, ret)

	val, ok := c.Get("other")
	require.True(t, ok)
	assert.Equal(t, "val", val)
}

func TestLRangeNegativeIndex(t *testing.T) {
	c := New()
	c.RPush("list", []any{"strawberry", "blueberry", "mango", "apple", "orange", "pineapple", "pear"})
	got := c.LRange("list", -5, -1)
	assert.Equal(t, []any{"mango", "apple", "orange", "pineapple", "pear"}, got)
}
