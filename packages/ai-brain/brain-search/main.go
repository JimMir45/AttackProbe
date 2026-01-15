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

// SearchResult represents a single search match
type SearchResult struct {
	FilePath    string
	RelPath     string
	LineNumber  int
	LineContent string
	Context     []string // lines before and after
	Score       int      // relevance score
}

// FileResult represents all matches in a file
type FileResult struct {
	FilePath  string
	RelPath   string
	Title     string
	Matches   []SearchResult
	Score     int
	UpdatedAt string
}

// Config holds search configuration
type Config struct {
	RootDir     string
	Query       string
	ContextLines int
	MaxResults  int
	ShowContext bool
	CaseSensitive bool
}

func main() {
	// Parse command line flags
	rootDir := flag.String("dir", "", "Root directory to search (default: auto-detect ai-brain)")
	contextLines := flag.Int("context", 2, "Number of context lines to show")
	maxResults := flag.Int("max", 20, "Maximum number of results")
	showContext := flag.Bool("c", true, "Show context lines")
	caseSensitive := flag.Bool("s", false, "Case sensitive search")
	help := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *help || flag.NArg() == 0 {
		printUsage()
		return
	}

	query := strings.Join(flag.Args(), " ")

	// Auto-detect ai-brain directory
	if *rootDir == "" {
		*rootDir = findBrainDir()
		if *rootDir == "" {
			fmt.Println("Error: Cannot find ai-brain directory. Use -dir to specify.")
			os.Exit(1)
		}
	}

	config := Config{
		RootDir:      *rootDir,
		Query:        query,
		ContextLines: *contextLines,
		MaxResults:   *maxResults,
		ShowContext:  *showContext,
		CaseSensitive: *caseSensitive,
	}

	results := search(config)
	printResults(results, config)
}

func printUsage() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║                    AI大脑 全文检索工具                         ║
║                      brain-search v1.0                        ║
╚═══════════════════════════════════════════════════════════════╝

用法: brain-search [选项] <关键词>

选项:
  -dir <path>    指定搜索目录 (默认: 自动检测ai-brain目录)
  -context <n>   显示上下文行数 (默认: 2)
  -max <n>       最大结果数 (默认: 20)
  -c             显示上下文 (默认: true)
  -s             区分大小写 (默认: false)
  -h             显示帮助

示例:
  brain-search 白皮书
  brain-search "技术指标"
  brain-search -max 10 部署
  brain-search -context 3 OWASP

支持:
  • 中英文关键词搜索
  • 多关键词搜索 (空格分隔, AND逻辑)
  • 结果按相关度排序
