package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	BrainDir      string `json:"brain_dir"`
	OllamaURL     string `json:"ollama_url"`
	EmbedModel    string `json:"embed_model"`
	ChatModel     string `json:"chat_model"`
	IndexFile     string `json:"index_file"`
	ChunkSize     int    `json:"chunk_size"`
	ChunkOverlap  int    `json:"chunk_overlap"`
	TopK          int    `json:"top_k"`
	Temperature   float64 `json:"temperature"`
}

// Document represents a markdown document
type Document struct {
	Path      string `json:"path"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
}

// Chunk represents a document chunk with embedding
type Chunk struct {
	ID        string    `json:"id"`
	DocPath   string    `json:"doc_path"`
	DocTitle  string    `json:"doc_title"`
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding"`
	StartLine int       `json:"start_line"`
}

// VectorIndex holds all chunks and metadata
type VectorIndex struct {
	Version   string    `json:"version"`
	CreatedAt string    `json:"created_at"`
	Config    Config    `json:"config"`
	Chunks    []Chunk   `json:"chunks"`
	DocCount  int       `json:"doc_count"`
}

// SearchResult represents a similarity search result
type SearchResult struct {
	Chunk      Chunk   `json:"chunk"`
	Score      float64 `json:"score"`
}

// OllamaEmbedRequest is the request format for Ollama embeddings
type OllamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbedResponse is the response format for Ollama embeddings
type OllamaEmbedResponse struct {
	Embedding []float64 `json:"embedding"`
}

// OllamaChatRequest is the request format for Ollama chat
type OllamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ChatMessage   `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaChatResponse is the response format for Ollama chat
type OllamaChatResponse struct {
	Message ChatMessage `json:"message"`
}

var defaultConfig = Config{
	BrainDir:     "",
	OllamaURL:    "http://10.10.10.123:11434",
	EmbedModel:   "nomic-embed-text",
	ChatModel:    "qwen2:0.5b",
	IndexFile:    "",
	ChunkSize:    500,
	ChunkOverlap: 50,
	TopK:         5,
	Temperature:  0.7,
}

func main() {
	// Subcommands
	indexCmd := flag.NewFlagSet("index", flag.ExitOnError)
	askCmd := flag.NewFlagSet("ask", flag.ExitOnError)

	// Global flags
	var configFile string
	var showHelp bool

	flag.StringVar(&configFile, "config", "", "Config file path")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.Parse()

	if showHelp || len(os.Args) < 2 {
		printUsage()
		return
	}

	// Load config
	config := loadConfig(configFile)

	switch os.Args[1] {
	case "index":
		// Index subcommand flags
		indexCmd.StringVar(&config.BrainDir, "dir", config.BrainDir, "Brain directory to index")
		indexCmd.StringVar(&config.OllamaURL, "ollama", config.OllamaURL, "Ollama API URL")
		indexCmd.StringVar(&config.EmbedModel, "model", config.EmbedModel, "Embedding model")
		indexCmd.Parse(os.Args[2:])

		if config.BrainDir == "" {
			config.BrainDir = findBrainDir()
		}
		if config.IndexFile == "" {
			config.IndexFile = filepath.Join(config.BrainDir, "_system", "index.json")
		}

		runIndex(config)

	case "ask", "query", "q":
		// Ask subcommand flags
		askCmd.StringVar(&config.BrainDir, "dir", config.BrainDir, "Brain directory")
		askCmd.StringVar(&config.OllamaURL, "ollama", config.OllamaURL, "Ollama API URL")
		askCmd.StringVar(&config.ChatModel, "model", config.ChatModel, "Chat model")
		askCmd.IntVar(&config.TopK, "k", config.TopK, "Number of chunks to retrieve")
		askCmd.Parse(os.Args[2:])

		if config.BrainDir == "" {
			config.BrainDir = findBrainDir()
		}
		if config.IndexFile == "" {
			config.IndexFile = filepath.Join(config.BrainDir, "_system", "index.json")
		}

		question := strings.Join(askCmd.Args(), " ")
		if question == "" {
			fmt.Println("请输入问题")
			return
		}

		runAsk(config, question)

	case "status":
		if config.BrainDir == "" {
			config.BrainDir = findBrainDir()
		}
		if config.IndexFile == "" {
			config.IndexFile = filepath.Join(config.BrainDir, "_system", "index.json")
		}
		runStatus(config)

	default:
		// Treat as question if no subcommand
		question := strings.Join(os.Args[1:], " ")
		if config.BrainDir == "" {
			config.BrainDir = findBrainDir()
		}
		if config.IndexFile == "" {
			config.IndexFile = filepath.Join(config.BrainDir, "_system", "index.json")
		}
		runAsk(config, question)
	}
}

func printUsage() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║                    AI大脑 智能问答工具                         ║
║                      brain-ask v1.0                           ║
╚═══════════════════════════════════════════════════════════════╝

用法:
  brain-ask index [选项]       构建向量索引
  brain-ask ask "问题"         提问并获取答案
  brain-ask "问题"             直接提问 (等同于 ask)
  brain-ask status             查看索引状态

索引选项:
  -dir <path>      指定大脑目录 (默认: 自动检测)
  -ollama <url>    Ollama API地址 (默认: http://10.10.10.123:11434)
  -model <name>    Embedding模型 (默认: nomic-embed-text)

问答选项:
  -dir <path>      指定大脑目录
  -ollama <url>    Ollama API地址
  -model <name>    对话模型 (默认: qwen2:0.5b)
  -k <n>           检索文档数 (默认: 5)

示例:
  # 首次使用，构建索引
  brain-ask index

  # 提问
  brain-ask "白皮书在哪里？"
  brain-ask "部署需要什么配置？"
  brain-ask "测试用例有多少个？"

  # 使用不同模型
  brain-ask -model qwen2:7b "解释一下越狱攻击"

环境要求:
  - Ollama 服务运行中 (需要 embedding 和 chat 模型)
  - 推荐模型: nomic-embed-text (embedding), qwen2 (chat)
`)
}

