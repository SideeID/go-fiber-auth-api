package middleware

import (
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

		if isVPNDetected(c) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: VPN usage detected")
		}

		if !isValidNetworkType(c) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Invalid network type")
		}

		if !isValidWiFiSSID(c) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Invalid WiFi network")
		}

		if !isValidCellularNetwork(c) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied: Insecure cellular network")
		}

		c.Set("X-Network-Security", "validated")
		c.Set("X-Client-IP", clientIP)
		c.Set("X-Security-Level", "high")

		return c.Next()
	}
}

// Mendapatkan IP address client yang sebenarnya
func getClientIP(c *fiber.Ctx) string {
	// Cek header X-Forwarded-For
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Cek header X-Real-IP
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
	
	// if networkType == "" {
	// 	return true
	// }

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
		"SMK-NET",
		"SEKOLAH-WIFI",
		"UJIKOM-NET",
		"STUDENT-NET",
		"SCHOOL-NETWORK",
		"PENDIDIKAN-NET",
		"WIFI lemod parahh",
	}

	for _, allowed := range allowedSSIDs {
		if strings.EqualFold(ssid, allowed) {
			return true
		}
	}

	schoolPatterns := []string{
		`^SMK-.*`,
		`^SEKOLAH-.*`,
		`^PENDIDIKAN-.*`,
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
