package swriter

import (
	"testing"
)

const xmaxInt = int(^uint(0) >> 1)

func testWrite(testNo int, initCap int, p []byte, t *testing.T) {

	sw := New(initCap)
	if sw.Len() != 0 || sw.Cap() != initCap {
		t.Fatalf("[test%d]invalid cap %d or len %d", testNo, sw.Cap(), sw.Len())
	}

	var n int
	var err error

	// capacityをまたがないwrite
	loop := initCap / len(p)
	for i := 0; i < loop; i++ {
		n, err = sw.Write(p)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(p) {
			t.Fatalf("[test%d]invalid return size %d", testNo, len(p))
		}
		if sw.Len() != len(p)*(i+1) || sw.Cap() != initCap {
			t.Fatalf("[test%d]invalid cap %d or len %d", testNo, sw.Cap(), sw.Len())
		}
	}

	// capacityをまたぐwrite
	n, err = sw.Write(p)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(p) {
		t.Fatalf("[test%d]invalid return size %d", testNo, len(p))
	}
	if sw.Len() != len(p)*(loop+1) || sw.Cap()%initCap != 0 {
		t.Fatalf("[test%d]invalid cap %d or len %d", testNo, sw.Cap(), sw.Len())
	}

	// bufの中身が正しいか確認
	out := make([]byte, sw.Len())
	sw.CopyTo(out)
	for i, v := range out {
		if v != p[i%len(p)] {
			t.Fatalf("[test%d]invalid sw val[%d] %d", testNo, i, v)
		}
	}

}

func TestWriter(t *testing.T) {

	testWrite(1, 4, []byte{1, 2}, t)
	testWrite(2, 5, []byte{1, 2}, t)
	testWrite(3, 1, []byte{1, 2, 3}, t)

	{
		initCap := 3
		sw := New(initCap)
		if sw.Len() != 0 || sw.Cap() != initCap {
			t.Fatalf("invalid cap %d or len %d", sw.Cap(), sw.Len())
		}
		sw.Grow(3)
		if sw.Len() != 3 || sw.Cap() != initCap {
			t.Fatalf("invalid cap %d or len %d", sw.Cap(), sw.Len())
		}
		sw.Grow(3)
		if sw.Len() != 6 || sw.Cap() != initCap*2 {
			t.Fatalf("invalid cap %d or len %d", sw.Cap(), sw.Len())
		}
	}
	{
		initCap := 3
		sw := New(initCap)
		sw.Write([]byte{1, 2})
		sw.Write([]byte{1, 2})
		sw.Write([]byte{1, 2})
		sw.Reset()
		if sw.Len() != 0 || sw.Tell() != 0 {
			t.Fatalf("invalid len %d off %d", sw.Len(), sw.Tell())
		}
	}
	{
		initCap := 3
		sw := New(initCap)
		sw.WriteAt([]byte{1, 2}, 2)
		if sw.Len() != 4 {
			t.Fatalf("invalid len %d", sw.Len())
		}
		if sw.Cap() != initCap*2 {
			t.Fatalf("invalid cap %d", sw.Cap())
		}
		if sw.Tell() != 4 {
			t.Fatalf("invalid Offset %d", sw.Tell())
		}
	}

}
