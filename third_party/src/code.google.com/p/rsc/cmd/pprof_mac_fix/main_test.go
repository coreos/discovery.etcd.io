// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var tests = []string{
	"testdata/mach_kernel_10_8_0",
	"testdata/mach_kernel_10_8_0_i386",
	"testdata/mach_kernel_11_4_2",
	"testdata/mach_kernel_12_4_0",
	"testdata/mach_kernel_13_0_0",
}

func TestAll(t *testing.T) {
	for _, tt := range tests {
		k := loadKernel(tt)
		errs := fixAnyVersion(k)
		if errs != nil {
			t.Errorf("%s: unrecognized kernel code", tt)
			for _, err := range errs {
				t.Errorf("%s", err)
			}
			continue
		}
		k1 := loadKernel(tt + "_fix")
		if !bytes.Equal(k.bsd_ast, k1.bsd_ast) {
			t.Errorf("%s: rewrite is incorrect\n", tt)
			n := 0
			for i := range k.bsd_ast {
				if k.bsd_ast[i] != k1.bsd_ast[i] {
					if n++; n > 20 {
						t.Logf("...")
						break
					}
					t.Logf("bsd_ast+%d: -%#02x +%#02x\n", i, k.bsd_ast[i], k1.bsd_ast[i])
				}
			}
			if err := ioutil.WriteFile("/tmp/newkernel", k.data, 0666); err != nil {
				t.Fatal(err)
			}
			want := gdbDisas(t, tt+"_fix", "bsd_ast")
			have := gdbDisas(t, "/tmp/newkernel", "bsd_ast")
			if err := ioutil.WriteFile("/tmp/want", want, 0666); err != nil {
				t.Fatal(err)
			}
			if err := ioutil.WriteFile("/tmp/have", have, 0666); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command("diff", "-u", "/tmp/have", "/tmp/want")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			os.Remove("/tmp/newkernel")
			os.Remove("/tmp/want")
			os.Remove("/tmp/have")
		}
	}
}

func gdbDisas(t *testing.T, file, name string) []byte {
	cmd := exec.Command("gdb", file)
	cmd.Stdin = strings.NewReader("disas " + name + "\n")
	disas, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gdb: %v", err)
	}
	return disas
}
