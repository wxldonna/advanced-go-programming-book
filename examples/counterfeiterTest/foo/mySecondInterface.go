package foo

// You only need **one** of these per package!
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o libraryTest/fake_interface3.go . MySpecialInterface3
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o libraryTest/fake_interface3.go . MySpecialInterface4

type MySpecialInterface3 interface {
	DoThings1(string, uint64) (int, error)
}

//counterfeiter:generate . MySpecialInterface2
type MySpecialInterface4 interface {
	DoThings2(string, uint64) (int, error)
}