func loadConfig(configFile string) Config {
	config := defaultConfig

	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err == nil {
			json.Unmarshal(data, &config)
		}
	}

	return config
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
	// Check if we're inside ai-brain
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

// ==================== Index Functions ====================

func runIndex(config Config) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    构建向量索引                                ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("📁 大脑目录: %s\n", config.BrainDir)
	fmt.Printf("🔗 Ollama: %s\n", config.OllamaURL)
	fmt.Printf("🧠 Embedding模型: %s\n", config.EmbedModel)
	fmt.Println()

	// Test Ollama connection
	fmt.Print("检查Ollama连接... ")
	if !testOllamaConnection(config.OllamaURL) {
		fmt.Println("❌ 失败")
		fmt.Println("\n错误: 无法连接到Ollama服务")
		fmt.Printf("请确保Ollama服务正在运行: %s\n", config.OllamaURL)
		fmt.Println("\n提示: 启动Ollama服务后重试")
		os.Exit(1)
	}
	fmt.Println("✅ 成功")

	// Check embedding model
	fmt.Printf("检查Embedding模型 (%s)... ", config.EmbedModel)
	if !testEmbeddingModel(config) {
		fmt.Println("❌ 失败")
		fmt.Printf("\n错误: Embedding模型 '%s' 不可用\n", config.EmbedModel)
		fmt.Println("\n请先拉取模型:")
		fmt.Printf("  ollama pull %s\n", config.EmbedModel)
		os.Exit(1)
	}
	fmt.Println("✅ 成功")
	fmt.Println()

	// Scan documents
	fmt.Println("扫描文档...")
	docs := scanDocuments(config.BrainDir)
	fmt.Printf("找到 %d 个Markdown文档\n\n", len(docs))

	if len(docs) == 0 {
		fmt.Println("未找到文档，退出")
		return
	}

	// Chunk documents
	fmt.Println("文档分块...")
	var allChunks []Chunk
	for _, doc := range docs {
		chunks := chunkDocument(doc, config)
		allChunks = append(allChunks, chunks...)
	}
	fmt.Printf("生成 %d 个文本块\n\n", len(allChunks))

	// Generate embeddings
	fmt.Println("生成向量嵌入...")
	startTime := time.Now()
	for i := range allChunks {
		if i%10 == 0 || i == len(allChunks)-1 {
			fmt.Printf("\r  进度: %d/%d (%.1f%%)", i+1, len(allChunks), float64(i+1)/float64(len(allChunks))*100)
		}
		embedding, err := getEmbedding(config, allChunks[i].Content)
		if err != nil {
			fmt.Printf("\n警告: 块 %d 嵌入失败: %v\n", i, err)
			continue
		}
		allChunks[i].Embedding = embedding
	}
	fmt.Printf("\n  耗时: %.1f秒\n\n", time.Since(startTime).Seconds())

	// Filter out chunks without embeddings
	var validChunks []Chunk
	for _, chunk := range allChunks {
		if len(chunk.Embedding) > 0 {
			validChunks = append(validChunks, chunk)
		}
	}

	// Create index
	index := VectorIndex{
		Version:   "1.0",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Config:    config,
		Chunks:    validChunks,
		DocCount:  len(docs),
	}

	// Save index
	fmt.Print("保存索引... ")
	if err := saveIndex(config.IndexFile, index); err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 成功")

	// Summary
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("索引构建完成!")
	fmt.Printf("  文档数: %d\n", len(docs))
	fmt.Printf("  文本块: %d\n", len(validChunks))
	fmt.Printf("  索引文件: %s\n", config.IndexFile)
	fmt.Printf("  文件大小: %.2f MB\n", float64(getFileSize(config.IndexFile))/1024/1024)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("现在可以使用 brain-ask \"问题\" 进行问答")
}

