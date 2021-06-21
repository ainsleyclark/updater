package updater

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

var (
	// ErrVersionMisMatch is returned by verifyInstallation if
	// the new downloaded version did not output the same
	// version parsed to Update().
	ErrVersionMisMatch = errors.New("version mismatch in updated executable")
)

// verifyInstallation verifies if the executable is installed
// correctly. The downloaded executable is run with the
// flag -version.
// Returns ErrVersionMisMatch if the versions could n ot be
// matched.
func (u *Updater) verifyInstallation() error {
	latestVersion, err := u.LatestVersion()
	if err != nil {
		return err
	}

	executable, err := os.Executable()
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
		return ErrVersionMisMatch
	}

	return nil
}
