package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ConversationArchive represents an archived conversation
type ConversationArchive struct {
	Date         string
	Project      string
	Topic        string
	Participants []string
	Summary      string
	Content      string
	Decisions    []string
	Tags         []string
}

func main() {
	// Command line flags
	project := flag.String("project", "", "Project name (e.g., llm-security-bas)")
	topic := flag.String("topic", "", "Conversation topic")
	summary := flag.String("summary", "", "Brief summary of the conversation")
	participants := flag.String("participants", "Claude,User", "Comma-separated list of participants")
	decisions := flag.String("decisions", "", "Comma-separated list of key decisions")
	tags := flag.String("tags", "", "Comma-separated list of tags")
	inputFile := flag.String("input", "", "Input file containing conversation (reads from stdin if not specified)")
	brainDir := flag.String("dir", "", "AI Brain directory (auto-detect if not specified)")
	list := flag.Bool("list", false, "List recent conversations")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *help {
		printUsage()
		return
	}

	// Auto-detect brain directory
	if *brainDir == "" {
		*brainDir = findBrainDir()
		if *brainDir == "" {
			fmt.Println("错误: 无法找到ai-brain目录，请使用 -dir 指定")
			os.Exit(1)
		}
	}

	// List mode
	if *list {
		listConversations(*brainDir, *project)
		return
	}

	// Validate required parameters
	if *project == "" {
		fmt.Println("错误: 请指定项目名称 -project <name>")
		fmt.Println("可用项目:")
		listProjects(*brainDir)
		os.Exit(1)
	}

	if *topic == "" {
		*topic = "对话记录"
	}

	// Read conversation content
	var content string
	var err error

	if *inputFile != "" {
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Printf("错误: 无法读取文件 %s: %v\n", *inputFile, err)
			os.Exit(1)
		}
		content = string(data)
	} else {
		// Read from stdin
		content, err = readFromStdin()
		if err != nil {
			fmt.Printf("错误: 无法从stdin读取: %v\n", err)
			os.Exit(1)
		}
	}

	if strings.TrimSpace(content) == "" {
		fmt.Println("错误: 对话内容为空")
		os.Exit(1)
	}

	// Parse participants
	participantList := []string{}
	if *participants != "" {
		for _, p := range strings.Split(*participants, ",") {
			participantList = append(participantList, strings.TrimSpace(p))
		}
	}

	// Parse decisions
	decisionList := []string{}
	if *decisions != "" {
		for _, d := range strings.Split(*decisions, ",") {
			decisionList = append(decisionList, strings.TrimSpace(d))
		}
	}

	// Parse tags
	tagList := []string{}
	if *tags != "" {
		for _, t := range strings.Split(*tags, ",") {
			tagList = append(tagList, strings.TrimSpace(t))
		}
	}

	// Create archive
	archive := ConversationArchive{
		Date:         time.Now().Format("2006-01-02"),
		Project:      *project,
		Topic:        *topic,
		Participants: participantList,
		Summary:      *summary,
		Content:      content,
		Decisions:    decisionList,
		Tags:         tagList,
	}

	// Save archive
	outputPath, err := saveArchive(*brainDir, archive)
	if err != nil {
		fmt.Printf("错误: 保存失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    对话归档成功                                ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("📁 项目: %s\n", archive.Project)
	fmt.Printf("📝 主题: %s\n", archive.Topic)
	fmt.Printf("📅 日期: %s\n", archive.Date)
	fmt.Printf("📄 文件: %s\n", outputPath)
	fmt.Println()
}

func printUsage() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║                    AI大脑 对话归档工具                         ║
║                     brain-archive v1.0                        ║
╚═══════════════════════════════════════════════════════════════╝

用法:
  brain-archive [选项]                 归档对话
  brain-archive -list [-project xxx]   列出历史对话

选项:
  -project <name>     项目名称 (必填)
  -topic <topic>      对话主题 (默认: "对话记录")
  -summary <text>     对话摘要
  -participants <p>   参与者，逗号分隔 (默认: Claude,User)
  -decisions <d>      关键决策，逗号分隔
  -tags <t>           标签，逗号分隔
  -input <file>       输入文件 (默认: 从stdin读取)
  -dir <path>         AI Brain目录 (默认: 自动检测)
  -list               列出历史对话
  -h                  显示帮助

示例:
  # 从stdin归档对话
  echo "对话内容..." | brain-archive -project llm-security-bas -topic "需求讨论"

  # 从文件归档
  brain-archive -project llm-security-bas -topic "架构设计" -input conversation.txt

  # 带完整元数据
  brain-archive -project llm-security-bas \
    -topic "API设计讨论" \
    -summary "讨论了API接口规范" \
    -decisions "使用RESTful,采用JWT认证" \
    -tags "API,设计" \
    -input conv.txt

  # 列出项目对话
  brain-archive -list -project llm-security-bas

归档目录:
  projects/{project}/_conversations/YYYY-MM-DD_主题.md
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

func readFromStdin() (string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	// Increase buffer size for large inputs
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func listProjects(brainDir string) {
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

func listConversations(brainDir, project string) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    对话归档记录                                ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	projectsDir := filepath.Join(brainDir, "projects")

	if project != "" {
		// List conversations for specific project
		listProjectConversations(projectsDir, project)
	} else {
		// List all projects and their conversations
		entries, err := os.ReadDir(projectsDir)
		if err != nil {
			fmt.Printf("错误: 无法读取项目目录: %v\n", err)
			return
		}

		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				listProjectConversations(projectsDir, entry.Name())
			}
		}
	}
}

