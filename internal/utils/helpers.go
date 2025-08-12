package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (most common for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header (nginx)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Check CF-Connecting-IP (Cloudflare)
	if cfip := r.Header.Get("CF-Connecting-IP"); cfip != "" {
		return cfip
	}

	// Fallback to RemoteAddr
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}

// Helper function to extract browser info
func getBrowserInfo(userAgent string) (browser, version string) {
	// Simple browser detection - you might want to use a library like "github.com/mileusna/useragent"
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "chrome"):
		return "Chrome", extractVersion(ua, "chrome/")
	case strings.Contains(ua, "firefox"):
		return "Firefox", extractVersion(ua, "firefox/")
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		return "Safari", extractVersion(ua, "version/")
	case strings.Contains(ua, "edge"):
		return "Edge", extractVersion(ua, "edge/")
	default:
		return "Unknown", ""
	}
}

func extractVersion(ua, prefix string) string {
	if idx := strings.Index(ua, prefix); idx != -1 {
		start := idx + len(prefix)
		end := start
		for end < len(ua) && (ua[end] >= '0' && ua[end] <= '9' || ua[end] == '.') {
			end++
		}
		return ua[start:end]
	}
	return ""
}

func GetDeviceType(r *http.Request) string {
	uaMobile := r.Header.Get("Sec-Ch-Ua-Mobile")
	if uaMobile == "?1" {
		return "Mobile"
	}
	return "Desktop"
}

func GetBrowser(r *http.Request) string {
	// Try Sec-CH-UA first (modern browsers)
	if secChUa := r.Header.Get("Sec-Ch-Ua"); secChUa != "" {
		// Parse Sec-CH-UA: "Google Chrome";v="91", "Chromium";v="91"
		if strings.Contains(secChUa, "Chrome") {
			return "Chrome"
		}
		if strings.Contains(secChUa, "Firefox") {
			return "Firefox"
		}
		if strings.Contains(secChUa, "Safari") {
			return "Safari"
		}
	}

	// Fallback to User-Agent parsing
	browser, _ := getBrowserInfo(r.Header.Get("User-Agent"))
	return browser
}

func GetBrowserVersion(r *http.Request) string {
	// Try Sec-CH-UA-Full-Version first
	if version := r.Header.Get("Sec-Ch-Ua-Full-Version"); version != "" {
		return strings.Trim(version, `"`)
	}

	// Fallback to User-Agent parsing
	_, version := getBrowserInfo(r.Header.Get("User-Agent"))
	return version
}

func GetOS(r *http.Request) string {
	if platform := r.Header.Get("Sec-Ch-Ua-Platform"); platform != "" {
		return strings.Trim(platform, `"`)
	}

	// Fallback to User-Agent parsing
	ua := r.Header.Get("User-Agent")
	switch {
	case strings.Contains(ua, "Windows"):
		return "Windows"
	case strings.Contains(ua, "Macintosh"):
		return "macOS"
	case strings.Contains(ua, "Linux"):
		return "Linux"
	case strings.Contains(ua, "Android"):
		return "Android"
	case strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad"):
		return "iOS"
	default:
		return "Unknown"
	}
}

func GetOSVersion(r *http.Request) string {
	if version := r.Header.Get("Sec-Ch-Ua-Platform-Version"); version != "" {
		return strings.Trim(version, `"`)
	}
	return ""
}

func GetLocation(r *http.Request) string {
	// Try various location headers
	if city := r.Header.Get("X-Geo-City"); city != "" {
		return city
	}
	if city := r.Header.Get("CF-IPCity"); city != "" { // Cloudflare
		return city
	}
	if country := r.Header.Get("CF-IPCountry"); country != "" { // Cloudflare
		return country
	}
	if country := r.Header.Get("X-Country-Code"); country != "" {
		return country
	}
	return ""
}
