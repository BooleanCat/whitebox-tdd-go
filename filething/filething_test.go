package filething

import (
	"errors"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileThing", func() {
	var (
		fileThing FileThing
		someFile  string
	)

	BeforeEach(func() {
		someFile = createSomeTempFile()
		fileThing = New(someFile)
	})

	AfterEach(func() {
		os.Remove(someFile)
		Expect(someFile).NotTo(BeAnExistingFile())
	})

	Describe("#Remove", func() {
		var removeErr error

		JustBeforeEach(func() {
			removeErr = fileThing.Remove()
		})

		It("does not return an error", func() {
			Expect(removeErr).NotTo(HaveOccurred())
		})

		It("removes a file", func() {
			Expect(someFile).NotTo(BeAnExistingFile())
		})

		Context("when FileThing.Path doesn't exist", func() {
			BeforeEach(func() {
				err := os.Remove(someFile)
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not return an error", func() {
				Expect(removeErr).NotTo(HaveOccurred())
			})
		})

		Context("when deleting FileThing.Path fails", func() {
			BeforeEach(func() {
				fileThing.remove = failToRemove
			})

			It("returns an error", func() {
				Expect(removeErr).To(HaveOccurred())
			})

			It("reports the correct error", func() {
				Expect(removeErr).To(MatchError("I failed"))
			})
		})
	})
})

func createSomeTempFile() string {
	tempFile, err := ioutil.TempFile("", "")
	Expect(err).NotTo(HaveOccurred())
	defer tempFile.Close()
	return tempFile.Name()
}

func failToRemove(path string) error {
	return errors.New("I failed")
}
