package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CompletionManager struct {
	config    *Config
	tokenGen  *TokenGenerator
	searchMgr *SearchManager
}

func NewCompletionManager() *CompletionManager {
	config, err := LoadConfig()
	if err != nil {
		config = DefaultConfig()
	}

	return &CompletionManager{
		config:    config,
		tokenGen:  NewTokenGenerator(),
		searchMgr: NewSearchManager(),
	}
}

func (cm *CompletionManager) GetShareCodeSuggestions(partial string) []string {
	if len(partial) < 2 {
		return []string{}
	}

	suggestions := make([]string, 0)

	adjectives := []string{
		"blue", "red", "green", "yellow", "purple", "orange", "pink", "brown",
		"black", "white", "gray", "silver", "gold", "navy", "teal", "coral",
		"lime", "indigo", "violet", "maroon", "olive", "cyan", "magenta",
		"fast", "slow", "big", "small", "tall", "short", "wide", "narrow",
		"bright", "dark", "light", "heavy", "soft", "hard", "smooth", "rough",
		"warm", "cool", "fresh", "old", "new", "young", "ancient", "modern",
		"quiet", "loud", "calm", "wild", "gentle", "fierce", "brave", "shy",
	}
	nouns := []string{
		"whale", "dolphin", "shark", "turtle", "seahorse", "octopus", "jellyfish",
		"crab", "lobster", "starfish", "clam", "oyster", "mussel", "coral",
		"eagle", "hawk", "owl", "falcon", "raven", "crow", "sparrow", "robin",
		"lion", "tiger", "bear", "wolf", "fox", "deer", "rabbit", "squirrel",
		"elephant", "giraffe", "zebra", "rhino", "hippo", "camel", "llama",
		"dragon", "phoenix", "unicorn", "griffin", "pegasus", "centaur",
		"mountain", "river", "ocean", "forest", "desert", "island", "valley",
		"castle", "tower", "bridge", "temple", "palace", "cottage", "cabin",
		"diamond", "ruby", "emerald", "sapphire", "pearl", "crystal", "gem",
	}

	partial = strings.ToLower(partial)
	parts := strings.Split(partial, "-")

	switch len(parts) {
	case 1:
		for _, adj := range adjectives {
			if strings.HasPrefix(strings.ToLower(adj), partial) {
				suggestions = append(suggestions, adj)
			}
		}
	case 2:
		if len(parts[1]) > 0 {
			for _, noun := range nouns {
				if strings.HasPrefix(strings.ToLower(noun), parts[1]) {
					suggestions = append(suggestions, fmt.Sprintf("%s-%s", parts[0], noun))
				}
			}
		} else {
			for _, noun := range nouns {
				suggestions = append(suggestions, fmt.Sprintf("%s-%s", parts[0], noun))
			}
		}
	case 3:
		if len(parts[2]) > 0 {
			for i := 1; i <= 999; i++ {
				numStr := fmt.Sprintf("%d", i)
				if strings.HasPrefix(numStr, parts[2]) {
					suggestions = append(suggestions, fmt.Sprintf("%s-%s-%s", parts[0], parts[1], numStr))
				}
			}
		} else {
			for i := 1; i <= 10; i++ {
				suggestions = append(suggestions, fmt.Sprintf("%s-%s-%d", parts[0], parts[1], i))
			}
		}
	}

	if len(suggestions) > 10 {
		suggestions = suggestions[:10]
	}

	return suggestions
}

func (cm *CompletionManager) GetFileSuggestions(partial string) []string {
	if partial == "" {
		partial = "."
	}

	dir := filepath.Dir(partial)
	base := filepath.Base(partial)

	if dir == "." {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	suggestions := make([]string, 0)
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(base)) {
			fullPath := filepath.Join(dir, name)
			if entry.IsDir() {
				suggestions = append(suggestions, fullPath+"/")
			} else if strings.HasSuffix(name, ".env") || strings.HasSuffix(name, ".config") {
				suggestions = append(suggestions, fullPath)
			}
		}
	}

	return suggestions
}

