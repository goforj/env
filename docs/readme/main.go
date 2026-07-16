//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	apiStart = "<!-- api:embed:start -->"
	apiEnd   = "<!-- api:embed:end -->"
)

// main renders the reproducible documentation artifacts for this module.
func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("✔ API section updated in README.md")
}

// run replaces only the marked README API section so hand-written documentation remains untouched.
func run() error {
	root, err := findRoot()
	if err != nil {
		return err
	}

	funcs, err := parseFuncs(root)
	if err != nil {
		return err
	}

	api := renderAPI(funcs)

	readmePath := filepath.Join(root, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}

	out, err := replaceAPISection(string(data), api)
	if err != nil {
		return err
	}

	return os.WriteFile(readmePath, []byte(out), 0o644)
}

//
// ------------------------------------------------------------
// Data model
// ------------------------------------------------------------
//

// FuncDoc captures the metadata needed to render one documented function.
type FuncDoc struct {
	Name        string
	Receiver    string
	Group       string
	Behavior    string
	Fluent      string
	Description string
	Examples    []Example
}

// Example captures an executable snippet and its source location.
type Example struct {
	Label string
	Code  string
	Line  int
}

//
// ------------------------------------------------------------
// Parsing
// ------------------------------------------------------------
//

var (
	groupHeader    = regexp.MustCompile(`(?i)^\s*@group\s+(.+)$`)
	behaviorHeader = regexp.MustCompile(`(?i)^\s*@behavior\s+(.+)$`)
	fluentHeader   = regexp.MustCompile(`(?i)^\s*@fluent\s+(.+)$`)
	exampleHeader  = regexp.MustCompile(`(?i)^\s*Example:\s*(.*)$`)
)

// parseFuncs collects exported functions and methods into a receiver-qualified, source-derived model.
func parseFuncs(root string) ([]*FuncDoc, error) {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(
		fset,
		root,
		func(info os.FileInfo) bool {
			return !strings.HasSuffix(info.Name(), "_test.go")
		},
		parser.ParseComments,
	)
	if err != nil {
		return nil, err
	}

	pkgName, err := selectPackage(pkgs)
	if err != nil {
		return nil, err
	}

	pkg, ok := pkgs[pkgName]
	if !ok {
		return nil, fmt.Errorf(`package %q not found`, pkgName)
	}

	funcs := map[string]*FuncDoc{}

	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Doc == nil {
				continue
			}

			if !ast.IsExported(fn.Name.Name) {
				continue
			}

			fd := &FuncDoc{
				Name:        fn.Name.Name,
				Receiver:    receiverName(fn),
				Group:       extractGroup(fn.Doc),
				Behavior:    extractBehavior(fn.Doc),
				Fluent:      extractFluent(fn.Doc),
				Description: extractDescription(fn.Doc),
				Examples:    extractExamples(fset, fn),
			}

			key := funcIdentity(fd.Receiver, fd.Name)
			if existing, ok := funcs[key]; ok {
				existing.Examples = append(existing.Examples, fd.Examples...)
			} else {
				funcs[key] = fd
			}
		}
	}

	out := make([]*FuncDoc, 0, len(funcs))
	for _, fd := range funcs {
		sort.Slice(fd.Examples, func(i, j int) bool {
			return fd.Examples[i].Line < fd.Examples[j].Line
		})
		out = append(out, fd)
	}

	return out, nil
}

// receiverName returns the named receiver for methods so root functions cannot overwrite them.
func receiverName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}
	receiver := fn.Recv.List[0].Type
	if pointer, ok := receiver.(*ast.StarExpr); ok {
		receiver = pointer.X
	}
	if identifier, ok := receiver.(*ast.Ident); ok {
		return identifier.Name
	}
	return ""
}

// funcIdentity creates the stable key used to keep methods and functions distinct.
func funcIdentity(receiver, name string) string {
	if receiver == "" {
		return name
	}
	return receiver + "." + name
}

