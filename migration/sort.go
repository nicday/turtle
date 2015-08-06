package migration

import "sort"

// sortedMigrations implements sort.Sort
type sortedMigrations struct {
	migrations map[string]*Migration
	direction  string
	sorted     []*Migration
}

// Len satisfies the sort.Sort interface.
func (s *sortedMigrations) Len() int {
	return len(s.migrations)
}

// Less satisfies the sort.Sort interface.
func (s *sortedMigrations) Less(i, j int) bool {
	if s.direction == "desc" {
		return s.migrations[s.sorted[i].ID].ID > s.migrations[s.sorted[j].ID].ID
	}
	return s.migrations[s.sorted[i].ID].ID < s.migrations[s.sorted[j].ID].ID
}

// Swap satisfies the sort.Sort interface.
func (s *sortedMigrations) Swap(i, j int) {
	s.sorted[i], s.sorted[j] = s.sorted[j], s.sorted[i]
}

// SortMigrations will sort migrations in either ascending or descending order based on their ID.
func SortMigrations(migrations map[string]*Migration, direction string) []*Migration {
	s := sortedMigrations{}
	s.migrations = migrations
	s.direction = direction
	s.sorted = make([]*Migration, len(migrations))
	i := 0
	for _, val := range migrations {
		s.sorted[i] = val
		i++
	}
	sort.Sort(&s)
	return s.sorted
}
