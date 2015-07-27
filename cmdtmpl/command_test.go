package cmdtmpl_test

import (
	. "github.com/nullstyle/mcdev/cmdtmpl"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cmdtmpl.NewCommand", func() {
	var cmd *Command
	var err error

	Context("when `args` is an empty slice", func() {
		BeforeEach(func() {
			cmd, err = NewCommand([]string{})
		})

		It("returns an error", func() {
			Expect(err).To(Equal(ErrInvalidCommand))
		})
	})

	Context("when `args` has a single element", func() {
		BeforeEach(func() {
			cmd, err = NewCommand([]string{"echo"})
		})
		It("returns without error", func() {
			Expect(err).To(BeNil())
		})
		It("assigns the single element to the Cmd attribute", func() {
			Expect(cmd.Cmd).To(Equal("echo"))
		})
	})

	Context("when `args` has multiple elements", func() {
		BeforeEach(func() {
			cmd, err = NewCommand([]string{"echo", "building", "{{.Pkg}}"})
		})
		It("returns without error", func() {
			Expect(err).To(BeNil())
		})
		It("assigns the first element to the Cmd attribute", func() {
			Expect(cmd.Cmd).To(Equal("echo"))
		})

		It("parses the rest to the Args template slice", func() {
			Expect(cmd.Args).To(HaveLen(2))
		})
	})

	Context("when `args` has an unparseable element", func() {
		BeforeEach(func() {
			cmd, err = NewCommand([]string{"echo", "{{.Pkg"})
		})

		It("returns an error", func() {
			Expect(err).ToNot(BeNil())
		})
	})
})
