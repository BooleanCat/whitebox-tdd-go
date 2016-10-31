package filething

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFilething(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filething Suite")
}
