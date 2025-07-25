package middleware

import (
	"encoding/json"
	"net"
	"regexp"
	"strconv"
	"strings"
	"ujikom-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func NetworkSecurityMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/api/v1/health" || c.Path() == "/api/v1/auth/login" || c.Path() == "/api/v1/auth/register" {
			return c.Next()
		}

		if adminOverride, ok := c.Locals("admin_override").(bool); ok && adminOverride {
			c.Set("X-Network-Security", "admin-bypassed")
			return c.Next()
		}

		adminKey := c.Get("X-Admin-Override")
		if adminKey == "DIMAS-ANJAY-MABAR" {
			c.Set("X-Network-Security", "admin-bypassed")
			return c.Next()
		}

		clientIP := getClientIP(c)
		
		if !isValidIPRange(clientIP) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Invalid IP range")
		}

		// if isVPNDetected(c) {
		// 	return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: VPN usage detected")
		// }

		// if !isValidNetworkType(c) {
		// 	return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Invalid network type")
		// }

		if !isValidWiFiSSID(c) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Invalid WiFi network")
		}

		if !isValidCellularNetwork(c) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Insecure cellular network")
		}

		if isFakeGPSDetected(c) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Fake GPS detected")
		}

		c.Set("X-Network-Security", "validated")
		c.Set("X-Client-IP", clientIP)
		c.Set("X-Security-Level", "high")

		return c.Next()
	}
}

// Mendapatkan IP address client yang sebenarnya
func getClientIP(c *fiber.Ctx) string {
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := c.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Cek header CF-Connecting-IP (Cloudflare)
	if cf := c.Get("CF-Connecting-IP"); cf != "" {
		return cf
	}

	return c.IP()
}

// Validasi IP range yang diizinkan
func isValidIPRange(ip string) bool {
	allowedRanges := []string{
		"10.0.0.0/8",       // Private network
		"172.16.0.0/12",    // Private network
		"192.168.0.0/16",   // Private network
		"127.0.0.0/8",     
		"103.0.0.0/8",     
		"114.0.0.0/8",     
		"202.0.0.0/8",     
		"103.156.71.94",
		"203.78.113.253",
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	for _, cidr := range allowedRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(clientIP) {
			return true
		}
	}

	return false
}

// Deteksi penggunaan VPN
func isVPNDetected(c *fiber.Ctx) bool {
	vpnHeaders := []string{
		"X-VPN-Client",
		"X-Forwarded-Proto",
		"X-Proxy-Authorization",
		"Via",
	}

	for _, header := range vpnHeaders {
		if c.Get(header) != "" {
			return true
		}
	}

	// Deteksi berdasarkan User-Agent
	userAgent := strings.ToLower(c.Get("User-Agent"))
	vpnKeywords := []string{
		"vpn",
		"proxy",
		"tunnel",
		"nordvpn",
		"expressvpn",
		"cyberghost",
		"protonvpn",
	}

	for _, keyword := range vpnKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}

	// Deteksi berdasarkan hop count (TTL)
	if ttl := c.Get("X-TTL"); ttl != "" {
		if ttlValue, err := strconv.Atoi(ttl); err == nil {
			// TTL yang terlalu rendah bisa menunjukkan VPN
			if ttlValue < 50 {
				return true
			}
		}
	}

	return false
}

// Validasi tipe jaringan
func isValidNetworkType(c *fiber.Ctx) bool {
	networkType := c.Get("X-Network-Type")

	allowedTypes := []string{
		"wifi",
		"cellular",
		"4g",
		"5g",
		"lte",
		"ethernet",
	}

	networkType = strings.ToLower(networkType)
	for _, allowed := range allowedTypes {
		if networkType == allowed {
			return true
		}
	}

	return false
}

// Validasi SSID WiFi sekolah
func isValidWiFiSSID(c *fiber.Ctx) bool {
	ssid := c.Get("X-WiFi-SSID")
	
	// if ssid == "" {
	// 	return true
	// }

	allowedSSIDs := []string{
		"JTI-3.01",
		"JTI-3.02",
		"JTI-3.03",
		"JTI-3.04",
		"JTI-3.05",
		"anjay",
	}

	for _, allowed := range allowedSSIDs {
		if strings.EqualFold(ssid, allowed) {
			return true
		}
	}

	schoolPatterns := []string{
		`^JTI-.*`,
		`^UJIKOM-.*`,
	}

	for _, pattern := range schoolPatterns {
		if matched, _ := regexp.MatchString(pattern, strings.ToUpper(ssid)); matched {
			return true
		}
	}

	return false
}

