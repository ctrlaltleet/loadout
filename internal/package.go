package internal

import (
	"fmt"
	"strings"
)

func ListPackages(cfg *Config, platform, selectStr string) {
	fmt.Println("Available packages for", platform)

	filters := ParseSelect(selectStr)
	selectedPackages := FilterPackages(cfg.Packages, filters)

	if len(filters) == 0 {
		for name, pkg := range cfg.Packages {
			printPackage(name, pkg, platform)
		}
		return
	}

	for name, pkg := range selectedPackages {
		printPackage(name, pkg, platform)
	}
}

func printPackage(name string, pkg Package, platform string) {
	fmt.Printf(" - %s (%s)", name, formatVersion(pkg.Version))
	if len(pkg.Tags) > 0 {
		fmt.Printf(" [tags: %s]", strings.Join(pkg.Tags, ", "))
	}

	if _, ok := pkg.PlatformAssets[platform]; ok {
		fmt.Printf(" [has platform asset]")
	}

	if len(pkg.GlobalAssets) > 0 {
		fmt.Printf(" [has global assets]")
	}

	fmt.Println()
}

func ParseSelect(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func FilterPackages(allPackages map[string]Package, filters []string) map[string]Package {
	selected := make(map[string]Package)

	for _, f := range filters {
		if strings.EqualFold(f, "all") {
			return allPackages
		}
	}

	for _, f := range filters {
		for name, pkg := range allPackages {
			if strings.EqualFold(name, f) {
				selected[name] = pkg
				continue
			}
			for _, tag := range pkg.Tags {
				if strings.EqualFold(tag, f) {
					selected[name] = pkg
					break
				}
			}
		}
	}

	return selected
}