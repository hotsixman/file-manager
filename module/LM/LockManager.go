package LM

import (
	"file-manager/module/types"
	"path/filepath"
	"sync"
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
		return 0, &types.LockManagerException{Code: types.ExceptionStatus.INVALID_STATUS}
	}

	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		return 0, &types.LockManagerException{Code: types.ExceptionStatus.PATH_IS_NOT_ABS}
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	// check
	currentLockStatus := this.locks[path]
	if currentLockStatus > 0 {
		return currentLockStatus, &types.LockManagerException{Code: types.ExceptionStatus.ALREADY_LOCKED}
	}

	// ancestors check
	ancestors := getAncestors(path)
	for _, ancestor := range ancestors {
		ancestorLockStatus := this.locks[ancestor]
		if ancestorLockStatus > 0 {
			return 5, &types.LockManagerException{Code: types.ExceptionStatus.ANCESTOR_LOCKED}
		}
	}

	/*
		@update 0.1.0
		자손이 잠겨있어도 잠굴 수 있음.
		이 경우 자손이 잠금을 풀어도 조상은 잠금이 풀리지 않음.
		이 때 다시 자손을 잠구는 것은 불가능.
		// decendents check
		if this.ancestorLocks[path] > 0 {
			return 5, &types.LockManagerException{Code: types.ExceptionStatus.DECENDENT_LOCKED}
		}
	*/

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
		return &types.LockManagerException{Code: types.ExceptionStatus.PATH_IS_NOT_ABS}
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	// check
	currentLockStatus := this.locks[path]
	if currentLockStatus == 0 {
		return &types.LockManagerException{Code: types.ExceptionStatus.ALREADY_UNLOCKED}
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

func (this *LockManager) Check(path string) (*LockInfo, error) {
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		return nil, &types.LockManagerException{Code: types.ExceptionStatus.PATH_IS_NOT_ABS}
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
		Blocked:         ancestorLocked || (currentLockStatus > 0),
		AncestorLocked:  ancestorLocked,
		DecendentLocked: decendentLocked,
	}, nil
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
