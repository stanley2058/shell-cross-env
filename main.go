package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BashSession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	reader *bufio.Reader
}

func NewBashSession() (*BashSession, error) {
	cmd := exec.Command("/bin/bash", "--norc", "--noprofile")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &BashSession{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		reader: bufio.NewReader(stdout),
	}, nil
}

func (b *BashSession) RunCommand(command string) (string, error) {
	// Add a unique marker to know when command ends
	marker := "END_OF_COMMAND"
	fullCmd := fmt.Sprintf("%s; echo %s\n", command, marker)

	if _, err := b.stdin.Write([]byte(fullCmd)); err != nil {
		return "", err
	}

	var output strings.Builder
	for {
		line, err := b.reader.ReadString('\n')
		if err != nil {
			return output.String(), err
		}
		if strings.Contains(line, marker) {
			break
		}
		output.WriteString(line)
	}

	return output.String(), nil
}

func (b *BashSession) Close() {
	b.stdin.Write([]byte("exit\n"))
	b.cmd.Wait()
}

func main() {
	args := os.Args
	binary := args[0]
	filename := filepath.Base(binary)

	if len(args) < 5 || args[1] != "--to" || args[3] != "source" {
		fmt.Printf("Usage: %s --to <fish|bash> source <file1> <file2>\n", filename)
		os.Exit(1)
	}
	targetShell := args[2]
	if targetShell != "fish" && targetShell != "bash" {
		fmt.Printf("Usage: %s --to <fish|bash> source <file1> <file2>\n", filename)
		os.Exit(1)
	}
	sourceFiles := args[4:]

	bash, err := NewBashSession()
	if err != nil {
		panic(err)
	}
	defer bash.Close()

	baseEnv, _ := bash.RunCommand("env")
	baseAlias, _ := bash.RunCommand("alias")

	baseEnvMap := make(map[string]string)
	for line := range strings.SplitSeq(baseEnv, "\n") {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			baseEnvMap[parts[0]] = parts[1]
		}
	}
	baseAliasMap := make(map[string]string)
	for line := range strings.SplitSeq(baseAlias, "\n") {
		if strings.HasPrefix(line, "alias ") {
			parts := strings.SplitN(strings.Replace(line, "alias ", "", 1), "=", 2)
			baseAliasMap[parts[0]] = parts[1]
		}
	}

	for _, file := range sourceFiles {
		if _, err := bash.RunCommand(fmt.Sprintf("source %s", file)); err != nil {
			panic(err)
		}
	}

	newEnv, _ := bash.RunCommand("env")
	newAlias, _ := bash.RunCommand("alias")

	addedEnv := make(map[string]string)
	for line := range strings.SplitSeq(newEnv, "\n") {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if val, ok := baseEnvMap[parts[0]]; ok && val == parts[1] {
				continue
			}
			addedEnv[parts[0]] = parts[1]
		}
	}

	addedAlias := make(map[string]string)
	for line := range strings.SplitSeq(newAlias, "\n") {
		if strings.HasPrefix(line, "alias ") {
			parts := strings.SplitN(strings.Replace(line, "alias ", "", 1), "=", 2)
			if val, ok := baseAliasMap[parts[0]]; ok && val == parts[1] {
				continue
			}
			addedAlias[parts[0]] = parts[1]
		}
	}

	if targetShell == "fish" {
		for key, val := range addedEnv {
			fmt.Printf("set -gx %s %s\n", key, val)
		}
		for key, val := range addedAlias {
			fmt.Printf("alias %s=%s\n", key, val)
		}
	} else {
		for key, val := range addedEnv {
			fmt.Printf("export %s=%s\n", key, val)
		}
		for key, val := range addedAlias {
			fmt.Printf("alias %s=%s\n", key, val)
		}
	}
}
