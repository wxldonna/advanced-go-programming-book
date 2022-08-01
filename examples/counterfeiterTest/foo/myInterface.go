package foo

// You only need **one** of these per package!
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// You will add lots of directives like these in the same package...
//counterfeiter:generate . MySpecialInterface1
type MySpecialInterface1 interface {
	DoThings1(string, uint64) (int, error)
}

//counterfeiter:generate . MySpecialInterface2
type MySpecialInterface2 interface {
	DoThings2(string, uint64) (int, error)
}
