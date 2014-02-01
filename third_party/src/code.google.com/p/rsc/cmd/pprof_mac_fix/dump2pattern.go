// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program analyzes pprof_mac_fix -dump output to generate the
// "before" assembly listing and a suggested pattern template.

// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var dumpRE = regexp.MustCompile(`(?ms)^version: ([^0-9\n]+([0-9][^:\n]+)[^\n]+)\n.*\n\(gdb\)[^\n]*current_thread:(.*?)^\(gdb\).*^\(gdb\)[^\n]*bsd_ast:(.*?)^\(gdb\)`)

var (
	testAsm *os.File
)

func main() {
	log.SetFlags(0)
	data, err := ioutil.ReadFile("dump")
	if err != nil {
		log.Fatal(err)
	}

	m := dumpRE.FindStringSubmatch(string(data))
	if m == nil {
		log.Fatal("cannot parse -dump output")
	}

	vers := strings.Replace(m[2], ".", "_", -1)

	name := "testdata/mach_kernel_" + vers + ".s"
	if _, err := os.Stat(name); err == nil {
		log.Fatalf("%s already exists", name)
	}
	testAsm, err = os.Create(name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(testAsm, "_version:\n\t.ascii \"%s\\0\"\n", m[1])
	for _, name := range extraFuncs {
		fmt.Fprintf(testAsm, "\n.globl _%s\n_%s:\n\tret\n", name, name)
	}
	dodis("current_thread", m[3])
	dodis("bsd_ast", m[4])
}

var extraFuncs = []string{
	"main",
	"psignal_internal",
	"task_vtimer_clear",
	"task_vtimer_set",
}

var (
	instRE = regexp.MustCompile(`(?-s)^0x[0-9a-f]+\s+<(\w+)\+([0-9]+)>:\s+(.*)$`)
	hexRE  = regexp.MustCompile(`(?-s)^\s+([0-9a-f]{2}( [0-9a-f]{2})*)\s*$`)
	callRE = regexp.MustCompile(`callq?\s+0x[0-9a-f]+\s+<(\w+(?:\+[0-9]+)?)>`)
)

func dodis(name string, dis string) {
	fmt.Fprintf(testAsm, "\n.globl _%s\n_%s:\n", name, name)
	lines := strings.Split(dis, "\n")
	for len(lines) > 0 && !strings.HasPrefix(lines[0], "0x") {
		lines = lines[1:]
	}
	for len(lines) >= 2 && !strings.HasPrefix(lines[len(lines)-2], "0x") {
		lines = lines[:len(lines)-1]
	}
	if len(lines) < 2 || len(lines)%2 != 0 {
		log.Fatalf("cannot parse disassembly for %s", name)
	}
	nextOffset := 0
	lastCall := ""
	var insts []string
	var hexes [][]string
	for i := 0; i < len(lines); i += 2 {
		if !instRE.MatchString(lines[i]) || !hexRE.MatchString(lines[i+1]) {
			log.Fatalf("cannot parse disassembly for %s", name)
		}
		m := instRE.FindStringSubmatch(lines[i])
		name := m[1]
		offset, _ := strconv.Atoi(m[2])
		inst := m[3]
		if offset != nextOffset {
			log.Fatalf("out of sync in disassembly for %s: have %s+%d, want %s+%d", name, name, nextOffset, name, offset)
		}
		hex := strings.Fields(lines[i+1])
		nextOffset = offset + len(hex)
		mcall := callRE.FindStringSubmatch(inst)
		fmt.Fprintf(testAsm, "// %s\n", lines[i])
		call := ""
		if mcall != nil {
			call = mcall[1]
		}
		if (lastCall == "task_vtimer_clear" || lastCall == "task_vtimer_set") && strings.Contains(call, "+") {
			call = "psignal_internal"
		}
		switch call {
		case "psignal_internal", "task_vtimer_clear", "task_vtimer_set":
			insts = append(insts, "call "+call)
			fmt.Fprintf(testAsm, "\tcall _%s\n", call)
		default:
			insts = append(insts, strings.Replace(inst, "\t", " ", -1))
			fmt.Fprintf(testAsm, "\t")
			for i, x := range hex {
				if i > 0 {
					fmt.Fprintf(testAsm, " ")
				}
				fmt.Fprintf(testAsm, ".byte 0x%s;", x)
			}
			fmt.Fprintf(testAsm, "\n")
		}
		hexes = append(hexes, hex)
		if call != "" {
			lastCall = call
		}
	}

	if name != "bsd_ast" {
		return
	}

	nfound := 0
	for i, inst := range insts {
		if inst != "call psignal_internal" {
			continue
		}

		// Found end. Work backward.
		state := 0
		j := i - 1
		for ; j >= 0; j-- {
			switch {
			case state == 0 && insts[j] == "call task_vtimer_clear":
				state = 1
			case state == 1 && insts[j] == "call task_vtimer_set":
				state = 2
			case state < 2 && strings.HasPrefix(insts[j], "cmp"):
				state = -1
			case state >= 2 && strings.HasPrefix(insts[j], "cmp"):
				state++
			case state < 2 && strings.HasPrefix(insts[j], "mov") && strings.HasPrefix(insts[j+1], "test"):
				state = -1
			case state >= 2 && strings.HasPrefix(insts[j], "mov") && strings.HasPrefix(insts[j+1], "test"):
				state++
			}
			if state == 4 {
				break
			}
		}
		if j < 0 {
			continue
		}
		nfound++

		// We think insts[j:i+1] and hexes[j:i+1] are worth pursuing.
		fmt.Printf("Starting at %s\n", lines[2*j])
		off0 := -1
		for k := j; k <= i; k++ {
			var buf bytes.Buffer
			hex := hexes[k]
			if hex[0] == "e8" && len(hex) == 5 {
				fmt.Fprintf(&buf, "0xe8 0x00/0x00 0x00/0x00 0x00/0x00 0x00/0x00")
			} else {
				for i, x := range hex {
					if i > 0 {
						buf.WriteString(" ")
					}
					fmt.Fprintf(&buf, "0x%s", x)
				}
			}
			m := instRE.FindStringSubmatch(lines[2*k])
			off, _ := strconv.Atoi(m[2])
			if off0 < 0 {
				off0 = off
			}
			fmt.Printf("    %-47s // %2d %s\n", buf.String(), off-off0, insts[k])
		}
		fmt.Printf("\n")
	}

	if nfound == 0 {
		fmt.Printf("no code patterns found in %s\n", name)
	}
}
