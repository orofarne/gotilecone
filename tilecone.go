package tilecone

// #cgo LDFLAGS: -ltilecone
// #include <stdlib.h>
// #include <tilecone/db.h>
//
// void *data_at_offset(void *data, size_t offset) {
//   return (char *)data + offset;
// }
//
// struct tc_tile *tile_at_pos(struct tc_tile *data, size_t pos) {
//   return data + pos;
// }
//
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

type DB struct {
	db C.tc_db
}

func NewDB(path string, mmappool int) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	db := &DB{}
	db.db = C.tc_new_db(cpath, C.int(mmappool))
	if 0 == C.tc_db_ok(db.db) {
		errCStr := C.tc_last_error(db.db)
		return nil, errors.New(C.GoString(errCStr))
	}
	runtime.SetFinalizer(db, func(obj interface{}) { obj.(*DB).Free() })
	return db, nil
}

func (db *DB) Free() {
	C.tc_free_db(db.db)
}

func (db *DB) SetTile(x uint64, y uint64, data []byte) error {
	rc := C.tc_set_tile(db.db, C.uint64_t(x), C.uint64_t(y), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	if 0 != rc {
		errCStr := C.tc_last_error(db.db)
		return errors.New(C.GoString(errCStr))
	}
	return nil
}

func (db *DB) GetTiles(zoom uint16, x uint64, y uint64) (tiles [][]byte, err error) {
	var pData unsafe.Pointer
	var pTiles *C.struct_tc_tile
	var lTiles C.size_t

	rc := C.tc_get_tiles(db.db, C.uint16_t(zoom), C.uint64_t(x), C.uint64_t(y), &pData, &pTiles, &lTiles)
	if 0 != rc {
		errCStr := C.tc_last_error(db.db)
		err = errors.New(C.GoString(errCStr))
		return
	}

	for i := uint(0); i < uint(lTiles); i++ {
		tilePointer := C.tile_at_pos(pTiles, C.size_t(i))
		dataPointer := C.data_at_offset(pData, C.size_t(tilePointer.offset))
		data := C.GoBytes(dataPointer, C.int(tilePointer.size))
		tiles = append(tiles, data)
	}

	return
}

func (db *DB) BucketZoom() uint64 {
	return uint64(C.tc_bucket_zoom(db.db))
}

func (db *DB) TileZoom() uint64 {
	return uint64(C.tc_tile_zoom(db.db))
}
