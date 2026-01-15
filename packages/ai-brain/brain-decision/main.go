package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Decision represents a decision record
type Decision struct {
	ID           string
	Title        string
	Date         string
	Project      string
	Phase        string // 00-立项, 01-设计, etc.
	Background   string
	Options      []Option
	Decision     string
	Rationale    string
	Impact       string
	Participants []string
	Tags         []string
	Status       string // proposed, approved, superseded, deprecated
	SupersededBy string
}

// Option represents a decision option
type Option struct {
	Name        string
	Description string
	Pros        []string
	Cons        []string
}

func main() {
	// Subcommands
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	// Global flags
	var brainDir string
	var help bool

	flag.StringVar(&brainDir, "dir", "", "AI Brain directory")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.Parse()

	if help || len(os.Args) < 2 {
		printUsage()
		return
	}

	// Auto-detect brain directory
	if brainDir == "" {
		brainDir = findBrainDir()
		if brainDir == "" {
			fmt.Println("错误: 无法找到ai-brain目录，请使用 -dir 指定")
			os.Exit(1)
		}
	}

	switch os.Args[1] {
	case "new", "create", "add":
		// New decision flags
		var project, title, phase string
		var interactive bool
		newCmd.StringVar(&project, "project", "", "Project name (required)")
		newCmd.StringVar(&title, "title", "", "Decision title (required)")
		newCmd.StringVar(&phase, "phase", "00-立项", "Project phase (00-立项, 01-设计, etc.)")
		newCmd.BoolVar(&interactive, "i", false, "Interactive mode")
		newCmd.Parse(os.Args[2:])

		if project == "" {
			fmt.Println("错误: 请指定项目名称 -project <name>")
			listProjects(brainDir)
			os.Exit(1)
		}

		if title == "" && !interactive {
			fmt.Println("错误: 请指定决策标题 -title <title> 或使用交互模式 -i")
			os.Exit(1)
		}

		if interactive {
			runInteractiveNew(brainDir, project, phase)
		} else {
			runNew(brainDir, project, title, phase)
		}

	case "list", "ls":
		var project, status string
		listCmd.StringVar(&project, "project", "", "Filter by project")
		listCmd.StringVar(&status, "status", "", "Filter by status (proposed, approved, superseded, deprecated)")
		listCmd.Parse(os.Args[2:])

		runList(brainDir, project, status)

	case "show", "view":
		showCmd.Parse(os.Args[2:])
		if len(showCmd.Args()) == 0 {
			fmt.Println("错误: 请指定决策文件路径")
			os.Exit(1)
		}
		runShow(brainDir, showCmd.Args()[0])

	case "template":
		printTemplate()

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║                    AI大脑 决策记录工具                         ║
║                     brain-decision v1.0                       ║
╚═══════════════════════════════════════════════════════════════╝

用法:
  brain-decision new [选项]        创建新决策记录
  brain-decision list [选项]       列出决策记录
  brain-decision show <path>       查看决策详情
  brain-decision template          显示决策模板

创建选项:
  -project <name>    项目名称 (必填)
  -title <title>     决策标题 (必填，除非使用 -i)
  -phase <phase>     项目阶段 (默认: 00-立项)
  -i                 交互模式

列表选项:
  -project <name>    按项目筛选
  -status <status>   按状态筛选 (proposed, approved, superseded, deprecated)

项目阶段:
  00-立项    立项阶段决策
  01-设计    设计阶段决策
  02-开发    开发阶段决策
  03-测试    测试阶段决策
  04-部署    部署阶段决策
  05-运营    运营阶段决策

示例:
  # 交互式创建决策
  brain-decision new -project llm-security-bas -i

  # 快速创建决策
  brain-decision new -project llm-security-bas -title "技术选型决策" -phase 01-设计

  # 列出所有决策
  brain-decision list

  # 列出特定项目的决策
  brain-decision list -project llm-security-bas

  # 显示决策模板
  brain-decision template

决策状态:
  proposed    - 提议中，待讨论
  approved    - 已批准，执行中
  superseded  - 已被新决策替代
  deprecated  - 已废弃
`)
}

func findBrainDir() string {
	candidates := []string{
		"./ai-brain",
		"../ai-brain",
		"../../ai-brain",
		"../../../ai-brain",
		"/home/vackbot/vackbas/ai-brain",
	}

	cwd, _ := os.Getwd()
	if strings.Contains(cwd, "ai-brain") {
		parts := strings.Split(cwd, "ai-brain")
		candidates = append([]string{parts[0] + "ai-brain"}, candidates...)
	}

	for _, dir := range candidates {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if info, err := os.Stat(absDir); err == nil && info.IsDir() {
			if _, err := os.Stat(filepath.Join(absDir, "projects")); err == nil {
				return absDir
			}
		}
	}
	return ""
}

func listProjects(brainDir string) {
	fmt.Println("可用项目:")
	projectsDir := filepath.Join(brainDir, "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			fmt.Printf("  - %s\n", entry.Name())
		}
	}
}

func runInteractiveNew(brainDir, project, phase string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    创建决策记录 (交互模式)                      ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("📁 项目: %s\n", project)
	fmt.Printf("📂 阶段: %s\n", phase)
	fmt.Println()

	// Title
	fmt.Print("📝 决策标题: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	if title == "" {
		fmt.Println("错误: 标题不能为空")
		os.Exit(1)
	}

	// Background
	fmt.Println("\n📋 背景 (为什么需要这个决策? 输入空行结束):")
	background := readMultiLine(reader)

	// Options
	fmt.Println("\n📊 选项 (每行一个选项，格式: 选项名: 描述，输入空行结束):")
	optionsText := readMultiLine(reader)
	options := parseOptions(optionsText)

	// Decision
	fmt.Print("\n✅ 决策 (选择了哪个选项): ")
	decision, _ := reader.ReadString('\n')
	decision = strings.TrimSpace(decision)

	// Rationale
	fmt.Println("\n💡 理由 (为什么选择这个? 输入空行结束):")
	rationale := readMultiLine(reader)

	// Impact
	fmt.Println("\n⚡ 影响 (这个决策会带来什么影响? 输入空行结束):")
	impact := readMultiLine(reader)

	// Participants
	fmt.Print("\n👥 参与者 (逗号分隔): ")
	participantsStr, _ := reader.ReadString('\n')
	participants := parseList(participantsStr)

	// Tags
	fmt.Print("\n🏷️ 标签 (逗号分隔): ")
	tagsStr, _ := reader.ReadString('\n')
	tags := parseList(tagsStr)

	// Status
	fmt.Print("\n📌 状态 (proposed/approved, 默认 approved): ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)
	if status == "" {
		status = "approved"
	}

	// Create decision
	dec := Decision{
		ID:           generateID(),
		Title:        title,
		Date:         time.Now().Format("2006-01-02"),
		Project:      project,
		Phase:        phase,
		Background:   background,
		Options:      options,
		Decision:     decision,
		Rationale:    rationale,
		Impact:       impact,
		Participants: participants,
		Tags:         tags,
		Status:       status,
	}

	// Save
	outputPath, err := saveDecision(brainDir, dec)
	if err != nil {
		fmt.Printf("错误: 保存失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ 决策记录创建成功!")
	fmt.Printf("📄 文件: %s\n", outputPath)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func runNew(brainDir, project, title, phase string) {
	dec := Decision{
		ID:           generateID(),
		Title:        title,
		Date:         time.Now().Format("2006-01-02"),
		Project:      project,
		Phase:        phase,
		Background:   "[请填写决策背景]",
		Options:      []Option{{Name: "选项A", Description: "[描述]"}, {Name: "选项B", Description: "[描述]"}},
		Decision:     "[请填写最终决策]",
		Rationale:    "[请填写选择理由]",
		Impact:       "[请填写决策影响]",
		Participants: []string{"[参与者]"},
		Tags:         []string{},
		Status:       "proposed",
	}

	outputPath, err := saveDecision(brainDir, dec)
	if err != nil {
		fmt.Printf("错误: 保存失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    决策记录创建成功                            ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("📁 项目: %s\n", project)
	fmt.Printf("📂 阶段: %s\n", phase)
	fmt.Printf("📝 标题: %s\n", title)
	fmt.Printf("📄 文件: %s\n", outputPath)
	fmt.Println()
	fmt.Println("请编辑文件填写完整的决策内容。")
}

func runList(brainDir, project, status string) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    决策记录列表                                ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	projectsDir := filepath.Join(brainDir, "projects")
	var projects []string

	if project != "" {
		projects = []string{project}
	} else {
		entries, err := os.ReadDir(projectsDir)
		if err != nil {
			fmt.Printf("错误: 无法读取项目目录: %v\n", err)
			return
		}
		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				projects = append(projects, entry.Name())
			}
		}
	}

	totalCount := 0
	for _, proj := range projects {
		decisions := findDecisions(filepath.Join(projectsDir, proj))
		if len(decisions) == 0 {
			continue
		}

		// Filter by status
		if status != "" {
			var filtered []DecisionInfo
			for _, d := range decisions {
				if d.Status == status {
					filtered = append(filtered, d)
				}
			}
			decisions = filtered
		}

		if len(decisions) == 0 {
			continue
		}

		fmt.Printf("📁 %s\n", proj)
		for _, d := range decisions {
			statusIcon := getStatusIcon(d.Status)
			fmt.Printf("   %s [%s] %s (%s)\n", statusIcon, d.ID, d.Title, d.Date)
			fmt.Printf("      📄 %s\n", d.Path)
		}
		fmt.Println()
		totalCount += len(decisions)
	}

	if totalCount == 0 {
		fmt.Println("未找到决策记录")
		if project != "" {
			fmt.Printf("\n提示: 使用 brain-decision new -project %s -i 创建新决策\n", project)
		}
	} else {
		fmt.Printf("共找到 %d 条决策记录\n", totalCount)
	}
}

func runShow(brainDir, path string) {
	// Try to find the file
	fullPath := path
	if !filepath.IsAbs(path) {
		fullPath = filepath.Join(brainDir, path)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("错误: 无法读取文件: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(content))
}

func printTemplate() {
	fmt.Println(`
# 决策记录模板

---
id: ADR-XXXX
title: 决策标题
date: YYYY-MM-DD
status: proposed | approved | superseded | deprecated
participants: [参与者1, 参与者2]
tags: [标签1, 标签2]
superseded_by: (如果被替代，填写新决策ID)
---

# 决策: {决策标题}

## 背景

为什么需要这个决策? 描述当前面临的问题或机会。

## 选项

### 选项A: {选项名称}
{选项描述}

**优点**:
- 优点1
- 优点2

**缺点**:
- 缺点1
- 缺点2

### 选项B: {选项名称}
{选项描述}

**优点**:
- 优点1
- 优点2

**缺点**:
- 缺点1
- 缺点2

## 决策

我们选择 **选项X**。

## 理由

为什么选择这个选项? 详细说明决策依据。

## 影响

这个决策会带来什么影响?

- 正面影响:
  - ...
- 负面影响/风险:
  - ...
- 需要的后续行动:
  - ...

## 参与者

- @人员1 - 角色/职责
- @人员2 - 角色/职责
`)
}

// Helper types and functions

type DecisionInfo struct {
	ID     string
	Title  string
	Date   string
	Status string
	Path   string
}

func findDecisions(projectDir string) []DecisionInfo {
	var decisions []DecisionInfo

	// Walk through all phase directories
	phases := []string{"00-立项", "01-设计", "02-开发", "03-测试", "04-部署", "05-运营"}
	for _, phase := range phases {
		decisionDir := filepath.Join(projectDir, phase, "决策记录")
		entries, err := os.ReadDir(decisionDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				info := parseDecisionFile(filepath.Join(decisionDir, entry.Name()))
				if info.ID != "" {
					relPath, _ := filepath.Rel(filepath.Dir(filepath.Dir(filepath.Dir(decisionDir))), filepath.Join(decisionDir, entry.Name()))
					info.Path = relPath
					decisions = append(decisions, info)
				}
			}
		}
	}

	// Sort by date (newest first)
	sort.Slice(decisions, func(i, j int) bool {
		return decisions[i].Date > decisions[j].Date
	})

	return decisions
}

func parseDecisionFile(path string) DecisionInfo {
	content, err := os.ReadFile(path)
	if err != nil {
		return DecisionInfo{}
	}

	info := DecisionInfo{}
	lines := strings.Split(string(content), "\n")
	inFrontmatter := false

	for _, line := range lines {
		if line == "---" {
			if inFrontmatter {
				break
			}
			inFrontmatter = true
			continue
		}

		if inFrontmatter {
			if strings.HasPrefix(line, "id:") {
				info.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
			} else if strings.HasPrefix(line, "title:") {
				info.Title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			} else if strings.HasPrefix(line, "date:") {
				info.Date = strings.TrimSpace(strings.TrimPrefix(line, "date:"))
			} else if strings.HasPrefix(line, "status:") {
				info.Status = strings.TrimSpace(strings.TrimPrefix(line, "status:"))
			}
		}
	}

	// If no ID found, try to extract from filename
	if info.ID == "" {
		base := filepath.Base(path)
		if strings.HasPrefix(base, "ADR-") {
			parts := strings.SplitN(base, "_", 2)
			if len(parts) > 0 {
				info.ID = strings.TrimSuffix(parts[0], ".md")
			}
		}
	}

	return info
}

func getStatusIcon(status string) string {
	switch status {
	case "proposed":
		return "📋"
	case "approved":
		return "✅"
	case "superseded":
		return "🔄"
	case "deprecated":
		return "❌"
	default:
		return "📄"
	}
}

func readMultiLine(reader *bufio.Reader) string {
	var lines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func parseList(s string) []string {
	var result []string
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func parseOptions(text string) []Option {
	var options []Option
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		opt := Option{Name: strings.TrimSpace(parts[0])}
		if len(parts) > 1 {
			opt.Description = strings.TrimSpace(parts[1])
		}
		options = append(options, opt)
	}
	return options
}

func generateID() string {
	return fmt.Sprintf("ADR-%s", time.Now().Format("20060102-150405"))
}

func sanitizeFilename(name string) string {
	reg := regexp.MustCompile(`[<>:"/\\|?*\s]+`)
	safe := reg.ReplaceAllString(name, "_")
	if len(safe) > 50 {
		safe = safe[:50]
	}
	return strings.Trim(safe, "_")
}

func saveDecision(brainDir string, dec Decision) (string, error) {
	// Create decision directory
	decisionDir := filepath.Join(brainDir, "projects", dec.Project, dec.Phase, "决策记录")
	if err := os.MkdirAll(decisionDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// Generate filename
	safeTitle := sanitizeFilename(dec.Title)
	filename := fmt.Sprintf("%s_%s.md", dec.ID, safeTitle)
	outputPath := filepath.Join(decisionDir, filename)

	// Generate markdown content
	content := generateDecisionMarkdown(dec)

	// Write file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	relPath, _ := filepath.Rel(brainDir, outputPath)
	return relPath, nil
}

func generateDecisionMarkdown(dec Decision) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("id: %s\n", dec.ID))
	sb.WriteString(fmt.Sprintf("title: %s\n", dec.Title))
	sb.WriteString(fmt.Sprintf("date: %s\n", dec.Date))
	sb.WriteString(fmt.Sprintf("status: %s\n", dec.Status))

	if len(dec.Participants) > 0 {
		sb.WriteString(fmt.Sprintf("participants: [%s]\n", strings.Join(dec.Participants, ", ")))
	}

	if len(dec.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(dec.Tags, ", ")))
	}

	if dec.SupersededBy != "" {
		sb.WriteString(fmt.Sprintf("superseded_by: %s\n", dec.SupersededBy))
	}

	sb.WriteString("---\n\n")

	// Title
	sb.WriteString(fmt.Sprintf("# 决策: %s\n\n", dec.Title))

	// Background
	sb.WriteString("## 背景\n\n")
	sb.WriteString(dec.Background)
	sb.WriteString("\n\n")

	// Options
	sb.WriteString("## 选项\n\n")
	for i, opt := range dec.Options {
		sb.WriteString(fmt.Sprintf("### 选项%d: %s\n\n", i+1, opt.Name))
		if opt.Description != "" {
			sb.WriteString(opt.Description)
			sb.WriteString("\n\n")
		}
	}

	// Decision
	sb.WriteString("## 决策\n\n")
	sb.WriteString(dec.Decision)
	sb.WriteString("\n\n")

	// Rationale
	if dec.Rationale != "" {
		sb.WriteString("## 理由\n\n")
		sb.WriteString(dec.Rationale)
		sb.WriteString("\n\n")
	}

	// Impact
	sb.WriteString("## 影响\n\n")
	sb.WriteString(dec.Impact)
	sb.WriteString("\n\n")

	// Participants
	sb.WriteString("## 参与者\n\n")
	for _, p := range dec.Participants {
		sb.WriteString(fmt.Sprintf("- @%s\n", p))
	}
	sb.WriteString("\n")

	return sb.String()
}
