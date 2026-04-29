package types

type lockStatus struct {
	Idle       int
	Reading    int
	MovingFrom int
	MovingTo   int
	Deleting   int
	Locked     int
}

var LockStatus = lockStatus{0, 1, 2, 3, 4, 5}

type exceptionCode struct {
	INVALID_STATUS   int
	PATH_IS_NOT_ABS  int
	ALREADY_LOCKED   int
	ANCESTOR_LOCKED  int
	DECENDENT_LOCKED int
	ALREADY_UNLOCKED int
}

var ExceptionStatus = exceptionCode{1, 2, 3, 4, 5, 6}
