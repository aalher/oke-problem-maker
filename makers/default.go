package makers

import (
	"os"

	"k8s.io/klog/v2"
)

func init() {
	// Docker related problems
	ProblemGenerators["DockerHung"] = makeDockerHung
    ProblemGenerators["DockerContainerStartupFailure"] = makeDockerContainerStartupFailure

	// Filesystem related problems
	ProblemGenerators["Ext4FilesystemError"] = makeFilesystemError

	// IO Errors
	ProblemGenerators["IOErrors"] = makeIOErrors

	// Memory Errors
	ProblemGenerators["MemoryReadError"] = makeMemoryReadError

	// OOM Errors
	ProblemGenerators["OOMKilling"] = makeOOMKill

	// KernelOops
	ProblemGenerators["KernelOops"] = makeKernelOops

	// Net Errors
	ProblemGenerators["UnregisterNetDevice"] = makeUnregisterNetDevice
}

func makeDockerHung() {
	const msg = `INFO: task docker:20744 blocked for more than 120 seconds.
      Tainted: G         C    3.16.0-4-amd64 #1
"echo 0 > /proc/sys/kernel/hung_task_timeout_secs" disables this message.
docker          D ffff8801a8f2b078     0 20744      1 0x00000000
 ffff8801a8f2ac20 0000000000000082 0000000000012f00 ffff880057a17fd8
 0000000000012f00 ffff8801a8f2ac20 ffffffff818bb4a0 ffff880057a17d80
 ffffffff818bb4a4 ffff8801a8f2ac20 00000000ffffffff ffffffff818bb4a8
Call Trace:
 [<ffffffff81510915>] ? schedule_preempt_disabled+0x25/0x70
 [<ffffffff815123c3>] ? __mutex_lock_slowpath+0xd3/0x1c0
 [<ffffffff815124cb>] ? mutex_lock+0x1b/0x2a
 [<ffffffff814175bc>] ? copy_net_ns+0x6c/0x130
 [<ffffffff8108bdf4>] ? create_new_namespaces+0xf4/0x180
 [<ffffffff8108beec>] ? copy_namespaces+0x6c/0x90
 [<ffffffff810654f6>] ? copy_process.part.25+0x966/0x1c30
 [<ffffffff81066991>] ? do_fork+0xe1/0x390
 [<ffffffff811c442c>] ? __alloc_fd+0x7c/0x120
 [<ffffffff81514079>] ? stub_clone+0x69/0x90
 [<ffffffff81513d0d>] ? system_call_fast_compare_end+0x10/0x15`

	writeKernelMessageOrDie(msg)
}

func makeDockerContainerStartupFailure() {
    const msg = "Error response from daemon: OCI runtime create failed: container_linux.go:345"

    writeKernelMessageOrDie(msg)
}

func makeFilesystemError() {
	const ext4ErrorTrigger = "/sys/fs/ext4/sda1/trigger_fs_error"

	msg := []byte("fake filesystem error from problem-maker")
	err := os.WriteFile(ext4ErrorTrigger, msg, 0200)
	if err != nil {
		klog.Fatalf("Failed writing log to %q: %v", ext4ErrorTrigger, err)
	}
}

func makeIOErrors() {
	const msg = "Buffer I/O error on dev sda1, logical block 123456, async page read"

	writeKernelMessageOrDie(msg)
}

func makeMemoryReadError() {
	const msg = "CE memory read error 0x0000000000000020"

	writeKernelMessageOrDie(msg)
}

func makeOOMKill() {
	const msg = `Memory cgroup out of memory: Kill process 1012 (heapster) score 1035 or sacrifice child
Killed process 1012 (heapster) total-vm:327128kB, anon-rss:306328kB, file-rss:11132kB, shmem-rss:12345kB`

	writeKernelMessageOrDie(msg)
}

func makeKernelOops() {
	const msg = "BUG: unable to handle kernel NULL pointer dereference at 0000000000000010"

	writeKernelMessageOrDie(msg)
}

func makeUnregisterNetDevice() {
	const msg = "unregister_netdevice: waiting for eth0 to become free. Usage count = 1"

	writeKernelMessageOrDie(msg)
}
