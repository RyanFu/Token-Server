package trick

import (
	"testing"
	"bytes"
)

func TestEqual(t *testing.T) {
	str := "testing"
	byte := []byte("testing")
	if !Equal(str, byte) {
		t.Fail()
	}
}

func BenchmarkEqual(b *testing.B) {
	str := "testing"
	byte := []byte("testing")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equal(str, byte)
	}
}

func TestString2Bytes(t *testing.T) {
	str := "testing"
	byte := []byte("testing")
	if bytes.Compare(String2Bytes(str), byte) != 0{
		t.Fail()
	}
}

func BenchmarkString2Bytes(b *testing.B) {
	str := "testing"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		String2Bytes(str)
	}
}

func TestBytes2String(t *testing.T) {
	str := "testing"
	byte := []byte("testing")
	if Bytes2String(byte) != str {
		t.Fail()
	}
}

func BenchmarkBytes2String(b *testing.B) {
	byte := []byte("testing")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Bytes2String(byte)
	}
}
