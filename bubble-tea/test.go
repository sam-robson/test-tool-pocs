package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Test struct {
    Description string `json:"description"`
    Command     string `json:"command"` 
    Type        string `json:"type"`
}

type Config struct {
    Tests map[string]Test `json:"tests"`
}

type item struct {
    path        string
    test        Test
    selected    bool
    isAvailable bool
}

type model struct {
    items      []item
    cursor     int
    rootDir    string
    currentDir string
}

const tickRate = time.Second / 10

type tickMsg struct{}

func tick() tea.Cmd {
    return tea.Tick(tickRate, func(time.Time) tea.Msg {
        return tickMsg{}
    })
}

func (m model) Init() tea.Cmd {
    return tick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        return m, tick()
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "up":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down":
            if m.cursor < len(m.items)-1 {
                m.cursor++
            }
        case "space", " ":
            if m.items[m.cursor].isAvailable {
                m.items[m.cursor].selected = !m.items[m.cursor].selected
            }
        case "enter":
            // Exit TUI temporarily
            fmt.Print("\033[2J") // Clear screen
            
        // Run tests synchronously
        for _, item := range m.items {
            if !item.selected {
                continue
            }

            fmt.Printf("\n=== Running %s ===\n", item.path)
            
            parts := strings.Fields(item.test.Command)
            testPath := filepath.Join(m.rootDir, item.path)
            
            cmd := exec.Command(parts[0], append(parts[1:], testPath)...)
            cmd.Dir = m.rootDir
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            if err := cmd.Run(); err != nil {
                fmt.Printf("\nError running test %s: %v\n", item.path, err)
            }
        }
        // Reset selections
        for i := range m.items {
            m.items[i].selected = false
        }

        // Wait for any key press
        fmt.Printf("\nPress any key to return to test selection...")
        reader := bufio.NewReader(os.Stdin)
        reader.ReadByte()

        // Clear screen before returning to TUI
        fmt.Print("\033[H\033[2J") // Move cursor to home position and clear screen

        // Return to TUI
        return m, tea.EnterAltScreen
        }
    }
    return m, nil
}

func (m *model) View() string {
    titleStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("12")).
        Bold(true).
        MarginBottom(1).
        Padding(0, 2)

    dirStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        MarginBottom(2).
        Padding(0, 2)
    
    containerStyle := lipgloss.NewStyle().
        Padding(1).
        MarginTop(1)
    
    itemStyle := lipgloss.NewStyle().
        PaddingLeft(2)
    
    disabledStyle := itemStyle.Copy().
        Foreground(lipgloss.Color("241"))
    
    content := strings.Builder{}
    
    content.WriteString(titleStyle.Render("Test Runner") + "\n")
    content.WriteString(dirStyle.Render(fmt.Sprintf("Directory: %s", m.currentDir)) + "\n")
    
    // Test items
    for i, item := range m.items {
        cursor := "  " // Two spaces for alignment
        if m.cursor == i {
            cursor = "❯ "
        }

        checkbox := "[ ]"
        if item.selected {
            checkbox = "[✓]"
        }

        // Format the test details with fixed widths
        testInfo := fmt.Sprintf("%-40s  %-15s  %s",
            item.path,
            "["+item.test.Type+"]",
            item.test.Description,
        )

        style := itemStyle
        if !item.isAvailable {
            style = disabledStyle
        }

        line := style.Render(fmt.Sprintf("%s%s %s", cursor, checkbox, testInfo))
        content.WriteString(line + "\n")
    }

    footerStyle := lipgloss.NewStyle().
        MarginTop(2).
        Foreground(lipgloss.Color("241")).
        Padding(0, 2)
    
    content.WriteString(footerStyle.Render("space: select • enter: run • q: quit"))

    // Wrap everything in container
    return containerStyle.Render(content.String())
}

func findRepoRoot() (string, error) {
    dir, err := os.Getwd()
    if err != nil {
        return "", err
    }

    for {
        if _, err := os.Stat(filepath.Join(dir, "config.json")); err == nil {
            return dir, nil
        }

        parent := filepath.Dir(dir)
        if parent == dir {
            return "", fmt.Errorf("config.json not found")
        }
        dir = parent
    }
}

func loadConfig(rootDir string) (Config, error) {
    data, err := os.ReadFile(filepath.Join(rootDir, "config.json"))
    if err != nil {
        return Config{}, err
    }

    var config Config
    err = json.Unmarshal(data, &config)
    return config, err
}

func main() {
    rootDir, err := findRepoRoot()
    if err != nil {
        fmt.Printf("Error finding repository root: %v\n", err)
        os.Exit(1)
    }

    config, err := loadConfig(rootDir)
    if err != nil {
        fmt.Printf("Error loading config: %v\n", err)
        os.Exit(1)
    }

    currentDir, err := os.Getwd()
    if err != nil {
        fmt.Printf("Error getting current directory: %v\n", err)
        os.Exit(1)
    }

    // Convert to relative path from root
    relDir, err := filepath.Rel(rootDir, currentDir)
    if err != nil {
        fmt.Printf("Error getting relative path: %v\n", err)
        os.Exit(1)
    }

    // Handle root directory case
    if relDir == "." {
        relDir = "" // Empty string will match all paths
    }

    var items []item
    for path, test := range config.Tests {
        isAvailable := strings.HasPrefix(path, relDir)
        items = append(items, item{
            path:        path,
            test:        test,
            isAvailable: isAvailable,
        })
    }

    initialModel := model{
        items:      items,
		cursor:    0,
        rootDir:    rootDir,
        currentDir: relDir,
    }

    p := tea.NewProgram(&initialModel)
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running program: %v\n", err)
        os.Exit(1)
    }
}