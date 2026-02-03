package cost

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry represents a cached parsing result.
type CacheEntry struct {
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	FileModTime  time.Time `json:"file_mod_time"`
	LastParsed   time.Time `json:"last_parsed"`
	EntriesCount int       `json:"entries_count"`
	TotalCost    float64   `json:"total_cost"`
	Checksum     string    `json:"checksum"`
}

// Cache manages incremental parsing cache.
type Cache struct {
	Entries  map[string]*CacheEntry `json:"entries"`
	CacheDir string                 `json:"-"`
}

// LoadCache loads the cache from disk.
func LoadCache() (*Cache, error) {
	cacheDir := getCacheDir()
	if cacheDir == "" {
		return &Cache{Entries: make(map[string]*CacheEntry)}, nil
	}

	cachePath := filepath.Join(cacheDir, "cost_cache.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return &Cache{
			Entries:  make(map[string]*CacheEntry),
			CacheDir: cacheDir,
		}, nil
	}

	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return &Cache{
			Entries:  make(map[string]*CacheEntry),
			CacheDir: cacheDir,
		}, nil
	}

	cache.CacheDir = cacheDir
	if cache.Entries == nil {
		cache.Entries = make(map[string]*CacheEntry)
	}

	return &cache, nil
}

// Save saves the cache to disk.
func (c *Cache) Save() error {
	if c.CacheDir == "" {
		c.CacheDir = getCacheDir()
	}

	if c.CacheDir == "" {
		return nil
	}

	// Ensure directory exists
	if err := os.MkdirAll(c.CacheDir, 0755); err != nil {
		return err
	}

	cachePath := filepath.Join(c.CacheDir, "cost_cache.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

// IsValid checks if a cache entry is still valid for a file.
func (c *Cache) IsValid(path string) bool {
	entry, ok := c.Entries[path]
	if !ok {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if file has been modified
	if info.Size() != entry.FileSize {
		return false
	}

	if !info.ModTime().Equal(entry.FileModTime) {
		return false
	}

	return true
}

// Get returns a cached entry if valid.
func (c *Cache) Get(path string) (*CacheEntry, bool) {
	if !c.IsValid(path) {
		return nil, false
	}
	return c.Entries[path], true
}

// Set stores a cache entry.
func (c *Cache) Set(path string, entries []Entry) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	var totalCost float64
	for _, e := range entries {
		totalCost += e.CostUSD
	}

	c.Entries[path] = &CacheEntry{
		FilePath:     path,
		FileSize:     info.Size(),
		FileModTime:  info.ModTime(),
		LastParsed:   time.Now(),
		EntriesCount: len(entries),
		TotalCost:    totalCost,
		Checksum:     computeChecksum(path),
	}
}

// ParseWithCache parses a JSONL file using cache for unchanged files.
// Returns (entries, fromCache, error).
func (c *Cache) ParseWithCache(path string) ([]Entry, bool, error) {
	if cached, ok := c.Get(path); ok {
		// Cache hit: return cached summary as a synthetic entry
		// This avoids re-parsing unchanged files for cost aggregation
		return []Entry{{
			Timestamp: cached.LastParsed,
			CostUSD:   cached.TotalCost,
		}}, true, nil
	}

	entries, err := ParseJSONL(path)
	if err != nil {
		return nil, false, err
	}

	c.Set(path, entries)
	return entries, false, nil
}

// GetCachedSummary returns cached summary without parsing.
// Returns (totalCost, entryCount, ok).
func (c *Cache) GetCachedSummary(path string) (float64, int, bool) {
	if cached, ok := c.Get(path); ok {
		return cached.TotalCost, cached.EntriesCount, true
	}
	return 0, 0, false
}

// Cleanup removes cache entries for files that no longer exist.
func (c *Cache) Cleanup() {
	for path := range c.Entries {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			delete(c.Entries, path)
		}
	}
}

// getCacheDir returns the visor cache directory.
func getCacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cache", "visor")
}

// computeChecksum computes a simple checksum for cache validation.
func computeChecksum(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}

	data := fmt.Sprintf("%s:%d:%d", path, info.Size(), info.ModTime().Unix())
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