`)
}

func findBrainDir() string {
	// Try common locations
	candidates := []string{
		"./ai-brain",
		"../ai-brain",
		"../../ai-brain",
		"/home/vackbot/vackbas/ai-brain",
	}

	// Also try from current working directory
	cwd, _ := os.Getwd()
	candidates = append(candidates, filepath.Join(cwd, "ai-brain"))

	for _, dir := range candidates {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if info, err := os.Stat(absDir); err == nil && info.IsDir() {
			// Verify it's an ai-brain directory by checking for projects folder
			if _, err := os.Stat(filepath.Join(absDir, "projects")); err == nil {
				return absDir
			}
		}
	}
	return ""
}

func search(config Config) []FileResult {
	var results []FileResult
	keywords := parseKeywords(config.Query, config.CaseSensitive)

	err := filepath.Walk(config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories and non-markdown files
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// Search in file
		fileResult := searchFile(path, config.RootDir, keywords, config)
		if fileResult != nil && len(fileResult.Matches) > 0 {
			results = append(results, *fileResult)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	// Sort by score (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit results
	if len(results) > config.MaxResults {
		results = results[:config.MaxResults]
	}

	return results
}

func parseKeywords(query string, caseSensitive bool) []string {
	keywords := strings.Fields(query)
	if !caseSensitive {
		for i, kw := range keywords {
			keywords[i] = strings.ToLower(kw)
		}
	}
	return keywords
}

func searchFile(filePath, rootDir string, keywords []string, config Config) *FileResult {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	relPath, _ := filepath.Rel(rootDir, filePath)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil
	}

	result := &FileResult{
		FilePath: filePath,
		RelPath:  relPath,
		Title:    extractTitle(lines),
		Matches:  []SearchResult{},
	}

	// Extract updated date from frontmatter
	result.UpdatedAt = extractDate(lines)

	// Search each line
	for lineNum, line := range lines {
		searchLine := line
		if !config.CaseSensitive {
			searchLine = strings.ToLower(line)
		}

		// Check if all keywords match
		allMatch := true
		matchCount := 0
		for _, kw := range keywords {
			if strings.Contains(searchLine, kw) {
				matchCount++
			} else {
				allMatch = false
			}
		}

		// For single keyword, just need to match
		// For multiple keywords, require all to match in the line OR file
		if len(keywords) == 1 && matchCount > 0 || allMatch {
			match := SearchResult{
				FilePath:    filePath,
				RelPath:     relPath,
				LineNumber:  lineNum + 1,
				LineContent: line,
				Score:       matchCount * 10,
			}

			// Add context
			if config.ShowContext {
				match.Context = getContext(lines, lineNum, config.ContextLines)
			}

			result.Matches = append(result.Matches, match)
			result.Score += match.Score
		}
	}

	// Boost score for title matches
	titleLower := strings.ToLower(result.Title)
	for _, kw := range keywords {
		if strings.Contains(titleLower, kw) {
			result.Score += 50
		}
	}

	// Boost score for recent files
	if result.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02", result.UpdatedAt); err == nil {
			daysSince := time.Since(t).Hours() / 24
			if daysSince < 7 {
				result.Score += 20
			} else if daysSince < 30 {
				result.Score += 10
			}
		}
	}

	return result
}

func extractTitle(lines []string) string {
	inFrontmatter := false
	for _, line := range lines {
		if line == "---" {
			if inFrontmatter {
				break // End of frontmatter
			}
			inFrontmatter = true
			continue
		}

		if inFrontmatter {
			if strings.HasPrefix(line, "title:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			}
		} else {
			// Look for first H1
			if strings.HasPrefix(line, "# ") {
				return strings.TrimPrefix(line, "# ")
			}
		}
	}
	return ""
}

func extractDate(lines []string) string {
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
			if strings.HasPrefix(line, "updated:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "updated:"))
			}
		}
	}
	return ""
}

func getContext(lines []string, lineNum int, contextLines int) []string {
	var context []string
	start := lineNum - contextLines
	end := lineNum + contextLines + 1

	if start < 0 {
		start = 0
	}
	if end > len(lines) {
		end = len(lines)
	}

	for i := start; i < end; i++ {
		prefix := "  "
		if i == lineNum {
			prefix = "> "
		}
		context = append(context, fmt.Sprintf("%s%d: %s", prefix, i+1, lines[i]))
	}

	return context
}

func printResults(results []FileResult, config Config) {
	if len(results) == 0 {
		fmt.Println("\n未找到匹配结果")
		fmt.Printf("搜索关键词: %s\n", config.Query)
		fmt.Printf("搜索目录: %s\n", config.RootDir)
		return
	}

	// Count total matches
	totalMatches := 0
	for _, r := range results {
		totalMatches += len(r.Matches)
	}

	fmt.Println()
	fmt.Printf("╔═══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  搜索结果: \"%s\"                                       \n", config.Query)
	fmt.Printf("║  找到 %d 个文件, %d 处匹配                                \n", len(results), totalMatches)
	fmt.Printf("╚═══════════════════════════════════════════════════════════════╝\n")

	for i, fileResult := range results {
		fmt.Println()
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📄 [%d] %s\n", i+1, fileResult.RelPath)
		if fileResult.Title != "" {
			fmt.Printf("   标题: %s\n", fileResult.Title)
		}
		fmt.Printf("   匹配: %d处  相关度: %d\n", len(fileResult.Matches), fileResult.Score)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

		// Show matches (limit to first 3 per file)
		showCount := len(fileResult.Matches)
		if showCount > 3 {
			showCount = 3
		}

		for j := 0; j < showCount; j++ {
			match := fileResult.Matches[j]
			fmt.Println()

			if config.ShowContext && len(match.Context) > 0 {
				for _, ctx := range match.Context {
					// Highlight the matching line
					if strings.HasPrefix(ctx, "> ") {
						fmt.Printf("   \033[33m%s\033[0m\n", ctx)
					} else {
						fmt.Printf("   %s\n", ctx)
					}
				}
			} else {
				fmt.Printf("   行 %d: %s\n", match.LineNumber, highlightKeywords(match.LineContent, config.Query))
			}
		}

		if len(fileResult.Matches) > 3 {
			fmt.Printf("\n   ... 还有 %d 处匹配\n", len(fileResult.Matches)-3)
		}
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("提示: 使用 -max <n> 调整结果数量, -context <n> 调整上下文行数")
}

func highlightKeywords(line, query string) string {
	keywords := strings.Fields(query)
	result := line

	for _, kw := range keywords {
		// Case insensitive replace with highlighting
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(kw))
		result = re.ReplaceAllStringFunc(result, func(match string) string {
			return fmt.Sprintf("\033[33m%s\033[0m", match)
		})
	}

	return result
}
