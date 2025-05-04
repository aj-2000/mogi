package app

/*
#cgo LDFLAGS: -L../renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "../renderer/include/renderer.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strconv"
	"sync"
	"unsafe"
)

// TODO: should we use a different package for font manager, renderer and texture manager?
type FontData = C.FontData

type fontManager struct {
	mu    sync.Mutex
	cache map[string]*FontData
}

func NewFontManager() *fontManager {
	return &fontManager{cache: make(map[string]*FontData)}
}

func (fm *fontManager) load(path string, size float32) (*FontData, error) {
	key := path + "|" + strconv.FormatFloat(float64(size), 'f', -1, 32)
	fm.mu.Lock()
	if font, ok := fm.cache[key]; ok {
		fm.mu.Unlock()
		return font, nil
	}
	fm.mu.Unlock()

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	font := C.load_font(cPath, C.float(size))
	if font == nil {
		return nil, fmt.Errorf("failed to load font: %s", path)
	}

	fm.mu.Lock()
	fm.cache[key] = font
	fm.mu.Unlock()
	return font, nil
}

func (fm *fontManager) unload(path string, size float32) {
	key := path + "|" + strconv.FormatFloat(float64(size), 'f', -1, 32)
	fm.mu.Lock()
	if font, ok := fm.cache[key]; ok {
		C.destroy_font(font)
		delete(fm.cache, key)
	}
	fm.mu.Unlock()
}

func (fm *fontManager) destroy() {
	fm.mu.Lock()
	for _, font := range fm.cache {
		C.destroy_font(font)
	}
	fm.cache = nil
	fm.mu.Unlock()
}