func (cm *CompletionManager) GetTemplateSuggestions(partial string) []string {
	templateManager := NewTemplateManager()
	templates := templateManager.ListTemplates()

	suggestions := make([]string, 0)
	partial = strings.ToLower(partial)

	for _, template := range templates {
		if strings.HasPrefix(strings.ToLower(template.Name), partial) ||
			strings.Contains(strings.ToLower(template.Description), partial) ||
			strings.Contains(strings.ToLower(template.Category), partial) {
			suggestions = append(suggestions, template.Name)
		}
	}

	return suggestions
}

func (cm *CompletionManager) GetCategorySuggestions(partial string) []string {
	categories := []string{
		"database", "api", "cloud", "container", "cache", "email", "logging", "general",
		"production", "development", "testing", "staging", "environment", "configuration",
	}

	suggestions := make([]string, 0)
	partial = strings.ToLower(partial)

	for _, category := range categories {
		if strings.HasPrefix(strings.ToLower(category), partial) {
			suggestions = append(suggestions, category)
		}
	}

	return suggestions
}

func (cm *CompletionManager) GetCommandSuggestions(partial string) []string {
	commands := []string{
		"share", "get", "template", "bulk", "search", "suggest",
		"list", "show", "create", "search",
	}

	suggestions := make([]string, 0)
	partial = strings.ToLower(partial)

	for _, cmd := range commands {
		if strings.HasPrefix(strings.ToLower(cmd), partial) {
			suggestions = append(suggestions, cmd)
		}
	}

	return suggestions
}

func (cm *CompletionManager) GetFlagSuggestions(command, partial string) []string {
	flagMap := map[string][]string{
		"share":    {"--expiry", "--readonly", "--output"},
		"get":      {"--output"},
		"template": {"--output"},
		"bulk":     {"--expiry", "--readonly", "--prefix", "--group-by"},
		"search":   {"--files", "--sensitive", "--categories", "--exact", "--case-sensitive", "--regex", "--output"},
		"suggest":  {"--files"},
	}

	flags, exists := flagMap[command]
	if !exists {
		return []string{}
	}

	suggestions := make([]string, 0)
	partial = strings.ToLower(partial)

	for _, flag := range flags {
		if strings.HasPrefix(strings.ToLower(flag), partial) {
			suggestions = append(suggestions, flag)
		}
	}

	return suggestions
}

func (cm *CompletionManager) GenerateBashCompletion() string {
	return `# DevLink bash completion
_devlink() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    case "${prev}" in
        devlink)
            COMPREPLY=( $(compgen -W "env share get template bulk search suggest" -- "${cur}") )
            return 0
            ;;
        env)
            COMPREPLY=( $(compgen -W "share get template bulk search suggest" -- "${cur}") )
            return 0
            ;;
        share|get)
            COMPREPLY=( $(compgen -f -X "!*.env" -- "${cur}") )
            return 0
            ;;
        template)
            COMPREPLY=( $(compgen -W "list show create search" -- "${cur}") )
            return 0
            ;;
        bulk)
            COMPREPLY=( $(compgen -W "share" -- "${cur}") )
            return 0
            ;;
        search)
            COMPREPLY=( $(compgen -W "suggest" -- "${cur}") )
            return 0
            ;;
    esac
}

complete -F _devlink devlink`
}

func (cm *CompletionManager) GenerateZshCompletion() string {
	return `# DevLink zsh completion
_devlink() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    _arguments -C \
        '1: :->cmds' \
        '*:: :->args'

    case $state in
        cmds)
            _values 'devlink commands' \
                'env[Environment operations]' \
                'share[Share environment file]' \
                'get[Get shared environment file]' \
                'template[Template operations]' \
                'bulk[Bulk operations]' \
                'search[Search operations]'
            ;;
        args)
            case $line[1] in
                env)
                    _values 'env commands' \
                        'share[Share environment file]' \
                        'get[Get shared environment file]' \
                        'template[Template operations]' \
                        'bulk[Bulk operations]' \
                        'search[Search operations]'
                    ;;
                share|get)
                    _files -g "*.env"
                    ;;
                template)
                    _values 'template commands' \
                        'list[List templates]' \
                        'show[Show template]' \
                        'create[Create from template]' \
                        'search[Search templates]'
                    ;;
            esac
            ;;
    esac
}

compdef _devlink devlink`
}
