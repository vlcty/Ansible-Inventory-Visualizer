package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-ini/ini"
	"github.com/olekukonko/tablewriter"
)

const (
	REPLACE_ALL int = -1
	EXIT_ERROR  int = 1
)

func main() {
	inventoryFilePath := flag.String("inventory", "", "Path to the inventory file")
	rawServerFilter := flag.String("filter", "", "Comma seperated list to filter servers")
	flag.Parse()

    CheckInventoryFilePath(*inventoryFilePath)
    serverFilter := ParseServerFilter(*rawServerFilter)

	table := BuildTable(*inventoryFilePath, serverFilter)
	fmt.Println(table)
}

func ParseServerFilter(rawServerFilter string) []string {
    serverFilter := make([]string, 0)

    for _, server := range strings.Split(rawServerFilter, ",") {
        /*
         * This bugged me a bit in go. Extract from the strings.Split documentation:
         *    If s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.
         * However my string is empty, so strings.Split doesn't add s but an empty string there is never a slice with length of 0
         * Filter empty server out
         */
        if server != "" {
            serverFilter = append(serverFilter, server)
        }
    }

    return serverFilter
}

func BuildTable(inventoryFilePath string, serverFilter []string) string {
	cfg := LoadPassedIniFile(inventoryFilePath)

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	headers := ExtractHeaders(cfg)
	servers := ExtractServers(cfg, serverFilter)
	data := make([][]string, 0)

	table.SetHeader(headers)

	for _, currentServer := range servers {
		row := make([]string, len(headers)+1)

		for index, currentSection := range headers {
			if index == 0 {
				row[index] = currentServer
			} else if cfg.Section(currentSection).HasKey(currentServer) {
				row[index] = "X"
			} else {
				row[index] = ""
			}
		}

		data = append(data, row)
	}

	table.AppendBulk(data)
	table.SetBorder(false)
	table.Render()

	return tableString.String()
}

func LoadPassedIniFile(filePath string) *ini.File {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		KeyValueDelimiters: " \n",
	}, filePath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EXIT_ERROR)
	}

	return cfg
}

func ExtractHeaders(cfg *ini.File) []string {
	headers := make([]string, 0)

	headers = append(headers, "Server")

	for _, item := range cfg.SectionStrings() {
		if !IsUnwantedSection(item) {
			headers = append(headers, item)
		}
	}

	sort.Strings(headers)

	return headers
}

func ExtractServers(cfg *ini.File, serverFilter []string) []string {
	servers := make([]string, 0)

	for _, currentSection := range cfg.Sections() {
		if IsUnwantedSection(currentSection.Name()) {
			continue
		}

		for _, currentServer := range currentSection.KeyStrings() {
            if len(serverFilter) != 0 && ! IsServerInSlice(currentServer, serverFilter) {
                continue
            }

			if ! IsServerInSlice(currentServer, servers) {
				servers = append(servers, currentServer)
			}
		}
	}

	sort.Strings(servers)

	return servers
}

func IsUnwantedSection(sectionName string) bool {
	unwanted := false

	if sectionName == "DEFAULT" || strings.HasSuffix(sectionName, ":vars") || strings.Index(sectionName, ":") != -1 {
		unwanted = true
	}

	return unwanted
}

func IsServerInSlice(server string, servers []string) bool {
	found := false

	for _, currentServer := range servers {
		if currentServer == server {
			found = true
			break
		}
	}

	return found
}

func CheckInventoryFilePath(inventoryFilePath string) {
    if inventoryFilePath == "" {
        fmt.Fprintln(os.Stderr, "No inventory file path given")
        os.Exit(EXIT_ERROR)
    } else {
        _, err := os.Stat(inventoryFilePath)

        if os.IsNotExist(err) {
            fmt.Fprintln(os.Stderr, "Inventory file not found")
            os.Exit(EXIT_ERROR)
        }

    }
}
