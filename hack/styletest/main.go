package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/charmbracelet/x/exp/charmtone"
	"golang.org/x/term"

	"github.com/DanStough/epok/internal/styles"
)

// This script is a hack to display the application colors and style.

func main() {

	view := flag.String("view", "", "style view to display")
	flag.Parse()

	isDark := styles.TerminalIsDark()
	fmt.Printf("Your terminal is dark: %t\n", isDark)

	fmt.Println("Terminal has interactive stdout:\t", term.IsTerminal(int(os.Stdout.Fd())))
	fmt.Println("Terminal has interactive stdin:\t", term.IsTerminal(int(os.Stdin.Fd())))
	fmt.Println("Terminal has interactive stderr:\t", term.IsTerminal(int(os.Stderr.Fd())))

	switch *view {
	case "charm":
		// Print the standard colors that come with charm
		printCharmColors()
	case "fang-palette":
		// Print the palette colors that are the default for fang.
		theme := styles.NewCharmTheme()
		printPalette(theme, isDark)
	case "epok-palette":
		// Print the palette colors that are selected for Epok
		theme := styles.NewEpokTheme()
		printPalette(theme, isDark)
	default:
		fmt.Println("Valid options are charm, fang-palette, epok-palette")
	}
}

// Charm Color Palette
func printCharmColors() {
	labelStyle := lipgloss.NewStyle().Width(3).Align(lipgloss.Right)
	nameStyle := lipgloss.NewStyle().Width(12).Align(lipgloss.Right)
	swatchStyle := lipgloss.NewStyle().Width(6)

	keys := charmtone.Keys()

	data := [][]string{}
	for idx, key := range keys {
		row := []string{
			fmt.Sprintf("%d", idx),
			key.String(),
			"", // Background filled in later
			palette(key),
		}
		data = append(data, row)
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			key := keys[row]
			color := lipgloss.Color(key.Hex())
			switch col {
			case 2:
				return swatchStyle.Background(color)
			case 1, 3:
				return nameStyle.Foreground(color)
			default:
				return labelStyle.Foreground(color)
			}
		})

	lipgloss.Println(t)
}

func palette(key charmtone.Key) string {
	switch {
	case key.IsPrimary():
		return "primary"
	case key.IsSecondary():
		return "secondary"
	case key.IsTertiary():
		return "tertiary"
	default:
		return ""
	}
}

type combo struct {
	name  string
	color color.Color
}

func printPalette(theme *styles.Theme, isDark bool) {
	c := lipgloss.LightDark(isDark)

	mapping := []combo{
		{"Primary", theme.Primary},
		{"Accent", c(theme.LightTheme.Accent, theme.DarkTheme.Accent)},
		{"AccentAlt", c(theme.LightTheme.AccentAlt, theme.DarkTheme.AccentAlt)},
		{"AccentAltTwo", c(theme.LightTheme.AccentAltTwo, theme.DarkTheme.AccentAltTwo)},
		{"Text", c(theme.LightTheme.Text, theme.DarkTheme.Text)},
		{"TextSubdued", c(theme.LightTheme.TextSubdued, theme.DarkTheme.TextSubdued)},
		{"TextMostSubdued", c(theme.LightTheme.TextMostSubdued, theme.DarkTheme.TextMostSubdued)},
		{"CodeBackground", c(theme.LightTheme.CodeBackground, theme.DarkTheme.CodeBackground)},
		{"ErrorText", theme.Error},
		{"ErrorBackground", theme.ErrorBackground},
	}

	nameStyle := lipgloss.NewStyle().Width(18).Align(lipgloss.Right)
	swatchStyle := lipgloss.NewStyle().Width(6)

	data := [][]string{}
	for _, pair := range mapping {
		row := []string{
			pair.name,
			"", // Background filled in later
		}
		data = append(data, row)
	}

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Rows(data...).
		StyleFunc(func(row, col int) lipgloss.Style {
			pair := mapping[row]
			color := pair.color
			switch col {
			case 1:
				return swatchStyle.Background(color)
			default:
				return nameStyle.Foreground(color)
			}
		})

	if _, err := lipgloss.Println(t); err != nil {
		panic(err)
	}
}
