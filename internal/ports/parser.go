package ports

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const ( MinPort = 1; MaxPort = 65535 )

func ParseExpression(expression string) ([]int, error) {
	if strings.TrimSpace(expression) == "" { return nil, fmt.Errorf("port expression cannot be empty") }
	seen := make(map[int]struct{})
	for _, rawItem := range strings.Split(expression, ",") {
		item := strings.TrimSpace(rawItem); if item == "" { return nil, fmt.Errorf("empty port item in expression") }
		parts := strings.Split(item, "-")
		switch len(parts) {
		case 1:
			port, err := parsePort(parts[0]); if err != nil { return nil, err }; seen[port] = struct{}{}
		case 2:
			start, err := parsePort(parts[0]); if err != nil { return nil, err }; end, err := parsePort(parts[1]); if err != nil { return nil, err }
			if start > end { return nil, fmt.Errorf("invalid port range %d-%d: start is greater than end", start, end) }
			for port := start; port <= end; port++ { seen[port] = struct{}{} }
		default: return nil, fmt.Errorf("invalid port item %q: expected PORT or START-END", item)
		}
	}
	result := make([]int, 0, len(seen)); for port := range seen { result = append(result, port) }; sort.Ints(result); return result, nil
}

func parsePort(value string) (int, error) {
	value = strings.TrimSpace(value); if value == "" { return 0, fmt.Errorf("port cannot be empty") }
	port, err := strconv.Atoi(value); if err != nil { return 0, fmt.Errorf("invalid port %q: must be an integer", value) }
	if port < MinPort || port > MaxPort { return 0, fmt.Errorf("port %d is outside valid range %d-%d", port, MinPort, MaxPort) }
	return port, nil
}
