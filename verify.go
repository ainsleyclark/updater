package updater

import (
	"errors"
	"os/exec"
	"strings"
)

// verifyInstallation verifies if the executable is installed correctly
// we are going to run the newly installed program by running it with -version
// if it outputs the good version then we assume the installation is good
func (u *Updater) verifyInstallation() error {
	latestVersion, err := u.LatestVersion()
	if err != nil {
		return err
	}
	executable, err := u.pkg.GetExecutable()
	if err != nil {
		return err
	}

	cmd := exec.Cmd{
		Path: executable,
		Args: []string{executable, "-version"},
	}

	output, err := cmd.Output()
	if err != nil {
		return err
	}
	strOutput := string(output)

	if !strings.Contains(strOutput, latestVersion) {
		return errors.New("version mismatch in latest updated")
	}
	return nil
}
