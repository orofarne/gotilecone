package tilecone

// #cgo LDFLAGS: -ltilecone
// #include <stdlib.h>
// #include <tilecone/db.h>
//
// void *data_at_offset(void *data, size_t offset) {
//   return (char *)data + offset;
// }
//
// struct tile *tile_at_pos(struct tile *data, size_t pos) {
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
	db C.db
}

func NewDB(path string, mmappool int) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	db := &DB{}
	db.db = C.new_db(cpath, C.int(mmappool))
	if 0 == C.db_ok(db.db) {
		errCStr := C.last_error(db.db)
		return nil, errors.New(C.GoString(errCStr))
	}
	runtime.SetFinalizer(db, func(obj interface{}) { obj.(*DB).Free() })
	return db, nil
}

func (db *DB) Free() {
	C.free_db(db.db)
}

func (db *DB) SetTile(x uint64, y uint64, data []byte) error {
	rc := C.set_tile(db.db, C.uint64_t(x), C.uint64_t(y), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	if 0 != rc {
		errCStr := C.last_error(db.db)
		return errors.New(C.GoString(errCStr))
	}
	return nil
}

func (db *DB) GetTiles(zoom uint16, x uint64, y uint64) (tiles [][]byte, err error) {
	var pData unsafe.Pointer
	var pTiles *C.struct_tile
	var lTiles C.size_t

	rc := C.get_tiles(db.db, C.uint16_t(zoom), C.uint64_t(x), C.uint64_t(y), &pData, &pTiles, &lTiles)
	if 0 != rc {
		errCStr := C.last_error(db.db)
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
