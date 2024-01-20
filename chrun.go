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
        cg()
		cmd := exec.Command("/proc/self/exe", append([]string{"chroot"}, os.Args[2:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWNS |
				syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWIPC |
				syscall.CLONE_NEWPID |
				syscall.CLONE_NEWNET |
				syscall.CLONE_NEWUSER,
			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getuid(),
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getgid(),
					Size:        1,
				},
			},
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


	fmt.Println(os.Getpid())
	cgdir := "/sys/fs/cgroup"
	newCgroupv2 := filepath.Join(cgdir, "ns")
	must(os.MkdirAll(newCgroupv2, 0755))
	defer os.RemoveAll(newCgroupv2)
	procsFile := filepath.Join(newCgroupv2, "cgroup.procs")
	must(ioutil.WriteFile(procsFile, []byte(strconv.Itoa(os.Getpid())),0644))
	//now we limit the number of threads in this container-process to 2
	// test in a shell if you can fork more than two times
	//pidFile := filepath.Join(newCgroupv2, "pids.max")
	//must(ioutil.WriteFile(pidFile, []byte(strconv.Itoa(11)), 0644))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
