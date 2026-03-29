package main

import (
	"path/filepath"
	"strconv"
	"sync"
)

const (
	StatusIdle       = 0
	StatusReading    = 1
	StatusMovingFrom = 2
	StatusMovingTo   = 3
	StatusDeleting   = 4
	StatusLocked     = 5
)

type LockManager struct {
	mu            sync.Mutex
	locks         map[string]int
	ancestorLocks map[string]int // 조상-자손 개수 map
}

func NewLockManager() *LockManager {
	return &LockManager{
		locks:         make(map[string]int),
		ancestorLocks: make(map[string]int),
	}
}

func (this *LockManager) Lock(path string, status int) (int, error) {
	if status < 1 || status > 5 {
		return 0, &LockManagerException{
			code: INVALID_STATUS,
		}
	}

	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		return 0, &LockManagerException{
			code: PATH_IS_NOT_ABS,
		}
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	// check
	currentLockStatus := this.locks[path]
	if currentLockStatus > 0 {
		return currentLockStatus, &LockManagerException{
			code: ALREADY_LOCKED,
		}
	}

	// ancestors check
	ancestors := getAncestors(path)
	for _, ancestor := range ancestors {
		ancestorLockStatus := this.locks[ancestor]
		if ancestorLockStatus > 0 {
			return 5, &LockManagerException{
				code: ANCESTOR_LOCKED,
			}
		}
	}

	// decendents check
	if this.ancestorLocks[path] > 0 {
		return 5, &LockManagerException{
			code: DECENDENT_LOCKED,
		}
	}

	// ancestorLocks 관리
	for _, ancestor := range ancestors {
		this.ancestorLocks[ancestor]++
	}

	// lock
	this.locks[path] = status

	return status, nil
}

func (this *LockManager) Unlock(path string) error {
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		return &LockManagerException{
			code: PATH_IS_NOT_ABS,
		}
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	// check
	currentLockStatus := this.locks[path]
	if currentLockStatus == 0 {
		return &LockManagerException{
			code: ALREADY_UNLOCKED,
		}
	}

	// ancestors check
	ancestors := getAncestors(path)

	// ancestorLocks 관리
	for _, ancestor := range ancestors {
		if this.ancestorLocks[ancestor] > 0 {
			this.ancestorLocks[ancestor]--
			if this.ancestorLocks[ancestor] == 0 {
				delete(this.ancestorLocks, ancestor)
			}
		}
	}

	// lock
	delete(this.locks, path)

	return nil
}

type LockInfo struct {
	Status          int  `json:"status"`
	Blocked         bool `json:"blocked"`
	AncestorLocked  bool `json:"ancestorLocked"`
	DecendentLocked bool `json:"decendentLocked"`
}

func (this *LockManager) CheckLocked(path string) (*LockInfo, error) {
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		return nil, &LockManagerException{
			code: PATH_IS_NOT_ABS,
		}
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	currentLockStatus := this.locks[path]

	ancestorLocked := false
	ancestors := getAncestors(path)
	for _, ancestor := range ancestors {
		ancestorLockStatus := this.locks[ancestor]
		if ancestorLockStatus > 0 {
			ancestorLocked = true
			break
		}
	}

	decendentLocked := (this.ancestorLocks[path] > 0)

	return &LockInfo{
		Status:          currentLockStatus,
		Blocked:         ancestorLocked || decendentLocked,
		AncestorLocked:  ancestorLocked,
		DecendentLocked: decendentLocked,
	}, nil
}

const (
	INVALID_STATUS   = 1
	PATH_IS_NOT_ABS  = 2
	ALREADY_LOCKED   = 3
	ANCESTOR_LOCKED  = 4
	DECENDENT_LOCKED = 5
	ALREADY_UNLOCKED = 6
)

type LockManagerException struct {
	code int
}

func (e *LockManagerException) Error() string {
	msg := ""
	switch e.code {
	case INVALID_STATUS:
		msg = "invalid status"
	case PATH_IS_NOT_ABS:
		msg = "path is not absolute"
	case ALREADY_LOCKED:
		msg = "already locked"
	case ANCESTOR_LOCKED:
		msg = "ancestor is locked"
	case DECENDENT_LOCKED:
		msg = "descendant is locked"
	default:
		msg = "unknown error"
	}
	return "LockError " + strconv.Itoa(e.code) + ": " + msg
}

func getAncestors(path string) []string {
	path = filepath.Clean(path)
	if path == "/" {
		return []string{}
	}

	ancestors := []string{}
	curr := path
	for {
		parent := filepath.Dir(curr)
		if parent == curr { // 더 이상 올라갈 곳이 없음 (루트 도달)
			break
		}
		ancestors = append(ancestors, parent)
		curr = parent
	}
	return ancestors
}