// funcDisplayName qualifies methods for unambiguous API links and headings.
func funcDisplayName(fn *FuncDoc) string {
	return funcIdentity(fn.Receiver, fn.Name)
}

// funcAnchor converts a display name into a stable Markdown anchor.
func funcAnchor(fn *FuncDoc) string {
	return strings.NewReplacer(".", "-", "_", "-").Replace(strings.ToLower(funcDisplayName(fn)))
}

// extractGroup honors the documentation grouping convention and places untagged APIs in Other.
func extractGroup(group *ast.CommentGroup) string {
	for _, c := range group.List {
		line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		if m := groupHeader.FindStringSubmatch(line); m != nil {
			return strings.TrimSpace(m[1])
		}
	}
	return "Other"
}

// extractBehavior normalizes behavior metadata used to classify generated API documentation.
func extractBehavior(group *ast.CommentGroup) string {
	for _, c := range group.List {
		line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		if m := behaviorHeader.FindStringSubmatch(line); m != nil {
			return strings.ToLower(strings.TrimSpace(m[1]))
		}
	}
	return ""
}

// extractFluent normalizes the fluent marker used to annotate chainable APIs.
func extractFluent(group *ast.CommentGroup) string {
	for _, c := range group.List {
		line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
		if m := fluentHeader.FindStringSubmatch(line); m != nil {
			return strings.ToLower(strings.TrimSpace(m[1]))
		}
	}
	return ""
}

// extractDescription stops before generator directives and examples so prose is not duplicated.
func extractDescription(group *ast.CommentGroup) string {
	var lines []string

	for _, c := range group.List {
		line := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))

		if exampleHeader.MatchString(line) ||
			groupHeader.MatchString(line) ||
			behaviorHeader.MatchString(line) ||
			fluentHeader.MatchString(line) {
			break
		}

		if len(lines) == 0 && line == "" {
			continue
		}

		lines = append(lines, line)
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// extractExamples retains source positions so documented cases render in declaration order.
func extractExamples(fset *token.FileSet, fn *ast.FuncDecl) []Example {
	var out []Example
	var current []string
	var label string
	var start int
	inExample := false

	flush := func() {
		if len(current) == 0 {
			return
		}

		out = append(out, Example{
			Label: label,
			Code:  strings.Join(normalizeIndent(current), "\n"),
			Line:  start,
		})

		current = nil
		label = ""
		inExample = false
	}

	for _, c := range fn.Doc.List {
		raw := strings.TrimPrefix(c.Text, "//")
		line := strings.TrimSpace(raw)

		if m := exampleHeader.FindStringSubmatch(line); m != nil {
			flush()
			inExample = true
			label = strings.TrimSpace(m[1])
			start = fset.Position(c.Slash).Line
			continue
		}

		if !inExample {
			continue
		}

		current = append(current, raw)
	}

	flush()
	return out
}

// selectPackage picks the primary package to document.
// Strategy:
//  1. If only one package exists, use it.
//  2. Prefer the non-"main" package with the most files.
//  3. Fall back to the first package alphabetically.
func selectPackage(pkgs map[string]*ast.Package) (string, error) {
	if len(pkgs) == 0 {
		return "", fmt.Errorf("no packages found")
	}

	if len(pkgs) == 1 {
		for name := range pkgs {
			return name, nil
		}
	}

	type candidate struct {
		name  string
		count int
	}

	candidates := make([]candidate, 0, len(pkgs))
	for name, pkg := range pkgs {
		candidates = append(candidates, candidate{
			name:  name,
			count: len(pkg.Files),
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].count == candidates[j].count {
			return candidates[i].name < candidates[j].name
		}
		return candidates[i].count > candidates[j].count
	})

	for _, cand := range candidates {
		if cand.name != "main" {
			return cand.name, nil
		}
	}

	return candidates[0].name, nil
}

//
// ------------------------------------------------------------
// Rendering
// ------------------------------------------------------------
//

