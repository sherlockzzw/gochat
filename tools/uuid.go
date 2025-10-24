package tools

import (
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
)

var (
	idLock = &sync.Mutex{}
)

// GenUUID -
func GenUUID() string {
	idLock.Lock()
	defer idLock.Unlock()
	v4 := uuid.NewV4()
	return strings.ReplaceAll(v4.String(), "-", "")
}
