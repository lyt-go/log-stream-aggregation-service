package service

import (
	"sort"
	"time"
)

type StreamStats struct {
	StreamID    string `json:"stream_id"`
	StreamName  string `json:"stream_name"`
	EntryCount  int    `json:"entry_count"`
	DebugCount  int    `json:"debug_count"`
	InfoCount   int    `json:"info_count"`
	WarnCount   int    `json:"warn_count"`
	ErrorCount  int    `json:"error_count"`
}

type LevelStats struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type TopStream struct {
	StreamID   string `json:"stream_id"`
	StreamName string `json:"stream_name"`
	EntryCount int    `json:"entry_count"`
}

type OverviewStats struct {
	StreamCount       int `json:"stream_count"`
	CollectorCount    int `json:"collector_count"`
	LogEntryCount     int `json:"log_entry_count"`
	SearchRecordCount int `json:"search_record_count"`
	AlertRuleCount    int `json:"alert_rule_count"`
	ActiveStreamCount int `json:"active_stream_count"`
	ActiveCollectorCount int `json:"active_collector_count"`
	ActiveAlertRuleCount int `json:"active_alert_rule_count"`
}

func (s *Service) StatsByStream() ([]StreamStats, error) {
	streams := s.store.ListStreams()
	entries := s.store.ListLogEntries()
	statsMap := make(map[string]*StreamStats)
	for _, st := range streams {
		statsMap[st.ID] = &StreamStats{
			StreamID:   st.ID,
			StreamName: st.Name,
		}
	}
	for _, e := range entries {
		if stat, ok := statsMap[e.StreamID]; ok {
			stat.EntryCount++
			switch e.Level {
			case "debug":
				stat.DebugCount++
			case "info":
				stat.InfoCount++
			case "warn":
				stat.WarnCount++
			case "error":
				stat.ErrorCount++
			}
		}
	}
	result := make([]StreamStats, 0, len(statsMap))
	for _, v := range statsMap {
		result = append(result, *v)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].EntryCount > result[j].EntryCount
	})
	return result, nil
}

func (s *Service) StatsByLevel() ([]LevelStats, error) {
	entries := s.store.ListLogEntries()
	counts := make(map[string]int)
	for _, e := range entries {
		counts[e.Level]++
	}
	result := make([]LevelStats, 0, len(counts))
	for level, count := range counts {
		result = append(result, LevelStats{Level: level, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result, nil
}

func (s *Service) TopStreams(n int) ([]TopStream, error) {
	all, err := s.StatsByStream()
	if err != nil {
		return nil, err
	}
	if n > len(all) {
		n = len(all)
	}
	result := make([]TopStream, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, TopStream{
			StreamID:   all[i].StreamID,
			StreamName: all[i].StreamName,
			EntryCount: all[i].EntryCount,
		})
	}
	return result, nil
}

func (s *Service) Overview() (*OverviewStats, error) {
	streams := s.store.ListStreams()
	collectors := s.store.ListCollectors()
	entries := s.store.ListLogEntries()
	searchRecords := s.store.ListSearchRecords()
	alertRules := s.store.ListAlertRules()
	ov := &OverviewStats{
		StreamCount:       len(streams),
		CollectorCount:    len(collectors),
		LogEntryCount:     len(entries),
		SearchRecordCount: len(searchRecords),
		AlertRuleCount:    len(alertRules),
	}
	for _, st := range streams {
		if st.Status == "active" {
			ov.ActiveStreamCount++
		}
	}
	for _, c := range collectors {
		if c.Status == "active" {
			ov.ActiveCollectorCount++
		}
	}
	for _, a := range alertRules {
		if a.Status == "active" {
			ov.ActiveAlertRuleCount++
		}
	}
	return ov, nil
}

func (s *Service) SearchLogs(query string, startTime, endTime time.Time, page, size int) (interface{}, int, error) {
	all := s.store.ListLogEntries()
	matched := make([]interface{}, 0)
	for _, e := range all {
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		if query != "" {
			found := false
			if containsString(e.Message, query) {
				found = true
			}
			for _, tag := range e.Tags {
				if containsString(tag, query) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		matched = append(matched, e)
	}
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []interface{}{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && len(substr) > 0 && indexOfSubstr(s, substr) >= 0
}

func indexOfSubstr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