func testOllamaConnection(url string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func testEmbeddingModel(config Config) bool {
	_, err := getEmbedding(config, "test")
	return err == nil
}

func scanDocuments(brainDir string) []Document {
	var docs []Document

	filepath.Walk(brainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(brainDir, path)
		doc := Document{
			Path:      relPath,
			Title:     extractTitle(string(content)),
			Content:   string(content),
			UpdatedAt: info.ModTime().Format("2006-01-02"),
		}
		docs = append(docs, doc)
		return nil
	})

	return docs
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
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
			if strings.HasPrefix(line, "title:") {
				return strings.TrimSpace(strings.TrimPrefix(line, "title:"))
			}
		} else {
			if strings.HasPrefix(line, "# ") {
				return strings.TrimPrefix(line, "# ")
			}
		}
	}
	return ""
}

func chunkDocument(doc Document, config Config) []Chunk {
	var chunks []Chunk

	// Remove frontmatter
	content := doc.Content
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			content = parts[2]
		}
	}

	// Split by sections (headers)
	sections := splitBySections(content)

	chunkID := 0
	for _, section := range sections {
		// If section is too long, split further
		if len(section) > config.ChunkSize*2 {
			subChunks := splitBySize(section, config.ChunkSize, config.ChunkOverlap)
			for _, subChunk := range subChunks {
				if len(strings.TrimSpace(subChunk)) < 50 {
					continue
				}
				chunks = append(chunks, Chunk{
					ID:       fmt.Sprintf("%s_%d", doc.Path, chunkID),
					DocPath:  doc.Path,
					DocTitle: doc.Title,
					Content:  strings.TrimSpace(subChunk),
				})
				chunkID++
			}
		} else if len(strings.TrimSpace(section)) >= 50 {
			chunks = append(chunks, Chunk{
				ID:       fmt.Sprintf("%s_%d", doc.Path, chunkID),
				DocPath:  doc.Path,
				DocTitle: doc.Title,
				Content:  strings.TrimSpace(section),
			})
			chunkID++
		}
	}

	return chunks
}

func splitBySections(content string) []string {
	var sections []string
	var currentSection strings.Builder

	headerRegex := regexp.MustCompile(`^#{1,4}\s+`)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if headerRegex.MatchString(line) {
			if currentSection.Len() > 0 {
				sections = append(sections, currentSection.String())
				currentSection.Reset()
			}
		}
		currentSection.WriteString(line)
		currentSection.WriteString("\n")
	}

	if currentSection.Len() > 0 {
		sections = append(sections, currentSection.String())
	}

	return sections
}

func splitBySize(text string, chunkSize, overlap int) []string {
	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunk := string(runes[i:end])
		chunks = append(chunks, chunk)
		if end == len(runes) {
			break
		}
	}

	return chunks
}

