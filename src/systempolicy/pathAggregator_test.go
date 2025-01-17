package systempolicy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ===================== //
// == Build Path Tree == //
// ===================== //

func TestAggregatePaths_1(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/UserDict.py", "/usr/lib/python2.7/UserDict.pyo"}

	results := AggregatePaths(paths)

	assert.Equal(t, len(results), 2)
	assert.False(t, results[0].isDir)
	assert.False(t, results[1].isDir)
}

func TestAggregatePaths_2(t *testing.T) {
	paths := []string{
		"/usr/lib/python2.7/UserDict.py",
		"/usr/lib/python2.7/UserDict.pyo",
		"/usr/lib/python2.7/UserDict.3",
		"/usr/lib/python2.7/UserDict.4",
	}

	results := AggregatePaths(paths)

	assert.Equal(t, len(results), 1)
	assert.True(t, results[0].isDir)
}

func TestAggregatePaths_3(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/", "/usr/lib/python2.8/", "/usr/lib/python2.9/", "/usr/lib/python2.10/"}

	results := AggregatePathsExt(paths)

	assert.Equal(t, len(results), 4)
	for _, str := range paths {
		assert.Contains(t, results, str)
	}
}

func TestAggregatePaths_4(t *testing.T) {
	paths := []string{
		"/usr/lib/python2.7/",
		"/usr/lib/python2.8/",
		"/usr/xyz/python2.9/",
		"/usr/lib/python2.10/",
		"/usr/lib/python2.11/",
	}

	results := AggregatePathsExt(paths)

	assert.Equal(t, len(results), 5)
	for _, str := range paths {
		assert.Contains(t, results, str)
	}
}

func TestAggregatePaths_5(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/", "/usr/lib/python2.7/xyz", "/usr/lib/python2.7/folder/xyz"}

	results := AggregatePathsExt(paths)

	assert.Equal(t, len(results), 2)
	assert.Contains(t, results, "/usr/lib/python2.7/")
	assert.Contains(t, results, "/usr/lib/python2.7/folder/xyz")
}

func TestMergeFileInDir_1(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/", "/usr/lib/python2.7/xyz"}

	dlist, flist := mergeFileInDir(paths)

	assert.Equal(t, len(dlist), 1)
	assert.Equal(t, len(flist), 0)
	assert.Contains(t, dlist, "/usr/lib/python2.7/")
}

func TestMergeFileInDir_2(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/", "/usr/lib/python2.7/xyz", "/usr/lib/python2.7/folder/xyz"}

	dlist, flist := mergeFileInDir(paths)

	assert.Equal(t, len(dlist), 1)
	assert.Equal(t, len(flist), 1)
	assert.Equal(t, flist[0], "/usr/lib/python2.7/folder/xyz")
	assert.Contains(t, dlist, "/usr/lib/python2.7/")
}

func TestMergeFileInDir_3(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/", "/usr/lib/python2.7/xyz",
		"/usr/lib/python2.7/folder/xyz", "/usr/lib/python2.7/folder/"}

	dlist, flist := mergeFileInDir(paths)

	assert.Equal(t, len(dlist), 2)
	assert.Equal(t, len(flist), 0)
	assert.Contains(t, dlist, "/usr/lib/python2.7/")
	assert.Contains(t, dlist, "/usr/lib/python2.7/folder/")
}

func TestAggregatePathsExt_1(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/", "/usr/lib/python2.7/xyz",
		"/usr/lib/python2.7/folder/xyz", "/usr/lib/python2.7/folder/"}

	results := AggregatePathsExt(paths)

	assert.Equal(t, len(results), 2)
	assert.Contains(t, results, "/usr/lib/python2.7/")
	assert.Contains(t, results, "/usr/lib/python2.7/folder/")
}

func TestAggregatePathsExt_2(t *testing.T) {
	paths := []string{"/usr/lib/python2.7/", "/usr/lib/python2.7/xyz",
		"/usr/lib/python2.7/folder/xyz", "/usr/lib/python2.7/folder/",
		"/usr/lib/python2.7/folder/abc", "/usr/lib/python2.7/folder/lmn"}

	results := AggregatePathsExt(paths)

	assert.Equal(t, len(results), 2)
	assert.Equal(t, results[0], "/usr/lib/python2.7/")
	assert.Equal(t, results[1], "/usr/lib/python2.7/folder/")
}

func TestAggregatePathsExt_3(t *testing.T) {
	paths := []string{
		"/usr/l11/l21/",
		"/usr/l12/l22/",
		"/usr/l13/l23/1", "/usr/l13/l23/2", "/usr/l13/l23/3", "/usr/l13/l23/4",
		"/usr/l14/l24/1", "/usr/l14/l24/2", "/usr/l14/l24/3", "/usr/l14/l24/4",
	}

	results := AggregatePathsExt(paths)

	assert.Equal(t, len(results), 4)
	assert.Contains(t, results, "/usr/l11/l21/")
	assert.Contains(t, results, "/usr/l12/l22/")
	assert.Contains(t, results, "/usr/l13/l23/")
	assert.Contains(t, results, "/usr/l14/l24/")
}

func TestAggregatePathsExt_4(t *testing.T) {
	paths := []string{
		"/usr/l11/l21/",
		"/usr/l11/l21/",
		"/usr/l11/l21/xyz",
		"/usr/l11/l21/xyz",
		"/usr/l11/l21/abc",
	}

	results := AggregatePathsExt(paths)

	assert.Equal(t, len(results), 1)
}
