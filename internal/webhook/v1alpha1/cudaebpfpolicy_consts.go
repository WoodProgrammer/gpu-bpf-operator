package v1alpha1

var ALLOWED_CUDA_EVENTS = []string{
	"cudaMalloc",
	"cudaFree",
	"cudaMemcpy",
	"cudaLaunchKernel",
	"cudaStreamCreate",
	"cudaStreamSynchronize",
	"cudaGetDevice",
	"cudaSetDevice",
	"cudaEventCreate",
	"cudaEventRecord",
	"cudaEventSynchronize",
}
