package namespace

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

func Init(newRoot string) error {
	logrus.Info("Init namespace start")

	if err := pivotRoot(newRoot); err != nil {
		return err
	}

	if err := syscall.Sethostname([]byte("judge")); err != nil {
		return err
	}

	logrus.Info("Init namespace end")
	return nil
}

func pivotRoot(newRoot string) error {
	if err := syscall.Mount(newRoot, newRoot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	// create put_old directory
	rootOld := filepath.Join(newRoot, "/.pivot_root")
	if err := os.MkdirAll(rootOld, 0700); err != nil {
		return err
	}

	// call pivotRoot
	if err := syscall.PivotRoot(newRoot, rootOld); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	// umount put_old, which now lives at /.pivot_root
	rootOld = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(rootOld, syscall.MNT_DETACH); err != nil {
		return err
	}

	// remove put_old
	if err := os.RemoveAll(rootOld); err != nil {
		return err
	}

	return nil
}
