package util

import (
	"regexp"
	"strings"
)

// GenerateSlug creates a URL-friendly slug from a title
func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove duplicate hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// IsValidSlug checks if a slug contains only valid characters
func IsValidSlug(slug string) bool {
	// Slug should only contain lowercase letters, numbers, and hyphens
	match, _ := regexp.MatchString("^[a-z0-9-]+$", slug)
	return match && !strings.HasPrefix(slug, "-") && !strings.HasSuffix(slug, "-")
}
