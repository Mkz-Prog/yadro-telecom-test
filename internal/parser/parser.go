package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Mkz-Prog/yadro-telecom-test/internal/domain"
)

// ParseLine принимает одну строку лога и превращает её в структуру domain.Event.
func ParseLine(line string) (*domain.Event, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	endBracket := strings.Index(line, "]")
	if !strings.HasPrefix(line, "[") || endBracket == -1 || endBracket < 9 {
		return nil, fmt.Errorf("invalid line format or missing time bracket: %q", line)
	}

	timeRaw := line[:endBracket+1]
	timeStr := line[1:endBracket]

	parsedTime, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time %s: %w", timeStr, err)
	}

	rest := strings.TrimSpace(line[endBracket+1:])
	parts := strings.Fields(rest)
	if len(parts) < 2 {
		return nil, fmt.Errorf("missing player ID or event ID in line: %q", line)
	}

	playerID, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid player ID %q: %w", parts[0], err)
	}

	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid event ID %q: %w", parts[1], err)
	}

	var extraParam string
	if len(parts) > 2 {
		extraParam = strings.Join(parts[2:], " ")
	}

	return &domain.Event{
		Time:       parsedTime,
		TimeRaw:    timeRaw,
		PlayerID:   playerID,
		ID:         eventID,
		ExtraParam: extraParam,
	}, nil
}
