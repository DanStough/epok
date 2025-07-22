// Package styles provides the pre-computed lipgloss styles and colors used throughout the application.
package styles

import (
	"image/color"
	"os"

	"github.com/DanStough/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/charmbracelet/x/term"
)

// How this file is organized:
//      Theme  		The palette of colors used for both light and dark modes.
//       |
//       V
//    Schemes    The implementation of the colors for either a light or dark theme and where they show up in the stylesheet
//       |
//       V
//     Sheet		The "style sheet" for the application, which is the application of the ColorScheme plus other styling primitives.
//
// The colors are borrowed from the charm color tones, mostly as a convenience.
//
// Theme colors are generally organized with more saturated hues up front, then in order of brightness.
//
//*** Theme guide ***
// * primary: used for help section headings and table borders
// * errorText: text color for the error banner
// * errorBackground: color for the error banner
//
// * accent: 		used to highlight text and also flags in the help
// * accentAlt: 	used for command names and args in the help
// * accentAltTwo: 	used for quoted strings in the help
//
// This gray palette is used for emphasis on different text
// * text: 				the default text color
// * subduedText:		tbd
// * mostSubduedText:	tbd
//
// * codeBackground: 	used the background for examples and usage in the help.
//
// User ../../hack/styletest to visualize these colors

func TerminalIsDark() bool {
	var isDark bool
	if term.IsTerminal(os.Stdout.Fd()) {
		isDark = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	}
	return isDark
}

type Theme struct {
	Primary         color.Color
	Error           color.Color
	ErrorBackground color.Color

	DarkTheme  OptionalTheme
	LightTheme OptionalTheme
}

type OptionalTheme struct {
	Accent       color.Color
	AccentAlt    color.Color
	AccentAltTwo color.Color

	Text            color.Color
	TextSubdued     color.Color
	TextMostSubdued color.Color

	CodeBackground color.Color
}

// NewCharmTheme returns a color theme based on the colors used by charmbracelet/fang
func NewCharmTheme() *Theme {
	return &Theme{
		Primary:         charmtone.Charple,
		Error:           charmtone.Butter,
		ErrorBackground: charmtone.Cherry,

		DarkTheme: OptionalTheme{
			Accent:          charmtone.Guac,
			AccentAlt:       charmtone.Cheeky,
			AccentAltTwo:    charmtone.Salmon,
			Text:            charmtone.Ash,
			TextSubdued:     charmtone.Squid,
			TextMostSubdued: charmtone.Oyster,
			CodeBackground:  lipgloss.Color("#2F2E36"),
		},
		LightTheme: OptionalTheme{
			Accent:       charmtone.Damson,
			AccentAlt:    charmtone.Pony,
			AccentAltTwo: charmtone.Coral,
			// Because the theme is based on a white background, the ordering is reversed.
			CodeBackground:  charmtone.Salt,
			TextMostSubdued: charmtone.Smoke,
			TextSubdued:     charmtone.Squid,
			Text:            charmtone.Charcoal,
		},
	}
}

// NewEpokTheme returns a color theme based on the custom colors selected by epok.
// [Color Wheel](https://www.canva.com/colors/color-wheel/).
func NewEpokTheme() *Theme {
	return &Theme{
		Primary:         lipgloss.Color("#FF8B00"),
		Error:           charmtone.Ash,
		ErrorBackground: lipgloss.Color("#FF4700"),

		DarkTheme: OptionalTheme{
			Accent:          lipgloss.Color("#00F4FF"), // This is a complementary color based on the FF0B00.
			AccentAlt:       lipgloss.Color("#FF0074"), // These are analogous colors with FF0B00 as the center
			AccentAltTwo:    lipgloss.Color("#FF00F1"),
			Text:            charmtone.Ash,
			TextSubdued:     charmtone.Squid,
			TextMostSubdued: charmtone.Oyster,
			CodeBackground:  lipgloss.Color("#2F2E36"),
		},
		LightTheme: OptionalTheme{
			Accent:       lipgloss.Color("#00B7BF"), // These are slightly darker versions of the dark accent colors
			AccentAlt:    lipgloss.Color("#CD005E"),
			AccentAltTwo: lipgloss.Color("#BC00B1"),
			// Because the theme is based on a white background, the ordering is reversed.
			CodeBackground:  lipgloss.Color("#FFF7F3"),
			TextMostSubdued: charmtone.Smoke,
			TextSubdued:     charmtone.Squid,
			Text:            charmtone.Charcoal,
		},
	}
}

