package main

import "fmt"

// GetPlatform returns the corresponding Platform Routing Value based on the specified region marker.
func GetPlatform(r string) string {
	platform := platformMap[r]
	if platform == "" {
		panic(fmt.Sprintf("Region %s does not exist", r))
	}
	return platform
}

// GetRegion returns the corresponding Regional Routing Value based on the given platform.
// Apparently in v5, some resources have different conventions for regions. Instead of "EUW1.api.riotgames.com", you need to specify "europe.api.riotgames.com"
func GetRegion(r string) string {
	region := regionMap[r]
	if region == "" {
		panic(fmt.Sprintf("Region %s does not exist", r))
	}
	return region
}

func checkSpecialChars(str string) bool {
	for _, r := range str {
		if ((r < 'A') || (r > 'z') || (r < '0') || (r > '9')) && (r != '-') && (r != '_') {
			return true
		}
	}
	return false
}
