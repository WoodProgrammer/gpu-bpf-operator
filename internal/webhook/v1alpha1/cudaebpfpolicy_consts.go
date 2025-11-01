package v1alpha1

var ALLOWED_GPU_EVENTS = []string{
	"nvidia_open",
	"nvidia_unlocked_ioctl",
	"nvidia_mmap",
	"nvidia_isr",
	"nvidia_isr_kthread_bh",
}