// FangColorScheme is the colorscheme used by `charmbracelet/fang`.
// Its signature is a fang.ColorSchemeFunc that can be used with fang.WithColorSchemeFunc on initialization.
// It is derived from fang.DefaultColorScheme.
func (t *Theme) FangColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	return fang.ColorScheme{
		Base:  c(t.LightTheme.Text, t.DarkTheme.Text), // Main Text
		Title: t.Primary,                              // Heading in Help

		// Codeblock Styles
		Codeblock: c(t.LightTheme.CodeBackground, t.DarkTheme.CodeBackground),   // This is the background for codeblocks and spans
		Program:   t.Primary,                                                    // Name of the program in the help codeblock. For simplicity this is collapsed with the primary color.
		Comment:   c(t.LightTheme.TextMostSubdued, t.DarkTheme.TextMostSubdued), // Styles comments in the example codeblock.

		// Styles Codeblocks and Usage
		Command:        c(t.LightTheme.AccentAlt, t.DarkTheme.AccentAlt),       // Styles command name in usage and examples.
		Flag:           c(t.LightTheme.Accent, t.DarkTheme.Accent),             // Color for flags in usage, examples and errors.
		Argument:       c(t.LightTheme.Text, t.DarkTheme.Text),                 // Styles usage and example args
		DimmedArgument: c(t.LightTheme.TextSubdued, t.LightTheme.TextSubdued),  // Alternate style of usage and example args
		QuotedString:   c(t.LightTheme.AccentAltTwo, t.DarkTheme.AccentAltTwo), // Quoted strings in usage and examples.

		Description: c(t.LightTheme.Text, t.DarkTheme.Text),               // text for flag and command descriptions
		FlagDefault: c(t.LightTheme.TextSubdued, t.DarkTheme.TextSubdued), // text for flag default values in descriptions
		ErrorHeader: [2]color.Color{
			t.Error,           // Text for the error message
			t.ErrorBackground, // Background for the error message
		},
	}
}

type scheme struct {
	Text        color.Color
	TextSubdued color.Color
	TextKeyword color.Color
	TextBanner  color.Color

	TableBorder            color.Color
	TableCellText          color.Color
	TableCellTextAlternate color.Color
}

// epokColorScheme provides the colors used is specific components of epok.
// It is not meant to be called externally, only through Theme.Sheet.
func (t *Theme) epokColorScheme() scheme {
	isDark := TerminalIsDark()
	c := lipgloss.LightDark(isDark)

	return scheme{
		Text:        c(t.LightTheme.Text, t.DarkTheme.Text),
		TextSubdued: c(t.LightTheme.TextSubdued, t.DarkTheme.TextSubdued),
		TextKeyword: c(t.LightTheme.Accent, t.DarkTheme.Accent),
		TextBanner:  t.Primary,

		TableBorder:            t.Primary,
		TableCellText:          c(t.LightTheme.Text, t.DarkTheme.Text),
		TableCellTextAlternate: c(t.LightTheme.TextSubdued, t.DarkTheme.TextSubdued),
	}
}

// Sheet represents all the lipgloss styles used within epok.
// This excludes the colors that are exported to `charmbracelet/fang`, which makes it's own styles.
type Sheet struct {
	Text        lipgloss.Style
	TextSubdued lipgloss.Style
	Keyword     lipgloss.Style
	Banner      lipgloss.Style

	Table TableStyle
}

// TableStyle are the table-specific lipgloss styles within the overall style sheet.
type TableStyle struct {
	Header          lipgloss.Style
	Cell            lipgloss.Style
	OddRow          lipgloss.Style
	EvenRow         lipgloss.Style
	Border          lipgloss.Style
	BorderThickness lipgloss.Border
}

// Sheet returns the computed styles used in epok based on common usages of lipgloss.
func (t *Theme) Sheet() *Sheet {
	es := t.epokColorScheme()

	baseCellStyle := lipgloss.NewStyle().Padding(0, 1) // Make this width configurable based on the table.

	return &Sheet{
		Text: lipgloss.NewStyle().
			Foreground(es.Text),
		TextSubdued: lipgloss.NewStyle().
			Foreground(es.TextSubdued),
		Keyword: lipgloss.NewStyle().
			Foreground(es.TextKeyword).
			Bold(true),
		Banner: lipgloss.NewStyle().
			Foreground(es.TextBanner),

		Table: TableStyle{
			Header:          lipgloss.NewStyle().Foreground(es.TableBorder).Bold(true).Align(lipgloss.Center),
			Cell:            baseCellStyle,
			OddRow:          baseCellStyle.Foreground(es.TableCellTextAlternate),
			EvenRow:         baseCellStyle.Foreground(es.TableCellText),
			Border:          lipgloss.NewStyle().Foreground(es.TableBorder),
			BorderThickness: lipgloss.ThickBorder(),
		},
	}
}
