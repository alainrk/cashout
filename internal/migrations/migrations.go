package migrations

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration represents a single database migration
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"size:255;not null;unique"`
	Name      string    `gorm:"size:255;not null"`
	AppliedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// Migrator handles database migrations
type Migrator struct {
	db *gorm.DB
}

// MigrationFunc is a function that performs a migration
type MigrationFunc func(*gorm.DB) error

// MigrationDefinition defines a migration with its metadata
type MigrationDefinition struct {
	Version  string
	Name     string
	Migrate  MigrationFunc
	Rollback MigrationFunc // Optional rollback function
}

// Available migrations - this will be populated by each migration file
var migrations []MigrationDefinition

// RegisterMigration adds a migration to the list of available migrations
func RegisterMigration(version, name string, migrateFn MigrationFunc) {
	migrations = append(migrations, MigrationDefinition{
		Version: version,
		Name:    name,
		Migrate: migrateFn,
	})
}

// RegisterMigrationWithRollback adds a migration with rollback capability
func RegisterMigrationWithRollback(version, name string, migrateFn, rollbackFn MigrationFunc) {
	migrations = append(migrations, MigrationDefinition{
		Version:  version,
		Name:     name,
		Migrate:  migrateFn,
		Rollback: rollbackFn,
	})
}

// NewMigrator creates a new migrator
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// EnsureMigrationTable makes sure the migration tracking table exists
func (m *Migrator) EnsureMigrationTable() error {
	return m.db.AutoMigrate(&Migration{})
}

// GetAppliedMigrations returns all migrations that have been applied
func (m *Migrator) GetAppliedMigrations() ([]Migration, error) {
	var appliedMigrations []Migration
	if err := m.db.Order("version").Find(&appliedMigrations).Error; err != nil {
		return nil, err
	}
	return appliedMigrations, nil
}

// MigrateUp applies all pending migrations
func (m *Migrator) MigrateUp() error {
	// Ensure the migration table exists
	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get all applied migrations
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map of applied migrations for faster lookup
	appliedMap := make(map[string]bool)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Version] = true
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Apply pending migrations
	for _, migration := range migrations {
		if !appliedMap[migration.Version] {
			fmt.Printf("Applying migration %s: %s\n", migration.Version, migration.Name)

			// Run the migration in a transaction
			err := m.db.Transaction(func(tx *gorm.DB) error {
				// Run the migration
				if err := migration.Migrate(tx); err != nil {
					return err
				}

				// Record that the migration has been applied
				return tx.Create(&Migration{
					Version:   migration.Version,
					Name:      migration.Name,
					AppliedAt: time.Now(),
				}).Error
			})
			if err != nil {
				return fmt.Errorf("failed to apply migration '%s': %w", migration.Version, err)
			}

			fmt.Printf("Migration %s applied successfully\n", migration.Version)
		}
	}

	return nil
}

// MigrateDown rolls back the last applied migration
func (m *Migrator) MigrateDown() error {
	// Ensure the migration table exists
	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get all applied migrations
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedMigrations) == 0 {
		fmt.Println("No migrations to roll back")
		return nil
	}

	// Get the last applied migration
	lastApplied := appliedMigrations[len(appliedMigrations)-1]

	// Find the migration in our list
	var migrationToRollback MigrationDefinition
	found := false
	for _, migration := range migrations {
		if migration.Version == lastApplied.Version {
			migrationToRollback = migration
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("could not find migration with version %s to roll back", lastApplied.Version)
	}

	// Check if rollback function exists
	if migrationToRollback.Rollback == nil {
		return fmt.Errorf("migration %s does not support rollback", lastApplied.Version)
	}

	fmt.Printf("Rolling back migration %s: %s\n", lastApplied.Version, lastApplied.Name)

	// Run the rollback in a transaction
	err = m.db.Transaction(func(tx *gorm.DB) error {
		// Run the rollback
		if err := migrationToRollback.Rollback(tx); err != nil {
			return err
		}

		// Remove the migration record
		return tx.Delete(&Migration{}, "version = ?", lastApplied.Version).Error
	})
	if err != nil {
		return fmt.Errorf("failed to roll back migration '%s': %w", lastApplied.Version, err)
	}

	fmt.Printf("Migration %s rolled back successfully\n", lastApplied.Version)
	return nil
}

// Status prints the status of all migrations
func (m *Migrator) Status() error {
	// Ensure the migration table exists
	if err := m.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get all applied migrations
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Create a map of applied migrations for faster lookup
	appliedMap := make(map[string]Migration)
	for _, migration := range appliedMigrations {
		appliedMap[migration.Version] = migration
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Print status
	fmt.Println("Migration Status:")
	fmt.Println("================")
	for _, migration := range migrations {
		applied, exists := appliedMap[migration.Version]
		if exists {
			fmt.Printf("[✓] %s: %s (applied at %s)\n", migration.Version, migration.Name, applied.AppliedAt.Format(time.RFC3339))
		} else {
			fmt.Printf("[ ] %s: %s (pending)\n", migration.Version, migration.Name)
		}
	}

	return nil
}