func getEmbedding(config Config, text string) ([]float64, error) {
	reqBody := OllamaEmbedRequest{
		Model:  config.EmbedModel,
		Prompt: text,
	}

	jsonData, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(config.OllamaURL+"/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result OllamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embedding, nil
}

func saveIndex(path string, index VectorIndex) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// ==================== Ask Functions ====================

func runAsk(config Config, question string) {
	// Load index
	index, err := loadIndex(config.IndexFile)
	if err != nil {
		fmt.Println("错误: 无法加载索引")
		fmt.Printf("索引文件: %s\n", config.IndexFile)
		fmt.Println("\n请先运行 brain-ask index 构建索引")
		os.Exit(1)
	}

	// Check Ollama
	if !testOllamaConnection(config.OllamaURL) {
		fmt.Println("错误: 无法连接到Ollama服务")
		fmt.Printf("请确保Ollama服务正在运行: %s\n", config.OllamaURL)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("💬 问题: %s\n", question)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Get question embedding
	fmt.Print("🔍 检索相关文档... ")
	questionEmbed, err := getEmbedding(config, question)
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		os.Exit(1)
	}

	// Search similar chunks
	results := searchSimilar(index.Chunks, questionEmbed, config.TopK)
	fmt.Printf("找到 %d 个相关片段\n", len(results))

	// Show retrieved chunks
	fmt.Println()
	fmt.Println("📚 参考文档:")
	for i, r := range results {
		fmt.Printf("  [%d] %s (相关度: %.2f)\n", i+1, r.Chunk.DocPath, r.Score)
	}

	// Build context
	var contextParts []string
	for _, r := range results {
		contextParts = append(contextParts, fmt.Sprintf("来源: %s\n%s", r.Chunk.DocPath, r.Chunk.Content))
	}
	context := strings.Join(contextParts, "\n\n---\n\n")

	// Generate answer
	fmt.Println()
	fmt.Print("🤖 生成回答... ")

	answer, err := generateAnswer(config, question, context)
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("完成")

	// Print answer
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📝 回答:")
	fmt.Println()
	fmt.Println(answer)
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func loadIndex(path string) (*VectorIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var index VectorIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	return &index, nil
}

func searchSimilar(chunks []Chunk, queryEmbed []float64, topK int) []SearchResult {
	var results []SearchResult

	for _, chunk := range chunks {
		if len(chunk.Embedding) == 0 {
			continue
		}
		score := cosineSimilarity(queryEmbed, chunk.Embedding)
		results = append(results, SearchResult{
			Chunk: chunk,
			Score: score,
		})
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > topK {
		results = results[:topK]
	}

	return results
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func generateAnswer(config Config, question, context string) (string, error) {
	systemPrompt := `你是AI大脑的智能助手，负责基于知识库内容回答用户问题。

规则:
1. 只基于提供的上下文信息回答，不要编造信息
2. 如果上下文中没有相关信息，诚实地说"根据知识库内容，我没有找到相关信息"
3. 回答要简洁、准确、有条理
4. 如果涉及多个文档，可以综合信息回答
5. 适当引用来源文档路径`

	userPrompt := fmt.Sprintf(`基于以下知识库内容回答问题。

【知识库内容】
%s

【问题】
%s

请基于上述内容回答问题:`, context, question)

	reqBody := OllamaChatRequest{
		Model: config.ChatModel,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Stream: false,
		Options: map[string]interface{}{
			"temperature": config.Temperature,
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Post(config.OllamaURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(body))
	}

	var result OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}

// ==================== Status Functions ====================

func runStatus(config Config) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    AI大脑索引状态                              ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("📁 大脑目录: %s\n", config.BrainDir)
	fmt.Printf("📄 索引文件: %s\n", config.IndexFile)
	fmt.Println()

	// Check index file
	index, err := loadIndex(config.IndexFile)
	if err != nil {
		fmt.Println("❌ 索引不存在或无法读取")
		fmt.Println("\n请运行 brain-ask index 构建索引")
		return
	}

	fmt.Println("✅ 索引状态: 正常")
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("版本: %s\n", index.Version)
	fmt.Printf("创建时间: %s\n", index.CreatedAt)
	fmt.Printf("文档数量: %d\n", index.DocCount)
	fmt.Printf("文本块数: %d\n", len(index.Chunks))
	fmt.Printf("Embedding模型: %s\n", index.Config.EmbedModel)
	fmt.Printf("文件大小: %.2f MB\n", float64(getFileSize(config.IndexFile))/1024/1024)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Check Ollama
	fmt.Println()
	fmt.Printf("🔗 Ollama服务: %s\n", config.OllamaURL)
	if testOllamaConnection(config.OllamaURL) {
		fmt.Println("   状态: ✅ 在线")
	} else {
		fmt.Println("   状态: ❌ 离线")
	}

	// Show document list
	fmt.Println()
	fmt.Println("📚 已索引文档:")
	docMap := make(map[string]int)
	for _, chunk := range index.Chunks {
		docMap[chunk.DocPath]++
	}

	var docs []string
	for doc := range docMap {
		docs = append(docs, doc)
	}
	sort.Strings(docs)

	for _, doc := range docs {
		fmt.Printf("   %s (%d块)\n", doc, docMap[doc])
	}
}

// ==================== Interactive Mode ====================

func runInteractive(config Config) {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    AI大脑 交互模式                             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("输入问题进行问答，输入 'quit' 或 'exit' 退出")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("❓ 请输入问题: ")
		if !scanner.Scan() {
			break
		}
		question := strings.TrimSpace(scanner.Text())

		if question == "" {
			continue
		}
		if question == "quit" || question == "exit" {
			fmt.Println("再见!")
			break
		}

		runAsk(config, question)
		fmt.Println()
	}
}
