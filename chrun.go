package main

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"github.com/containerd/cgroups/v3/cgroup2"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/codeclysm/extract"
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
	cg()
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

	var cgroupV2 bool
    if cgroups.Mode() == cgroups.Unified {
	cgroupV2 = true
    }
	fmt.Println(cgroupV2)
	res := cgroup2.Resources{}
	// dummy PID of -1 is used for creating a "general slice" to be used as a parent cgroup.
	// see https://github.com/containerd/cgroups/blob/1df78138f1e1e6ee593db155c6b369466f577651/v2/manager.go#L732-L735
	m, err := cgroup2.NewSystemd("/", "my-cgroup-abc.slice", -1, &res)
	//m, err := cgroup2.LoadSystemd("/", "my-cgroup-abc.slice")
	// https://www.kernel.org/doc/html/v5.0/admin-guide/cgroup-v2.html#threads
	cgType, err := m.GetType()
	fmt.Println(cgType)
	defer m.DeleteSystemd()

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
