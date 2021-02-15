package intermediate

import (
	"os/exec"

	lowlevel "github.com/yourfavoritedev/learngo/first/concurrency/error-propagation/low-level"
	myError "github.com/yourfavoritedev/learngo/first/concurrency/error-propagation/my-error"
)

// IntermediateErr represents an intermediate-level error from LowLevelErr
type IntermediateErr struct {
	error
}

// RunJob is a boilerplate function that wraps an IntermediateErr
func RunJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := lowlevel.IsGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{error: myError.WrapError(err, "cannot run job %q: requisite binaries not available", id)}
	} else if isExecutable == false {
		return myError.WrapError(nil, "cannot run job %q: requisite binaries not executable", id)
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}
