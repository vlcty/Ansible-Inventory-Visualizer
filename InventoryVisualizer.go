package main

import (
    "fmt"
    "os"
    "strings"
    "sort"
    // "io"
    "bytes"

    "github.com/go-ini/ini"
    "github.com/olekukonko/tablewriter"
)

const (
    REPLACE_ALL int = -1
    EXIT_NO_ERROR int = 0
    EXIT_ERROR int = 1
)

func main() {
    buf := &bytes.Buffer{}

    //BuildTable(ParseInventoryFilePath(), tablewriter.NewWriter(os.Stdout))
    BuildTable(ParseInventoryFilePath(), tablewriter.NewWriter(buf))

    fmt.Println(buf.String())
}

func BuildTable(filePath string, table *tablewriter.Table) {
    cfg := LoadPassedIniFile(filePath)

    headers := ExtractHeaders(cfg)

    table.SetHeader(headers)

    data := make([][]string, 0)

    for _, currentServer := range ExtractServers(cfg) {
        row := make([]string, len(headers) + 1)

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

    os.Exit(EXIT_NO_ERROR)
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
        if ! IsUnwantedSection(item) {
            headers = append(headers, item)
        }
    }

    sort.Strings(headers)

    return headers;
}

func ExtractServers(cfg *ini.File) []string {
    servers := make([]string, 0)

    for _, currentSection := range cfg.Sections() {
        if IsUnwantedSection(currentSection.Name()) {
            continue
        }

        for _, currentServer := range currentSection.KeyStrings() {
            if IsServerAlreadyKnown(currentServer, servers) == false {
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

func IsServerAlreadyKnown(server string, servers []string) bool {
    found := false

    for _, currentServer := range servers {
        if currentServer == server {
            found = true
            break
        }
    }

    return found
}

func ParseInventoryFilePath() string {
    if len(os.Args) > 1 {
        path := os.Args[1]

        _, err := os.Stat(path)

        if os.IsNotExist(err) {
            fmt.Fprintln(os.Stderr, "Inventory file not found")
            os.Exit(EXIT_ERROR)
        }

        return path
    } else {
        fmt.Fprintln(os.Stderr, "No inventory file given")
        os.Exit(EXIT_ERROR)
        return ""
    }
}


// func ExamineStringArray(array []string) {
//     for index, item := range array {
//         fmt.Printf("Index: %-2d Value: %s\n", index, item)
//     }
// }
