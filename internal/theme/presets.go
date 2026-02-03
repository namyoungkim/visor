package theme

// presets contains all predefined themes.
var presets = map[string]*Theme{
	"default":   defaultTheme,
	"powerline": powerlineTheme,
	"gruvbox":   gruvboxTheme,
	"nord":      nordTheme,
}

// defaultTheme uses standard terminal colors with ASCII separators.
var defaultTheme = &Theme{
	Name: "default",
	Colors: ColorPalette{
		Normal:   "white",
		Warning:  "yellow",
		Critical: "red",
		Good:     "green",
		Primary:  "cyan",
		Secondary: "blue",
		Muted:    "gray",
		Backgrounds: []string{},
	},
	Separators: DefaultSeparators(),
	Powerline:  false,
}

// powerlineTheme uses Powerline font glyphs with colored backgrounds.
var powerlineTheme = &Theme{
	Name: "powerline",
	Colors: ColorPalette{
		Normal:   "white",
		Warning:  "yellow",
		Critical: "red",
		Good:     "green",
		Primary:  "cyan",
		Secondary: "blue",
		Muted:    "gray",
		// Background colors cycle for adjacent widgets
		Backgrounds: []string{"#3c3836", "#504945", "#665c54", "#7c6f64"},
	},
	Separators: PowerlineSeparators(),
	Powerline:  true,
}

// gruvboxTheme uses Gruvbox color scheme (warm retro).
// https://github.com/morhetz/gruvbox
var gruvboxTheme = &Theme{
	Name: "gruvbox",
	Colors: ColorPalette{
		Normal:   "#ebdbb2", // fg
		Warning:  "#fabd2f", // bright yellow
		Critical: "#fb4934", // bright red
		Good:     "#b8bb26", // bright green
		Primary:  "#83a598", // bright blue
		Secondary: "#8ec07c", // bright aqua
		Muted:    "#928374", // gray
		Backgrounds: []string{"#3c3836", "#504945", "#665c54"},
	},
	Separators: DefaultSeparators(),
	Powerline:  false,
}

// gruvboxPowerlineTheme combines Gruvbox colors with Powerline styling.
var gruvboxPowerlineTheme = &Theme{
	Name: "gruvbox-powerline",
	Colors: ColorPalette{
		Normal:   "#ebdbb2",
		Warning:  "#fabd2f",
		Critical: "#fb4934",
		Good:     "#b8bb26",
		Primary:  "#83a598",
		Secondary: "#8ec07c",
		Muted:    "#928374",
		Backgrounds: []string{"#3c3836", "#504945", "#665c54"},
	},
	Separators: PowerlineSeparators(),
	Powerline:  true,
}

// nordTheme uses Nord color scheme (arctic, bluish).
// https://www.nordtheme.com/
var nordTheme = &Theme{
	Name: "nord",
	Colors: ColorPalette{
		Normal:   "#eceff4", // snow storm
		Warning:  "#ebcb8b", // aurora yellow
		Critical: "#bf616a", // aurora red
		Good:     "#a3be8c", // aurora green
		Primary:  "#81a1c1", // frost
		Secondary: "#88c0d0", // frost
		Muted:    "#4c566a", // polar night
		Backgrounds: []string{"#3b4252", "#434c5e", "#4c566a"},
	},
	Separators: DefaultSeparators(),
	Powerline:  false,
}

// nordPowerlineTheme combines Nord colors with Powerline styling.
var nordPowerlineTheme = &Theme{
	Name: "nord-powerline",
	Colors: ColorPalette{
		Normal:   "#eceff4",
		Warning:  "#ebcb8b",
		Critical: "#bf616a",
		Good:     "#a3be8c",
		Primary:  "#81a1c1",
		Secondary: "#88c0d0",
		Muted:    "#4c566a",
		Backgrounds: []string{"#3b4252", "#434c5e", "#4c566a"},
	},
	Separators: PowerlineSeparators(),
	Powerline:  true,
}

func init() {
	// Register additional themes
	presets["gruvbox-powerline"] = gruvboxPowerlineTheme
	presets["nord-powerline"] = nordPowerlineTheme
}
