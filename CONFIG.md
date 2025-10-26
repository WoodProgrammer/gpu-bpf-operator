{
  "action": "CREATE|DELETE|UPDATE"
  "hash": "sha256:5f4c...",
  "policies": [{
    "id": "cuda-malloc@d34d",
    "libPath": "/usr/lib/x86_64-linux-gnu/libcudart.so",
    "mode": "pidwatch",
    "processRegex": "^(python|trainer)$",
    "functions": ["cudaMalloc","cudaFree","cudaMemcpy"],
    "output": { "format": "ndjson" }
  }]
}