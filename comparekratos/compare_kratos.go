// Package comparekratos provides comparison functions between Kratos demo projects
// Compares source and fork projects, generates markdown difference reports
// Supports comparing code and generating tree structures
//
// comparekratos 提供 Kratos 演示项目之间的比较功能
// 比较源项目和 fork 项目，生成 markdown 差异报告
// 支持代码比较和生成目录树结构
package comparekratos

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/yylego/must"
	"github.com/yylego/osexec"
	"github.com/yylego/osexistpath/osmustexist"
	"github.com/yylego/printgo"
	"github.com/yylego/rese"
	"github.com/yylego/tint"
	"github.com/yylego/zaplog"
)

// ComparePath uses diff command to compare two paths and output results
// Ignores go.mod, go.sum, and bin differences
//
// ComparePath 使用 diff 命令比较两个路径的差异并输出结果
// 忽略 go.mod、go.sum 和 bin 的差异
func ComparePath(path0 string, path1 string) {
	path0 = osmustexist.ROOT(path0)
	path1 = osmustexist.ROOT(path1)
	zaplog.SUG.Debugln("path0:", path0)
	zaplog.SUG.Debugln("path1:", path1)
	output := rese.A1(osexec.NewCommandConfig().WithDebugMode(osexec.SHOW_COMMAND).WithExpectExit(1, "DIFFERENCES FOUND").
		Exec(
			"diff",
			"-ruN",
			"--exclude=go.mod",
			"--exclude=go.sum",
			"--exclude=bin",
			path0,
			path1,
		))
	if len(output) == 0 {
		tint.GREEN.ShowMessage("SAME")
	} else {
		tint.AMBER.ShowMessage("⬇⬇⬇")
		zaplog.SUG.Debugln(string(output))
		tint.AMBER.ShowMessage("⬆⬆⬆")
	}
}

// ShowReadableChanges shows formatted readable changes between two paths
// Red indicates deleted lines, green indicates added lines
//
// ShowReadableChanges 显示格式化的易读变更结果
// 红色显示删除的代码行，绿色显示新增的代码行
func ShowReadableChanges(path0, path1 string) {
	path0 = osmustexist.ROOT(path0)
	path1 = osmustexist.ROOT(path1)
	output := rese.A1(osexec.NewCommandConfig().WithExpectExit(1, "DIFFERENCES FOUND").
		Exec(
			"diff",
			"-ruN",
			"--exclude=go.mod",
			"--exclude=go.sum",
			"--exclude=bin",
			path0,
			path1,
		))

	if len(output) == 0 {
		tint.GREEN.ShowMessage("✅ NO CHANGES")
		return
	}
	tint.AMBER.ShowMessage("📋 FOUND DIFFERENCES")

	var sourcePath string
	var adds, cuts []string

	printFile := func() {
		if sourcePath != "" && (len(adds) > 0 || len(cuts) > 0) {
			zaplog.SUG.Debugln(tint.CYAN.Sprintf("📄 %s (+%d -%d)", sourcePath, len(adds), len(cuts)))
			for _, line := range cuts {
				zaplog.SUG.Debugln(tint.RED.Sprint("  - " + line))
			}
			for _, line := range adds {
				zaplog.SUG.Debugln(tint.GREEN.Sprint("  + " + line))
			}
		}
	}

	for _, line := range strings.Split(string(output), "\n") {
		switch {
		case strings.HasPrefix(line, "diff -ruN"):
			printFile()
			sourcePath, adds, cuts = "", nil, nil

		case strings.HasPrefix(line, "---"):
			// skip

		case strings.HasPrefix(line, "+++"):
			if parts := strings.Fields(line); len(parts) >= 2 {
				if strings.Contains(parts[1], path1+"/") {
					sourcePath = strings.TrimPrefix(parts[1], path1+"/")
				} else {
					sourcePath = filepath.Base(parts[1])
				}
			}

		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			adds = append(adds, line[1:])

		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			cuts = append(cuts, line[1:])
		}
	}

	printFile()
}

