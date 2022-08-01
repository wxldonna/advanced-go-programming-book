// +build !amd64

package trc

func GoroutineID() int64 {
	return -1
}
