package utils

// IsNumber 是否是数字
func IsNumber(c rune) bool {
	return c >= '0' && c <= '9'
}
