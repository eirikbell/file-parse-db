package fileparsedb

import (
	"reflect"
	"strconv"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFileDB(t *testing.T) {
	assert := assert.New(t)

	t.Run("Test values stored on object", func(t *testing.T){
		dbPath := "db/test.db"
		fileDir := "data/conviva"
		fileDb := NewFileDb(dbPath, fileDir, sortAsInt)
		assert.Equal(dbPath, fileDb.dbPath)
		assert.Equal(fileDir, fileDb.fileDir)
		assert.Equal(reflect.ValueOf(sortAsInt), reflect.ValueOf(fileDb.fileNameComparer))
	})
}

func TestSortFiles(t *testing.T) {
	assert := assert.New(t)

	files := []string{"1", "2", "11", "9", "20", "3", "test", "1pa"}
	fdb := NewFileDb("db/test.db", "data/test", strings.Compare)

	t.Run("Default sort by name", func(t *testing.T){
		fdb.fileNameComparer = nil
		res := fdb.sortFiles(files)
		fmt.Println(res)
		assert.Equal("1" , res[0])
		assert.Equal("11" , res[1])
		assert.Equal("1pa" , res[2])
		assert.Equal("2" , res[3])
		assert.Equal("20" , res[4])
		assert.Equal("3" , res[5])
		assert.Equal("9" , res[6])
		assert.Equal("test" , res[7])
	})

	t.Run("Sort by name", func(t *testing.T){
		fdb.fileNameComparer = strings.Compare
		res := fdb.sortFiles(files)
		fmt.Println(res)
		assert.Equal("1" , res[0])
		assert.Equal("11" , res[1])
		assert.Equal("1pa" , res[2])
		assert.Equal("2" , res[3])
		assert.Equal("20" , res[4])
		assert.Equal("3" , res[5])
		assert.Equal("9" , res[6])
		assert.Equal("test" , res[7])
	})

	t.Run("Sort by name descending", func(t *testing.T){
		fdb.fileNameComparer = func(s1, s2 string) int { return 0 - strings.Compare(s1, s2) }
		res := fdb.sortFiles(files)
		fmt.Println(res)
		assert.Equal("1" , res[7])
		assert.Equal("11" , res[6])
		assert.Equal("1pa" , res[5])
		assert.Equal("2" , res[4])
		assert.Equal("20" , res[3])
		assert.Equal("3" , res[2])
		assert.Equal("9" , res[1])
		assert.Equal("test" , res[0])
	})

	t.Run("Sort as int", func(t *testing.T){
		fdb.fileNameComparer = sortAsInt
		res := fdb.sortFiles(files)
		fmt.Println(res)
		assert.Equal("1" , res[0])
		assert.Equal("2" , res[1])
		assert.Equal("3" , res[2])
		assert.Equal("9" , res[3])
		assert.Equal("11" , res[4])
		assert.Equal("20" , res[5])	
		assert.Equal("1pa" , res[6])	
		assert.Equal("test" , res[7])
	})
}

func sortAsInt(s1, s2 string) int {
	if v1, err := strconv.ParseInt(s1, 10, 32); err == nil {
		if v2, err := strconv.ParseInt(s2, 10, 32); err == nil {
			return int(v1 - v2)
		}
		return -1
	}
	if _, err := strconv.ParseInt(s2, 10, 32); err == nil {
		return 1
	}
	return strings.Compare(s1, s2)
}