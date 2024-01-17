package main

import (
	"context"
	"os"
	"os/exec"
	"syscall"
    "path/filepath"
	"github.com/codeclysm/extract"
	"io/ioutil"
	"strconv"
	"fmt"
)

func main() {
	switch os.Args[1] {
	case "run":
		os.Mkdir("./assets/tmp/", 0750)
		defer os.RemoveAll("./assets/tmp/")
		cmd := exec.Command("/proc/self/exe", append([]string{"chroot"}, os.Args[2:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
			Unshareflags: syscall.CLONE_NEWNS,
		}
		must(cmd.Run())
	case "pull":
		image := os.Args[2]
		pullImage(image)
	case "chroot":
		chroot()
	default:
		panic("what?")
	}

}

func pullImage(image string) {
	cmd := exec.Command("./pull", image)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func chroot() {

	func() error {
		r, err := os.Open("/home/coder/gocontainerrt/assets/alpine.tar.gz")
		if err != nil {
			return err
		}
		defer r.Close()
		ctx := context.Background()
		return extract.Archive(ctx, r, "./assets/tmp/", nil)
	}()
	must(syscall.Chroot("./assets/tmp"))
	must(syscall.Chdir("/"))
	//cg()
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("mycontainer")))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	must(cmd.Run())

	must(syscall.Unmount("proc", 0))

}

func cg() {

	pid := os.Getpid()
	fmt.Println(pid)
	cgdir := "/sys/fs/cgroup"
	procsFile := filepath.Join(cgdir, "cgroup.procs")
	ioutil.WriteFile(procsFile, []byte(strconv.Itoa(pid)),0700)

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
