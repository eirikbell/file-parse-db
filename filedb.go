package fileparsedb

import (
	"strconv"
	"strings"
	"sort"
	"fmt"
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"io/ioutil"
)

// func main() {

// 	db_path := "db/parsed_files.db"

// 	fdb := NewFileDb(db_path, "data/conviva", strings.Compare)

// 	file, err := fdb.GetNextFile()

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if len(file) <= 0 {
// 		log.Info("NONE LEFT")
// 		return
// 	}

// 	log.Info("NEXT", file)

// 	err = fdb.MarkParsed(file)
	
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

type FileDB struct {
	dbPath string
	fileBucket string
	fileDir string
	fileNameComparer func(string, string) int
}

func NewFileDb(dbPath, fileDir string, fileNameComparer func(string, string) int) *FileDB {
	fdb := FileDB{dbPath: dbPath, fileBucket: "syncFiles", fileDir: fileDir, fileNameComparer: fileNameComparer}

	db, err := fdb.openDb()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(fdb.fileBucket))
		return err
	})

	return &fdb
}

func (fdb *FileDB) GetNextFile() (string, error) {
	fileList, err := fdb.listFiles()
	if err != nil {
		return "", err
	}

	err = fdb.ensureFilesInDb(fileList)
	if err != nil {
		return "", err
	}

	filesToProcess, err := fdb.getNonProcessedFiles()
	if err != nil {
		return "", err
	}

	if len(filesToProcess) <= 0 {
		return "", nil
	}

	return fdb.sortFiles(filesToProcess)[0], nil
}

func (fdb *FileDB) MarkParsed(file string) error {
	return fdb.storeFile(file, true)
}

func (fdb *FileDB) getNonProcessedFiles() ([]string, error) {
	db, err := fdb.openDb()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	result := []string{}

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(fdb.fileBucket))
	
		c := b.Cursor()
	
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			strVal := string(v[:])
			value, err := strconv.ParseBool(strVal)
			if err != nil {
				return err
			}

			if !value {
				result = append(result, string(k[:]))
			}
		}
	
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, err
}

func (fdb *FileDB) ensureFilesInDb(fileList []string) error {
	if fileList == nil || len(fileList) <= 0 {
		//error condition
	}

	for _, file := range fileList {
		exists, err := fdb.fileExists(file)
		if err != nil {
			return err
		}

		if !exists {
			fdb.storeFile(file, false)
		}
	}

	return nil
}

func (fdb *FileDB) storeFile(file string, value bool) error {
	db, err := fdb.openDb()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	strVal := strconv.FormatBool(value)

	err = db.Update(func(tx *bolt.Tx) error {
		fmt.Println("store file", file, value)
		b := tx.Bucket([]byte(fdb.fileBucket))
		err := b.Put([]byte(file), []byte(strVal))
		return err
	})

	return err
}

func (fdb *FileDB) fileExists(file string) (bool, error) {
	db, err := fdb.openDb()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer db.Close()

	result := false
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fdb.fileBucket))
		v := b.Get([]byte(file))
		fmt.Printf("The answer for %s is: %s\n", file, v)

		if v != nil {
			result = true
		}
		return nil
	})

	return result, err
}

func (fdb *FileDB) listFiles() ([]string, error) {
	dirContent, err := ioutil.ReadDir(fdb.fileDir)
	if err != nil {
		return nil, err
	}

	fileList := []string{}
	for _, file := range dirContent {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	return fileList, nil
}

func (fdb *FileDB) sortFiles(files []string) []string {
	comparer := fdb.fileNameComparer

	if comparer == nil {
		comparer = strings.Compare
	}

	sort.Slice(files, func(i, j int) bool { return comparer(files[i], files[j]) < 0})

	return files
}

func (fdb *FileDB) openDb() (*bolt.DB, error) {
	db, err := bolt.Open(fdb.dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return db, nil
}