// GenerateChangesFile generates markdown file of code differences
// Outputs diff in markdown format with syntax highlighting
//
// GenerateChangesFile 生成代码差异的 markdown 文件
// 以 markdown 格式输出 diff，支持语法高亮
func GenerateChangesFile(path0, path1, outputPath string) {
	path0 = osmustexist.ROOT(path0)
	path1 = osmustexist.ROOT(path1)
	output := rese.A1(osexec.NewCommandConfig().WithExpectExit(1, "DIFFERENCES FOUND").
		Exec(
			"diff",
			"-ruN",
			"-U3",
			"--exclude=go.mod",
			"--exclude=go.sum",
			"--exclude=bin",
			path0,
			path1,
		))

	if len(output) == 0 {
		content := "# Changes\n\n✅ NO CHANGES\n"
		must.Done(os.WriteFile(outputPath, []byte(content), 0644))
		zaplog.SUG.Debugln("Generated", outputPath, "(no changes)")
		return
	}

	var sourcePath string
	var diffLines []string
	var addCount, cutCount int

	ptx := printgo.NewPTX()
	ptx.Println("# Changes")
	ptx.Println()
	ptx.Println("Code differences compared to source project.")
	ptx.Println()

	processFile := func() {
		if sourcePath != "" && (addCount > 0 || cutCount > 0) {
			ptx.Printf("## %s (+%d -%d)\n\n", sourcePath, addCount, cutCount)
			ptx.Println("```diff")
			for _, line := range diffLines {
				ptx.Println(line)
			}
			ptx.Println("```")
			ptx.Println()
		}
	}

	for _, line := range strings.Split(string(output), "\n") {
		switch {
		case strings.HasPrefix(line, "diff -ruN"):
			processFile()
			sourcePath, diffLines, addCount, cutCount = "", nil, 0, 0

		case strings.HasPrefix(line, "---"):
			// skip

		case strings.HasPrefix(line, "+++"):
			if parts := strings.Fields(line); len(parts) >= 2 {
				if strings.Contains(parts[1], path1+"/") {
					sourcePath = strings.TrimPrefix(parts[1], path1+"/")
				} else {
					sourcePath = filepath.Base(parts[1])
				}
			}

		case strings.HasPrefix(line, "@@"):
			diffLines = append(diffLines, line)

		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			diffLines = append(diffLines, line)
			addCount++

		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			diffLines = append(diffLines, line)
			cutCount++

		case strings.HasPrefix(line, " "):
			diffLines = append(diffLines, line)
		}
	}

	processFile()

	must.Done(os.WriteFile(outputPath, ptx.Bytes(), 0644))
	zaplog.SUG.Debugln("Generated", outputPath, "with differences")
}

// GenerateTreeChanges generates tree structure of sibling projects
// Lists DIRs except excluded ones, outputs to markdown
//
// GenerateTreeChanges 生成兄弟项目的目录树结构
// 列出除了排除目录之外的所有目录，输出到 markdown
func GenerateTreeChanges(root string, excludeNames []string, outputPath string) {
	root = osmustexist.ROOT(root)
	zaplog.SUG.Debugln(root)

	var matchNames []string
	for _, item := range rese.A1(os.ReadDir(root)) {
		if item.IsDir() {
			name := item.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			if slices.Contains(excludeNames, name) {
				continue
			}
			zaplog.SUG.Debugln("match:", name)
			matchNames = append(matchNames, name)
		}
	}
	zaplog.SUG.Debugln("match:", matchNames)

	if len(matchNames) == 0 {
		content := "# Changes\n\n✅ NO CHANGES\n"
		must.Done(os.WriteFile(outputPath, []byte(content), 0644))
		return
	}

	ptx := printgo.NewPTX()
	ptx.Println("# Changes")
	ptx.Println()

	ptx.Println("## Overview")
	ptx.Println()
	ptx.Println("Sibling projects:")
	ptx.Println()
	for _, name := range matchNames {
		ptx.Fprintf("- [%s](#%s)", name, name)
		ptx.Println()
	}
	ptx.Println()

	ptx.Println("## Project Structures")
	ptx.Println()
	for idx, name := range matchNames {
		if idx > 0 {
			ptx.Println("---")
			ptx.Println()
		}

		ptx.Fprintf("### %s", name)
		ptx.Println()
		ptx.Println()
		ptx.Fprintf("**Location**: [%s](../%s)", name, name)
		ptx.Println()
		ptx.Println()
		ptx.Println("```bash")
		ptx.Fprintf("cd %s && tree --noreport", name)
		ptx.Println()
		ptx.Println("```")
		ptx.Println()

		subRoot := filepath.Join(root, name)
		treeOutput := rese.A1(osexec.ExecInPath(subRoot, "tree", "--noreport", "--charset=ascii", "--gitignore", "-I", "node_modules|.git|bin|.idea|.vscode"))

		ptx.Println("```")
		ptx.Write(treeOutput)
		ptx.Println("```")
		ptx.Println()
	}
	must.Done(os.WriteFile(outputPath, ptx.Bytes(), 0644))
}
