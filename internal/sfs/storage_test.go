package sfs_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sfs "github.com/robherley/shallow-fetch-sha/internal/sfs"
)

var _ = Describe("Storage", func() {
	Describe("NewFileSystem", func() {
		Context("In disk mode", func() {
			It("should successfully initialize when specifying a relative directory", func() {
				testRelDir := "./foobar"
				expectedAbs, _ := filepath.Abs(testRelDir)

				fs, err := sfs.NewFileSystem(testRelDir, sfs.DiskMode)

				Expect(fs).To(Not(BeNil()))
				Expect(err).To(BeNil())
				Expect(fs.GetWorkTree()).To(Not(BeNil()))
				Expect(fs.GetStorage()).To(Not(BeNil()))
				Expect(fs.GetWorkTree().Root()).To(Equal(expectedAbs))
			})

			It("should successfully initialize when specifying an absolute directory", func() {
				expectedAbs := "/tmp/foobar"

				fs, err := sfs.NewFileSystem(expectedAbs, sfs.DiskMode)

				Expect(fs).To(Not(BeNil()))
				Expect(err).To(BeNil())
				Expect(fs.GetWorkTree()).To(Not(BeNil()))
				Expect(fs.GetStorage()).To(Not(BeNil()))
				Expect(fs.GetWorkTree().Root()).To(Equal(expectedAbs))
			})
		})

		Context("In memory mode", func() {
			It("should successfully initialize, regardless of directory", func() {
				fs, err := sfs.NewFileSystem("", sfs.MemoryMode)

				Expect(fs).To(Not(BeNil()))
				Expect(err).To(BeNil())
				Expect(fs.GetWorkTree()).To(Not(BeNil()))
				Expect(fs.GetStorage()).To(Not(BeNil()))
				Expect(fs.GetWorkTree().Root()).To(Equal("/"))
			})
		})

		Context("In random mode", func() {
			It("should fail", func() {
				fs, err := sfs.NewFileSystem(".", sfs.FSMode("not-a-real-mode"))

				Expect(fs).To(BeNil())
				Expect(err).To(Not(BeNil()))
			})
		})
	})
})
