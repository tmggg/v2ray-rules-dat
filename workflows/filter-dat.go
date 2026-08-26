// Command filterdat keeps only the specified lists in a geosite.dat file.
//
// Usage: go run ./filter-dat.go <input.dat> <output.dat> <keep1,keep2,...>
// List names are case-insensitive. Entries not in the keep list are dropped.
package main

import (
	"fmt"
	"os"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: filterdat <input.dat> <output.dat> <keep1,keep2,...>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	var list router.GeoSiteList
	if err := proto.Unmarshal(data, &list); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}

	keep := make(map[string]bool)
	for _, name := range strings.Split(os.Args[3], ",") {
		if name = strings.TrimSpace(name); name != "" {
			keep[strings.ToUpper(name)] = true
		}
	}

	filtered := new(router.GeoSiteList)
	for _, entry := range list.Entry {
		if keep[strings.ToUpper(entry.CountryCode)] {
			filtered.Entry = append(filtered.Entry, entry)
		}
	}

	out, err := proto.Marshal(filtered)
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(os.Args[2], out, 0644); err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	fmt.Printf("Kept %d of %d lists, written to %s\n", len(filtered.Entry), len(list.Entry), os.Args[2])
}
