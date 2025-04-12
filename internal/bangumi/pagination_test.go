package bangumi

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaginate(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	fetch := func(offset int) ([]int, int, error) {
		end := offset + defaultPaginationLimit
		if end > len(data) {
			end = len(data)
		}
		return data[offset:end], len(data), nil
	}

	result, err := paginate(fetch)

	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestPaginate_FetchErrorFirstPage(t *testing.T) {
	fetch := func(offset int) ([]int, int, error) {
		return nil, 0, errors.New("network error")
	}

	result, err := paginate(fetch)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "network error")
}
