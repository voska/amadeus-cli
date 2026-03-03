package cmd

import "github.com/voska/amadeus-cli/internal/output"

type SchemaCmd struct {
	Command string `arg:"" optional:"" help:"Show schema for a specific command"`
}

func (c *SchemaCmd) Run(g *Globals) error {
	if c.Command != "" {
		return output.Write(g.Ctx, commandSchema(c.Command, g.Version))
	}
	return output.Write(g.Ctx, fullSchema(g.Version))
}

func fullSchema(version string) map[string]any {
	return map[string]any{
		"name":    "amadeus",
		"version": version,
		"commands": []map[string]any{
			{
				"name": "auth",
				"help": "Manage authentication",
				"subcommands": []map[string]any{
					{"name": "login", "help": "Authenticate with Amadeus API"},
					{"name": "status", "help": "Show authentication status"},
					{"name": "logout", "help": "Remove stored credentials"},
				},
			},
			{
				"name": "flights",
				"help": "Search and price flights",
				"subcommands": []map[string]any{
					{
						"name": "search",
						"help": "Search for flight offers",
						"flags": []map[string]any{
							{"name": "--from", "required": true, "help": "Origin IATA code (e.g., JFK)"},
							{"name": "--to", "required": true, "help": "Destination IATA code (e.g., CDG)"},
							{"name": "--date", "required": true, "help": "Departure date (YYYY-MM-DD)"},
							{"name": "--return", "help": "Return date for round-trip (YYYY-MM-DD)"},
							{"name": "--adults", "default": "1", "help": "Number of adults (1-9)"},
							{"name": "--children", "default": "0", "help": "Number of children (0-9)"},
							{"name": "--class", "help": "Travel class (ECONOMY, PREMIUM_ECONOMY, BUSINESS, FIRST)"},
							{"name": "--nonstop", "help": "Direct flights only"},
							{"name": "--currency", "help": "Currency code (e.g., USD, EUR)"},
							{"name": "--max-price", "help": "Maximum price (no decimals)"},
							{"name": "--max", "default": "10", "help": "Maximum number of results (1-250)"},
						},
					},
					{"name": "price", "help": "Confirm pricing for a flight offer", "flags": []map[string]any{{"name": "--offer-id", "required": true, "help": "Flight offer ID"}}},
					{"name": "seatmap", "help": "View seat map for a flight offer", "flags": []map[string]any{{"name": "--offer-id", "required": true, "help": "Flight offer ID"}}},
				},
			},
			{
				"name": "hotels",
				"help": "Search and book hotels",
				"subcommands": []map[string]any{
					{
						"name": "search",
						"help": "Search hotels by city or location",
						"flags": []map[string]any{
							{"name": "--city", "help": "City IATA code (e.g., PAR)"},
							{"name": "--lat", "help": "Latitude"},
							{"name": "--lng", "help": "Longitude"},
							{"name": "--radius", "default": "5", "help": "Search radius in km"},
							{"name": "--ratings", "help": "Hotel ratings to filter (1-5)"},
						},
					},
					{
						"name": "offers",
						"help": "Get offers for a specific hotel",
						"flags": []map[string]any{
							{"name": "--hotel-id", "required": true, "help": "Amadeus hotel ID"},
							{"name": "--checkin", "required": true, "help": "Check-in date (YYYY-MM-DD)"},
							{"name": "--checkout", "required": true, "help": "Check-out date (YYYY-MM-DD)"},
							{"name": "--adults", "default": "1", "help": "Number of adults"},
							{"name": "--rooms", "default": "1", "help": "Number of rooms"},
							{"name": "--currency", "help": "Currency code"},
						},
					},
					{
						"name": "book",
						"help": "Book a hotel offer",
						"flags": []map[string]any{
							{"name": "--offer-id", "required": true, "help": "Hotel offer ID from search results"},
							{"name": "--guest-name", "required": true, "help": "Guest full name"},
							{"name": "--guest-email", "required": true, "help": "Guest email address"},
						},
					},
				},
			},
			{
				"name": "airports",
				"help": "Search airports",
				"subcommands": []map[string]any{
					{"name": "search", "help": "Search airports by keyword (autocomplete)", "args": []string{"keyword"}},
				},
			},
			{
				"name": "airlines",
				"help": "Look up airlines",
				"subcommands": []map[string]any{
					{"name": "lookup", "help": "Look up airline by IATA code", "args": []string{"code"}},
				},
			},
			{
				"name": "schema",
				"help": "Show CLI schema for agent introspection",
				"args": []string{"command?"},
			},
			{
				"name": "exit-codes",
				"help": "Show exit code reference",
			},
		},
		"global_flags": []map[string]any{
			{"name": "--json", "short": "-j", "help": "Output JSON to stdout", "aliases": []string{"--machine"}},
			{"name": "--plain", "short": "-p", "help": "Output TSV, no color"},
			{"name": "--select", "short": "-s", "help": "Project output to fields (dot-path)", "aliases": []string{"--fields"}},
			{"name": "--results-only", "help": "Strip response metadata, return data array only"},
			{"name": "--test", "help": "Use test environment (test.api.amadeus.com)"},
			{"name": "--dry-run", "help": "Show what would happen without making API calls"},
			{"name": "--no-input", "help": "Never prompt for input"},
			{"name": "--verbose", "short": "-v", "help": "Verbose output to stderr"},
		},
		"environment_variables": []map[string]any{
			{"name": "AMADEUS_API_KEY", "help": "API key (overrides config file)"},
			{"name": "AMADEUS_API_SECRET", "help": "API secret (overrides config file)"},
			{"name": "AMADEUS_AUTO_JSON", "help": "Set to 1 to auto-detect non-TTY and output JSON"},
			{"name": "AMADEUS_CONFIG_DIR", "help": "Override config directory (default: ~/.config/amadeus)"},
		},
	}
}

func commandSchema(name, version string) map[string]any {
	schema := fullSchema(version)
	commands, ok := schema["commands"].([]map[string]any)
	if !ok {
		return map[string]any{"error": "command not found: " + name}
	}
	for _, cmd := range commands {
		if cmd["name"] == name {
			return cmd
		}
	}
	return map[string]any{"error": "command not found: " + name}
}