func listProjectConversations(projectsDir, project string) {
	convDir := filepath.Join(projectsDir, project, "_conversations")
	entries, err := os.ReadDir(convDir)
	if err != nil {
		return
	}

	fmt.Printf("📁 %s\n", project)

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			info, _ := entry.Info()
			modTime := ""
			if info != nil {
				modTime = info.ModTime().Format("2006-01-02 15:04")
			}
			fmt.Printf("   📄 %s  (%s)\n", entry.Name(), modTime)
		}
	}
	fmt.Println()
}

func saveArchive(brainDir string, archive ConversationArchive) (string, error) {
	// Create conversations directory
	convDir := filepath.Join(brainDir, "projects", archive.Project, "_conversations")
	if err := os.MkdirAll(convDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// Generate filename
	safeTopicName := sanitizeFilename(archive.Topic)
	filename := fmt.Sprintf("%s_%s.md", archive.Date, safeTopicName)
	outputPath := filepath.Join(convDir, filename)

	// Check if file exists, add suffix if needed
	counter := 1
	for {
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			break
		}
		filename = fmt.Sprintf("%s_%s_%d.md", archive.Date, safeTopicName, counter)
		outputPath = filepath.Join(convDir, filename)
		counter++
	}

	// Generate markdown content
	content := generateMarkdown(archive)

	// Write file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	relPath, _ := filepath.Rel(brainDir, outputPath)
	return relPath, nil
}

func sanitizeFilename(name string) string {
	// Remove or replace invalid characters
	reg := regexp.MustCompile(`[<>:"/\\|?*\s]+`)
	safe := reg.ReplaceAllString(name, "_")

	// Limit length
	if len(safe) > 50 {
		safe = safe[:50]
	}

	return strings.Trim(safe, "_")
}

func generateMarkdown(archive ConversationArchive) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", archive.Topic))
	sb.WriteString(fmt.Sprintf("date: %s\n", archive.Date))
	sb.WriteString(fmt.Sprintf("created: %s\n", time.Now().Format("2006-01-02 15:04:05")))

	if len(archive.Participants) > 0 {
		sb.WriteString(fmt.Sprintf("participants: [%s]\n", strings.Join(archive.Participants, ", ")))
	}

	if len(archive.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(archive.Tags, ", ")))
	}

	sb.WriteString("status: archived\n")
	sb.WriteString("---\n\n")

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", archive.Topic))

	// Metadata
	sb.WriteString("## 基本信息\n\n")
	sb.WriteString(fmt.Sprintf("- **日期**: %s\n", archive.Date))
	sb.WriteString(fmt.Sprintf("- **项目**: %s\n", archive.Project))
	if len(archive.Participants) > 0 {
		sb.WriteString(fmt.Sprintf("- **参与者**: %s\n", strings.Join(archive.Participants, ", ")))
	}
	sb.WriteString("\n")

	// Summary
	if archive.Summary != "" {
		sb.WriteString("## 摘要\n\n")
		sb.WriteString(archive.Summary)
		sb.WriteString("\n\n")
	}

	// Key decisions
	if len(archive.Decisions) > 0 {
		sb.WriteString("## 关键决策\n\n")
		for _, d := range archive.Decisions {
			sb.WriteString(fmt.Sprintf("- %s\n", d))
		}
		sb.WriteString("\n")
	}

	// Conversation content
	sb.WriteString("## 对话内容\n\n")
	sb.WriteString(archive.Content)
	sb.WriteString("\n")

	return sb.String()
}
