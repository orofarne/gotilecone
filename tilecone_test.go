package tilecone

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCanUseDB(t *testing.T) {
	dir := "__test_db_1"
	index := `[Tiles]
BucketZoom = 15
TileZoom = 19
BlockSize = 1024
  `
	os.Mkdir(dir, os.ModeDir|os.ModePerm)
	ioutil.WriteFile(
		fmt.Sprintf("%s%cindex.ini", dir, os.PathSeparator),
		[]byte(index), 0644)
	defer os.RemoveAll(dir)

	db, err := NewDB(dir, 10)
	if err != nil {
		t.Fatalf("NewDB error: %v", err)
	}
	if db == nil {
		t.Fatalf("db is nil")
	}

	defer db.Free()

	testData := "My test data"

	err = db.SetTile(316893, 163547, []byte(testData))
	if err != nil {
		t.Fatalf("SetTile error: %v", err)
	}

	resData, err := db.GetTiles(19, 316893, 163547)
	if err != nil {
		t.Fatalf("GetTile error: %v", err)
	}
	if len(resData) != 1 {
		t.Fatalf("len(resData) = %v (should be 1)", len(resData))
	}
	if string(resData[0]) != testData {
		t.Fatalf("testData != resData (\"%v\" != \"%v\")", testData, string(resData[0]))
	}
}

func TestCanReadDBInfo(t *testing.T) {
	dir := "__test_db_1"
	index := `[Tiles]
BucketZoom = 15
TileZoom = 19
BlockSize = 1024
  `
	os.Mkdir(dir, os.ModeDir|os.ModePerm)
	ioutil.WriteFile(
		fmt.Sprintf("%s%cindex.ini", dir, os.PathSeparator),
		[]byte(index), 0644)
	defer os.RemoveAll(dir)

	db, err := NewDB(dir, 10)
	if err != nil {
		t.Fatalf("NewDB error: %v", err)
	}
	if db == nil {
		t.Fatalf("db is nil")
	}

	defer db.Free()

	if db.BucketZoom() != 15 {
		t.Errorf("Invalid bucket zoom %v (should be %v)", db.BucketZoom(), 15)
	}
	if db.TileZoom() != 19 {
		t.Errorf("Invalid tile zoom %v (should be %v)", db.TileZoom(), 19)
	}
}
