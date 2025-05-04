package app

/*
#cgo LDFLAGS: -L../renderer/lib/Release -lrenderer -lglfw3 -lgdi32 -static
#include "../renderer/include/renderer.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

type textureManager struct {
	mu    sync.Mutex
	cache map[string]C.GLuint
}

func NewTextureManager() *textureManager {
	return &textureManager{cache: make(map[string]C.GLuint)}
}

func (tm *textureManager) load(path string) (C.GLuint, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tex, ok := tm.cache[path]; ok {
		return tex, nil
	}

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	tex := C.load_texture(cPath)
	if tex == 0 {
		return 0, fmt.Errorf("failed to load texture: %s", path)

	}
	return tex, nil
}

// TODO: smart unload
func (tm *textureManager) unload(path string) {
	tm.mu.Lock()
	if tex, ok := tm.cache[path]; ok {
		C.free_texture(tex)
		delete(tm.cache, path)
	}
	tm.mu.Unlock()
}

func (tm *textureManager) destroy() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	for path, tex := range tm.cache {
		C.free_texture(tex)
		delete(tm.cache, path)
	}
}
