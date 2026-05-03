package membership

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/1204244136/MDA/agent/go-service/pkg/i18n"
	"github.com/rs/zerolog/log"
)

// memberMarker is the emoji marker in task labels that indicates a member-only task.
// Tasks whose label contains this marker require membership (Gold, UserLevel >= 3).
const memberMarker = "🍊"

// unsupportedTiers are membership tiers no longer supported in MDA.
// When detected, a log warning is issued but the user is still processed at their actual level.
var unsupportedTiers = map[string]bool{
	"铜Doro会员": true,
	"银Doro会员": true,
}

// membershipLevels maps tier names to their numeric user level.
var membershipLevels = map[string]int{
	"普通用户":      0,
	"铜Doro会员":    1,
	"银Doro会员":    2,
	"金Doro会员":    3,
	"金Doro企业版": 4,
}

// monthlyCost maps tier names to their monthly cost in ORANGE units.
var monthlyCost = map[string]float64{
	"普通用户":      0,
	"铜Doro会员":    1,
	"银Doro会员":    3,
	"金Doro会员":    5,
	"金Doro企业版": 100,
}

// minMemberLevel is the minimum UserLevel required for member-only tasks in MDA.
const minMemberLevel = 3 // Gold tier

// MemberDataURL is the only data source for V6 membership data.
const MemberDataURL = "https://doropay.top/api/members/v6"

// memberOnlyEntries is populated at init by scanning task JSON files.
var memberOnlyEntries = map[string]bool{}

var loadOnce sync.Once

// taskFile represents the structure of a task JSON file.
type taskFile struct {
	Task []struct {
		Name  string `json:"name"`
		Entry string `json:"entry"`
		Label string `json:"label"`
	} `json:"task"`
}

// LoadMemberOnlyEntries scans task JSON files and populates memberOnlyEntries
// based on whether the task label contains the memberMarker.
func LoadMemberOnlyEntries() {
	loadOnce.Do(func() {
		entries := scanTaskFiles()
		for k, v := range entries {
			memberOnlyEntries[k] = v
		}
		log.Info().
			Int("count", len(memberOnlyEntries)).
			Str("marker", memberMarker).
			Msg("Loaded member-only entries from task files")
	})
}

// IsMemberOnly checks if a task entry requires membership.
func IsMemberOnly(entry string) bool {
	return memberOnlyEntries[entry]
}

func scanTaskFiles() map[string]bool {
	result := map[string]bool{}

	taskDirs := findTaskDirectories()
	if len(taskDirs) == 0 {
		log.Warn().Msg("No task directories found, membership check will have no member-only entries")
		return result
	}

	for _, dir := range taskDirs {
		entries := scanDirectory(dir)
		for k, v := range entries {
			result[k] = v
		}
	}

	return result
}

func findTaskDirectories() []string {
	// Standard search paths relative to working directory
	cwd, _ := os.Getwd()
	wd := filepath.Clean(cwd)

	candidates := []string{
		filepath.Join(wd, "tasks"),
		filepath.Join(wd, "assets", "tasks"),
		filepath.Join(filepath.Dir(wd), "tasks"),
		filepath.Join(filepath.Dir(wd), "assets", "tasks"),
	}

	var found []string
	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			found = append(found, dir)
		}
	}
	return found
}

// resolveLabel resolves an i18n key (e.g. "$task.Arena.label") to its display text.
// If the key doesn't start with "$", it's returned as-is.
func resolveLabel(label string) string {
	if !strings.HasPrefix(label, "$") {
		return label
	}
	return i18n.T(strings.TrimPrefix(label, "$"))
}

func scanDirectory(dir string) map[string]bool {
	result := map[string]bool{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Warn().Str("dir", dir).Err(err).Msg("Failed to read task directory")
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Scan subdirectories (e.g., preset/)
			subResult := scanDirectory(filepath.Join(dir, entry.Name()))
			for k, v := range subResult {
				result[k] = v
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Debug().Str("file", filePath).Err(err).Msg("Failed to read task file")
			continue
		}

		var tf taskFile
		if err := json.Unmarshal(data, &tf); err != nil {
			log.Debug().Str("file", filePath).Err(err).Msg("Failed to parse task file")
			continue
		}

		for _, t := range tf.Task {
			if t.Entry == "" {
				continue
			}
			resolvedLabel := resolveLabel(t.Label)
			if strings.Contains(resolvedLabel, memberMarker) {
				result[t.Entry] = true
				log.Debug().
					Str("entry", t.Entry).
					Str("label_key", t.Label).
					Str("label_resolved", resolvedLabel).
					Str("file", entry.Name()).
					Msg("Found member-only task")
			}
		}
	}

	return result
}