// Validasi keamanan jaringan seluler
func isValidCellularNetwork(c *fiber.Ctx) bool {
	// Ambil informasi carrier dari header
	carrier := c.Get("X-Carrier")
	networkType := c.Get("X-Network-Type")
	
	// if carrier == "" {
	// 	return true
	// }

	allowedCarriers := []string{
		"telkomsel",
		"indosat",
		"xl",
		"axis",
		"tri",
		"smartfren",
		"by.u",
	}

	carrier = strings.ToLower(carrier)
	for _, allowed := range allowedCarriers {
		if strings.Contains(carrier, allowed) {
			if isSecureCellularNetwork(networkType) {
				return true
			}
		}
	}

	return false
}

// Cek keamanan jaringan seluler
func isSecureCellularNetwork(networkType string) bool {
	secureNetworks := []string{
		"4g",
		"5g",
		"lte",
		"lte-a",
	}

	networkType = strings.ToLower(networkType)
	for _, secure := range secureNetworks {
		if networkType == secure {
			return true
		}
	}

	return false
}

func isFakeGPSDetected(c *fiber.Ctx) bool {
	body := c.Body()
	if len(body) == 0 {
		return false
	}

	var locationData struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.Unmarshal(body, &locationData); err != nil {
		return false
	}

	if isSuspiciousPrecision(locationData.Latitude, locationData.Longitude) {
		return true
	}

	if hasUnrealisticPrecision(locationData.Latitude, locationData.Longitude) {
		return true
	}

	mockHeaders := []string{
		"X-Mock-Location",
		"X-Fake-GPS",
		"X-Location-Spoofed",
	}

	for _, header := range mockHeaders {
		if c.Get(header) != "" {
			return true
		}
	}

	if accuracy := c.Get("X-GPS-Accuracy"); accuracy != "" {
		if acc, err := strconv.ParseFloat(accuracy, 64); err == nil {
			if acc == 0 || acc < 1 {
				return true
			}
		}
	}

	return false
}

func isSuspiciousPrecision(lat, lng float64) bool {
	latPrecision := getDecimalPlaces(lat)
	lngPrecision := getDecimalPlaces(lng)

	if latPrecision > 10 || lngPrecision > 10 {
		return true
	}

	if latPrecision > 12 || lngPrecision > 12 {
		return true
	}

	return false
}

func hasUnrealisticPrecision(lat, lng float64) bool {
	latStr := strconv.FormatFloat(lat, 'f', -1, 64)
	lngStr := strconv.FormatFloat(lng, 'f', -1, 64)

	if strings.Contains(latStr, "000") || strings.Contains(lngStr, "000") {
		if getDecimalPlaces(lat) > 10 || getDecimalPlaces(lng) > 10 {
			return true
		}
	}

	return false
}

func getDecimalPlaces(num float64) int {
	str := strconv.FormatFloat(num, 'f', -1, 64)
	if dotIndex := strings.Index(str, "."); dotIndex != -1 {
		return len(str) - dotIndex - 1
	}
	return 0
}

// Middleware untuk logging informasi jaringan
func NetworkInfoMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		networkInfo := map[string]string{
			"client_ip":     getClientIP(c),
			"user_agent":    c.Get("User-Agent"),
			"network_type":  c.Get("X-Network-Type"),
			"wifi_ssid":     c.Get("X-WiFi-SSID"),
			"carrier":       c.Get("X-Carrier"),
			"device_id":     c.Get("X-Device-ID"),
			"app_version":   c.Get("X-App-Version"),
		}

		// Set network info ke context untuk digunakan di controller
		c.Locals("network_info", networkInfo)

		return c.Next()
	}
}

func AdminNetworkOverrideMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		adminKey := c.Get("X-Admin-Override")
		
		if adminKey == "DIMAS-ANJAY-MABAR" {
			c.Set("X-Network-Security", "admin-override")
			c.Locals("admin_override", true)
		}

		return c.Next()
	}
}