// renderAPI groups and sorts the parsed model so README generation is reproducible.
func renderAPI(funcs []*FuncDoc) string {
	byGroup := map[string][]*FuncDoc{}

	for _, fd := range funcs {
		byGroup[fd.Group] = append(byGroup[fd.Group], fd)
	}

	groupNames := make([]string, 0, len(byGroup))
	for g := range byGroup {
		groupNames = append(groupNames, g)
	}
	sort.Strings(groupNames)

	var buf bytes.Buffer

	// ---------------- Index ----------------
	buf.WriteString("### <a id=\"api-index\"></a>API Index\n\n")
	buf.WriteString("| Group | Functions |\n")
	buf.WriteString("|------:|-----------|\n")

	for _, group := range groupNames {
		sort.Slice(byGroup[group], func(i, j int) bool {
			return funcDisplayName(byGroup[group][i]) < funcDisplayName(byGroup[group][j])
		})

		var links []string
		for _, fn := range byGroup[group] {
			links = append(links, fmt.Sprintf("[%s](#%s)", funcDisplayName(fn), funcAnchor(fn)))
		}

		buf.WriteString(fmt.Sprintf("| **%s** | %s |\n",
			group,
			strings.Join(links, " · "),
		))
	}

	buf.WriteString("\n\n")

	// ---------------- Details ----------------
	for _, group := range groupNames {
		buf.WriteString("## " + group + "\n\n")

		for _, fn := range byGroup[group] {
			anchor := funcAnchor(fn)

			header := funcDisplayName(fn)
			if fn.Fluent == "true" {
				header += " · fluent"
			}

			buf.WriteString(fmt.Sprintf("### <a id=\"%s\"></a>%s\n\n", anchor, header))

			if fn.Description != "" {
				buf.WriteString(fn.Description + "\n\n")
			}

			for _, ex := range fn.Examples {
				if ex.Label != "" {
					buf.WriteString(fmt.Sprintf("_Example: %s_\n\n", ex.Label))
				}

				buf.WriteString("```go\n")
				buf.WriteString(strings.TrimSpace(ex.Code))
				buf.WriteString("\n```\n\n")
			}
		}
	}

	return strings.TrimRight(buf.String(), "\n")
}

//
// ------------------------------------------------------------
// README replacement
// ------------------------------------------------------------
//

// replaceAPISection confines generated writes to the API markers and rejects malformed README structure.
func replaceAPISection(readme, api string) (string, error) {
	start := strings.Index(readme, apiStart)
	end := strings.Index(readme, apiEnd)

	if start == -1 || end == -1 || end < start {
		return "", fmt.Errorf("API anchors not found or malformed")
	}

	var out bytes.Buffer
	out.WriteString(readme[:start+len(apiStart)])
	out.WriteString("\n\n")
	out.WriteString(api)
	out.WriteString("\n")
	out.WriteString(readme[end:])

	return out.String(), nil
}

//
// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------
//

// findRoot anchors generation to the root module from either the root or docs directory.
func findRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	if fileExists(filepath.Join(wd, "go.mod")) {
		return wd, nil
	}
	parent := filepath.Join(wd, "..")
	if fileExists(filepath.Join(parent, "go.mod")) {
		return filepath.Clean(parent), nil
	}
	return "", fmt.Errorf("could not find project root")
}

// fileExists lets root discovery ignore candidate paths that are not present.
func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// normalizeIndent removes shared documentation padding without changing relative code indentation.
func normalizeIndent(lines []string) []string {
	min := -1

	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		n := len(l) - len(strings.TrimLeft(l, " \t"))
		if min == -1 || n < min {
			min = n
		}
	}

	if min <= 0 {
		return lines
	}

	out := make([]string, len(lines))
	for i, l := range lines {
		if len(l) >= min {
			out[i] = l[min:]
		} else {
			out[i] = strings.TrimLeft(l, " \t")
		}
	}

	return out
}
