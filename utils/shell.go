package utils

import (
	"os/exec"
	"strings"

	"github.com/weirwei/codereview/log"
)

// ShellExec Exec Shell Command
func ShellExec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	log.Debugf("command exec:%s", cmd.String())
	result, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Shell Exec failed: command: %s, err: %s, output: %s", cmd.String(), err.Error(), result)
		return "", err
	}
	log.Debugf("command exec result:%s", strings.ReplaceAll(string(result), "\n", "\\n"))
	return string(result), err
}
