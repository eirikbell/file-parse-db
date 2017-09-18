# Description
Wrapper library for BoltDB to synchroneously process files in a given order

# Usage
```
// db will be created if it doesnt exist from before
fdb := NewFileDb("/some/dir/dbname.db", "/folder/containing/files", "", strings.Compare)

// Gets the next non processed file
fileName, err := fdb.GetNextFile()
if err != nil {
	log.Fatal(err)
}

// If fileName is empty, no file is remaining for processing
if len(fileName) <= 0 {
	log.Info("NONE LEFT")
	return
}

// Do something with file
log.Info("NEXT", fileName)

err = fdb.MarkParsed(fileName)
if err != nil {
	log.Fatal(err)
}

```

# Require file prefix
```
fdb := NewFileDb("/some/dir/dbname.db", "/folder/containing/files", "fileNamePrefix", strings.Compare)
...

```


## Custom string comparer
```
func sortAsInt(s1, s2 string) int {
	if v1, err := strconv.ParseInt(s1, 10, 32); err == nil {
		if v2, err := strconv.ParseInt(s2, 10, 32); err == nil {
            // if both strings are int
			return int(v1 - v2)
		}
        // if only first string is int, sort this first
		return -1
	}
	if _, err := strconv.ParseInt(s2, 10, 32); err == nil {
        // if only second string is int, sort this first
		return 1
	}
    // if neither strings are ints, fallback to lexicographic order
	return strings.Compare(s1, s2)
}

fdb := NewFileDb("/some/dir/dbname.db", "/folder/containing/files", "", sortAsInt)

```