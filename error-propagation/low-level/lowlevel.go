package lowlevel

import (
	"os"

	myError "github.com/yourfavoritedev/learngo/first/concurrency/error-propagation/my-error"
)

// LowLevelErr represents a low-level error message
type LowLevelErr struct {
	error
}

// IsGloballyExec encapsulate a LowLevelErr
func IsGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{error: myError.WrapError(err, err.Error())}
	}
	return info.Mode().Perm()&0100 == 0100, nil
}
