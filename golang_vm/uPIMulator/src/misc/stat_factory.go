package misc

import (
	"fmt"
	"slices"
)

type StatFactory struct {
	name  string
	stats map[string]int64
}

func (this *StatFactory) Init(name string) {
	this.name = name
	this.stats = make(map[string]int64)
}

func (this *StatFactory) Name() string {
	return this.name
}

func (this *StatFactory) Stats() []string {
	stats := make([]string, 0)
	for stat, _ := range this.stats {
		stats = append(stats, stat)
	}

	slices.Sort(stats)
	return stats
}

func (this *StatFactory) Value(stat string) int64 {
	return this.stats[stat]
}

func (this *StatFactory) Increment(stat string, value int64) {
	this.stats[stat] += value
}

func (this *StatFactory) ToLines() []string {
	lines := make([]string, 0)
	for stat, value := range this.stats {
		line := fmt.Sprintf("%s_%s: %d", this.name, stat, value)
		lines = append(lines, line)
	}
	return lines
}
