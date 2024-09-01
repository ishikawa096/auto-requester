package utils

// Colorize the stdout messages
func Red(str string) string {
	return "\033[31m" + str + "\033[0m"
}

func Green(str string) string {
	return "\033[32m" + str + "\033[0m"
}

func Yellow(str string) string {
	return "\033[33m" + str + "\033[0m"
}

func Blue(str string) string {
	return "\033[34m" + str + "\033[0m"
}

func Magenta(str string) string {
	return "\033[35m" + str + "\033[0m"
}

func Cyan(str string) string {
	return "\033[36m" + str + "\033[0m"
}
