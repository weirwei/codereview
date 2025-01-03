package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/weirwei/codereview/log"
	"golang.org/x/net/html"
)

// WriteFile
func WriteFile(content, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Errorf("Failed to create file:%v", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		log.Errorf("Failed to write file:", err)
		return
	}
}

// MatchFileLanguage Match code language based on file extension
func MatchFileLanguage(fileName string) string {
	switch filepath.Ext(fileName) {
	case ".go":
		return "go"
	default:
		return ""
	}
}

// EstimateTokens Estimate tokens
func EstimateTokens(s string) int {
	tokenCount := 0
	runes := []rune(s)

	for i := 0; i < len(runes); {
		if unicode.IsSpace(runes[i]) {
			i++
			continue
		}

		if unicode.IsPunct(runes[i]) || unicode.IsSymbol(runes[i]) {
			tokenCount++
			i++
		} else if unicode.Is(unicode.Han, runes[i]) {
			tokenCount++
			i++
		} else {
			end := i + 4
			if end > len(runes) {
				end = len(runes)
			}
			tokenCount++
			i = end
		}
	}

	return tokenCount
}

// ExtractHtmlTagContent Extract html tag content
func ExtractHtmlTagContent(tag string, htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var outputs string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				outputs += extractText(c)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return outputs, nil
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text.WriteString(extractText(c))
	}
	return text.String()
}

func RmDuplication[T comparable](l []T) (res []T) {
	m := make(map[T]bool)
	for _, v := range l {
		if !m[v] {
			res = append(res, v)
			m[v] = true
		}
	}
	return
}

func ToJson(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(data)
}

func GetLangByFilepath(filepath string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath[strings.LastIndex(filepath, "."):], "."))
	switch ext {
	case "go":
		return "go"
	case "py", "python":
		return "python"
	case "js", "jsx":
		return "javascript"
	case "ts", "tsx":
		return "typescript"
	case "java":
		return "java"
	case "c":
		return "c"
	case "cpp", "cc", "cxx":
		return "cpp"
	case "cs":
		return "csharp"
	case "rb":
		return "ruby"
	case "php":
		return "php"
	case "swift":
		return "swift"
	case "kt":
		return "kotlin"
	case "rs":
		return "rust"
	default:
		return ""
	}
}
