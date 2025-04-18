package clipboard

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"

	"github.com/pechorka/dev-tools/pkg/errs"
)

func Read() ([]byte, error) {
	switch os := runtime.GOOS; os {
	case "linux":
		var aggErr error
		for name, args := range map[string][]string{
			"xclip": {"-selection", "clipboard", "-o"},
		} {
			out, err := runCommand(name, args...)
			if err != nil {
				aggErr = errs.Join(aggErr, err)
				continue
			}

			return out, nil
		}

		return nil, aggErr
	case "darwin":
		return runCommand("pbpaste")
	case "windows":
		return runCommand("powershell", "-NoProfile", "-Command", "Get-Clipboard")
	default:
		return nil, errs.New("unsupported os %s", os)
	}
}

func runCommand(name string, args ...string) ([]byte, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, errs.Wrap(err, "failed to lookup program %s", name)
	}

	cmd := exec.Command(path, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
