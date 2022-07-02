package apierror

import (
	"net"
	"net/url"
	"os"
	"syscall"
)

// IsNetErr 是否是网络错误
func IsNetErr(err error) bool {
	b := underlyingErrorIs(err, syscall.ECONNREFUSED)
	if !b {
		b = underlyingErrorIs(err, syscall.ECONNABORTED)
	}
	return b
}

// underlyingErrorIs 递归调用底层错误类型，判断是否是目标类型
func underlyingErrorIs(err, target error) bool {
	err = underlyingError(err)
	if err == target {
		return true
	} else if err == nil {
		return false
	} else {
		return underlyingErrorIs(err, target)
	}
}

// underlyingError returns the underlying error for known os error types.
func underlyingError(err error) error {
	switch err := err.(type) {
	case *url.Error:
		return err.Err
	case *net.OpError:
		return err.Err
	case *os.SyscallError:
		return err.Err
	}
	return nil
}
