Whitebox testing with Golang and Ginkgo
=======================================

The problem
-----------

Let's start with a straightforward problem. There exists a type called `FileThing` and you've been asked to add a method to it: `Remove`. The requirements for this function are:

- It can be used like `filething.New(someFilePath).Remove()`
- It will delete the file at `someFilePath`
- It will return an error if deletion fails
- It will not return an error if the file doesn't exist

You decide to test drive this (right?) and start with the easiest tests first - you might end up with something like this:

```go
package filething_test

import (
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
	})
})

func createSomeTempFile() string {
	tempFile, err := ioutil.TempFile("", "")
	Expect(err).NotTo(HaveOccurred())
	defer tempFile.Close()
	return tempFile.Name()
}
```

And your implementation might look something like this:

```go
package filething

type FileThing struct {
	Path string
}

func New(path string) FileThing {
	return FileThing{Path: path}
}

func (fileThing FileThing) Remove() error {
	os.Remove(fileThing.Path)
	return nil
}
```

Notice how you've managed to satisfy most of your requirements without even considering the error that comes back from `os.Remove`. There's just one more requirement to implement here: "It will return an error if deletion fails".

So how do you get deletion to fail? You could remove the file... oh but that's not considered an error. How about you `chmod` the file and make it not writeable... that doesn't quite feel right and besides it only addresses one kind of failure, you want to deal with all kinds of failure.

What you really want is to dictate in very explicit terms the kind of error you want to happen and check that it bubbles up.

Enter whitebox testing.

A solution
----------

A wonderful thing about go is that function signatures are types - including, from the example above, `os.Remove`. You need a way to modify the behaviour of this function while hiding implementation details from consumers of `FileThing`. This implies a private member variable only visible to the implementor.

It's perfectly fine to call your test package `filething`, now that it is part of the implementation package it has access to private stuff inside that same package.

So you imagine a way to control the behaviour of `FileThing.Remove`'s internals, a way to control its `Remover`:

```go
type Remover func(string) error
```

You write your tests:

```go
Describe("#Remove", func() {
  ...

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

func failToRemove(path string) error {
	return errors.New("I failed")
}
```

You want this to be transparent to the consumer of `FileThing`, so in making the tests pass, you add some default value:

```go
package filething

import "os"

type Remover func(string) error

type FileThing struct {
	Path   string
	remove Remover
}

func New(path string) FileThing {
	return FileThing{
		Path:   path,
		remove: os.Remove,
	}
}

func (fileThing FileThing) Remove() error {
	err := fileThing.remove(fileThing.Path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
```

Just like that, tests are passing and you have a way to be very explicit in your about failures from other functions.
