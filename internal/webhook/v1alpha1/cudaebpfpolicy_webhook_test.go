/*
Copyright 2025 WoodProgrammer.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gpuv1alpha1 "github.com/WoodProgrammer/gpu-bpf-operator/api/v1alpha1"
)

var _ = Describe("CudaEBPFPolicy Webhook", func() {
	var (
		obj       *gpuv1alpha1.CudaEBPFPolicy
		oldObj    *gpuv1alpha1.CudaEBPFPolicy
		validator CudaEBPFPolicyCustomValidator
		defaulter CudaEBPFPolicyCustomDefaulter
		ctx       context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		obj = &gpuv1alpha1.CudaEBPFPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-policy",
				Namespace: "default",
			},
		}
		oldObj = &gpuv1alpha1.CudaEBPFPolicy{}
		validator = CudaEBPFPolicyCustomValidator{}
		Expect(validator).NotTo(BeNil(), "Expected validator to be initialized")
		defaulter = CudaEBPFPolicyCustomDefaulter{}
		Expect(defaulter).NotTo(BeNil(), "Expected defaulter to be initialized")
	})

	Context("When creating CudaEBPFPolicy under Defaulting Webhook", func() {
		It("Should apply defaults for outputFormat when empty", func() {
			By("simulating a scenario where outputFormat is empty")
			obj.Spec.OutputFormat = ""
			obj.Spec.Mode = ""

			By("calling the Default method to apply defaults")
			err := defaulter.Default(ctx, obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that the default values are set")
			Expect(obj.Spec.OutputFormat).To(Equal("ndjson"))
		})

		It("Should apply defaults for mode when empty", func() {
			By("simulating a scenario where mode is empty")
			obj.Spec.Mode = ""
			obj.Spec.OutputFormat = ""

			By("calling the Default method to apply defaults")
			err := defaulter.Default(ctx, obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that the default values are set")
			Expect(obj.Spec.Mode).To(Equal("pidwatch"))
		})

		It("Should not override existing values", func() {
			By("simulating a scenario where values are already set")
			obj.Spec.OutputFormat = "prometheus"
			obj.Spec.Mode = "systemwide"

			By("calling the Default method")
			err := defaulter.Default(ctx, obj)
			Expect(err).NotTo(HaveOccurred())

			By("checking that the existing values are preserved")
			Expect(obj.Spec.OutputFormat).To(Equal("prometheus"))
			Expect(obj.Spec.Mode).To(Equal("systemwide"))
		})
	})

	Context("When creating or updating CudaEBPFPolicy under Validating Webhook", func() {
		It("Should deny creation if functions array is empty", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least one function must be specified"))
		})

		It("Should deny creation if function name is empty", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "",
					Kind: "uprobe",
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("function name must be specified"))
		})

		It("Should deny creation if function kind is invalid", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "invalidprobe",
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("kind"))
		})

		It("Should deny creation if duplicate function names exist", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
				{
					Name: "cudaStreamCreate",
					Kind: "uretprobe",
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Duplicate"))
		})

		It("Should deny creation if argument index is negative", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
					Args: []gpuv1alpha1.Arg{
						{
							Index: -1,
							Name:  "stream",
						},
					},
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("argument index must be non-negative"))
		})

		It("Should deny creation if argument name is empty", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
					Args: []gpuv1alpha1.Arg{
						{
							Index: 0,
							Name:  "",
						},
					},
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("argument name must be specified"))
		})

		It("Should deny creation if duplicate argument indices exist", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
					Args: []gpuv1alpha1.Arg{
						{
							Index: 0,
							Name:  "stream",
						},
						{
							Index: 0,
							Name:  "flags",
						},
					},
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Duplicate"))
		})

		It("Should deny creation if mode is invalid", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "invalidmode"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("mode"))
		})

		It("Should deny creation if outputFormat is invalid", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"
			obj.Spec.OutputFormat = "invalidformat"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("output"))
		})

		It("Should deny creation if libPath is empty", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
			}
			obj.Spec.LibPath = ""
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("libPath must be specified"))
		})

		It("Should deny creation if image is empty", func() {
			By("simulating an invalid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = ""
			obj.Spec.Mode = "pidwatch"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("image must be specified"))
		})

		It("Should admit creation with valid spec", func() {
			By("simulating a valid creation scenario")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
				{
					Name: "cudaStreamSynchronize",
					Kind: "uretprobe",
					Args: []gpuv1alpha1.Arg{
						{
							Index: 0,
							Name:  "stream",
						},
					},
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "pidwatch"
			obj.Spec.OutputFormat = "ndjson"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should admit creation with all valid probe kinds", func() {
			By("simulating a valid creation scenario with all probe kinds")
			obj.Spec.Functions = []gpuv1alpha1.Function{
				{Name: "func1", Kind: "uprobe"},
				{Name: "func2", Kind: "uretprobe"},
				{Name: "func3", Kind: "kprobe"},
				{Name: "func4", Kind: "kretprobe"},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:latest"
			obj.Spec.Mode = "systemwide"
			obj.Spec.OutputFormat = "prometheus"

			By("validating the creation")
			_, err := validator.ValidateCreate(ctx, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should validate updates correctly", func() {
			By("simulating a valid update scenario")
			oldObj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
			}
			oldObj.Spec.LibPath = "/usr/lib/libcuda.so"
			oldObj.Spec.Image = "test-image:v1"
			oldObj.Spec.Mode = "pidwatch"

			obj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
				{
					Name: "cudaMemcpy",
					Kind: "uprobe",
				},
			}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:v2"
			obj.Spec.Mode = "systemwide"

			By("validating the update")
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Should deny update with invalid spec", func() {
			By("simulating an invalid update scenario")
			oldObj.Spec.Functions = []gpuv1alpha1.Function{
				{
					Name: "cudaStreamCreate",
					Kind: "uprobe",
				},
			}
			oldObj.Spec.LibPath = "/usr/lib/libcuda.so"
			oldObj.Spec.Image = "test-image:v1"
			oldObj.Spec.Mode = "pidwatch"

			obj.Spec.Functions = []gpuv1alpha1.Function{}
			obj.Spec.LibPath = "/usr/lib/libcuda.so"
			obj.Spec.Image = "test-image:v2"
			obj.Spec.Mode = "pidwatch"

			By("validating the update")
			_, err := validator.ValidateUpdate(ctx, oldObj, obj)
			Expect(err).To(HaveOccurred())
		})
	})

})
