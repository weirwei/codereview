package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

// ShellExec Exec Shell Command
func ShellExec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	Debugf("command exec:%s", cmd.String())
	result, err := cmd.CombinedOutput()
	if err != nil {
		Errorf("Shell Exec failed: command: %s, err: %s, output: %s", cmd.String(), err.Error(), result)
		return "", err
	}
	Debugf("command exec result:%s", strings.ReplaceAll(string(result), "\n", "\\n"))
	return string(result), err
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

// WriteFile
func WriteFile(content, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		Errorf("Failed to create file:%v", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		Errorf("Failed to write file:", err)
		return
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
