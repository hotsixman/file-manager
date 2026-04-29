package types

import "strconv"

type LockManagerException struct {
	Code int
}

func (e *LockManagerException) Error() string {
	msg := ""
	switch e.Code {
	case ExceptionStatus.INVALID_STATUS:
		msg = "invalid status"
	case ExceptionStatus.PATH_IS_NOT_ABS:
		msg = "path is not absolute"
	case ExceptionStatus.ALREADY_LOCKED:
		msg = "already locked"
	case ExceptionStatus.ANCESTOR_LOCKED:
		msg = "ancestor is locked"
	case ExceptionStatus.DECENDENT_LOCKED:
		msg = "descendant is locked"
	default:
		msg = "unknown error"
	}
	return "LockError " + strconv.Itoa(e.Code) + ": " + msg
}
