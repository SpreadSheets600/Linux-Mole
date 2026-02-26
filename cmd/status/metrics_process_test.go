package main

import "testing"

func TestParseTopProcesses_DarwinFormat(t *testing.T) {
	out := ` %CPU %MEM COMMAND
42.5  3.1 /Applications/Foo.app/Contents/MacOS/Foo
25.0  1.2 /usr/bin/python3
`

	procs := parseTopProcesses(out, 5, true)
	if len(procs) != 2 {
		t.Fatalf("len(procs) = %d, want 2", len(procs))
	}
	if procs[0].Name != "Foo" {
		t.Fatalf("procs[0].Name = %q, want %q", procs[0].Name, "Foo")
	}
	if procs[1].Name != "python3" {
		t.Fatalf("procs[1].Name = %q, want %q", procs[1].Name, "python3")
	}
}

func TestParseTopProcesses_LinuxNoHeader(t *testing.T) {
	out := ` 9.1  2.0 systemd
 8.0  1.5 node
 invalid line
 4.2  0.5 dockerd
`

	procs := parseTopProcesses(out, 2, false)
	if len(procs) != 2 {
		t.Fatalf("len(procs) = %d, want 2", len(procs))
	}
	if procs[0].Name != "systemd" || procs[1].Name != "node" {
		t.Fatalf("unexpected names: %+v", procs)
	}
}

func TestProcessPSCommands_LinuxIncludesFallback(t *testing.T) {
	cmds, ok := processPSCommands("linux")
	if !ok {
		t.Fatal("expected linux ps commands to be supported")
	}
	if len(cmds) < 2 {
		t.Fatalf("len(cmds) = %d, want at least 2", len(cmds))
	}
}
