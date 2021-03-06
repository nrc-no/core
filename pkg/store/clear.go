package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/sql/convert"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

const clearProtectionEnvVar = "VERY_DANGEROUSLY_ALLOW_TO_CLEAR_DB"

func checkClearProtection(ctx context.Context) error {
	l := logging.NewLogger(ctx)
	requiredEnvVar, ok := os.LookupEnv(clearProtectionEnvVar)
	if len(requiredEnvVar) == 0 || !ok || requiredEnvVar != "yes" {
		l.Error(fmt.Sprintf("Cannot run the Clear method. You must set the %s=yes environment variable\n", clearProtectionEnvVar))
		return errors.New("required environment variable not present")
	}
	return nil
}

// Clear deletes everything from the database. Useful for development.
// Be careful in production...
func Clear(ctx context.Context, db *gorm.DB) error {

	if err := checkClearProtection(ctx); err != nil {
		return err
	}

	l := logging.NewLogger(ctx)

	l.Info("listing databases")
	var databases []*Database
	if err := db.Find(&databases).Error; err != nil {
		l.Error("failed to list databases", zap.Error(err))
		return err
	}

	l.Info("deleting databases")
	for _, database := range databases {
		l.Info("deleting database", zap.String("database_id", database.ID))
		if err := convert.DeleteDatabaseSchemaIfExist(db, database.ID); err != nil {
			l.Error("failed to delete database", zap.Error(err))
			return err
		}
	}

	l.Info("deleting field options")
	if err := db.Where("field_id = field_id").Delete(&Option{}).Error; err != nil {
		l.Error("failed to delete field options", zap.Error(err))
		return err
	}

	l.Info("deleting fields")
	if err := db.Where("id = id").Delete(&Field{}).Error; err != nil {
		l.Error("failed to delete fields", zap.Error(err))
		return err
	}

	l.Info("deleting forms")
	if err := db.Where("id = id").Delete(&Form{}).Error; err != nil {
		l.Error("failed to delete forms", zap.Error(err))
		return err
	}

	l.Info("deleting folders")
	if err := db.Where("id = id").Delete(&Folder{}).Error; err != nil {
		l.Error("failed to delete folders", zap.Error(err))
		return err
	}

	l.Info("deleting databases")
	if err := db.Where("id = id").Delete(&Database{}).Error; err != nil {
		l.Error("failed to delete databases", zap.Error(err))
		return err
	}

	l.Info("deleting identity providers")
	if err := db.Where("id = id").Delete(&IdentityProvider{}).Error; err != nil {
		l.Error("failed to delete identity providers", zap.Error(err))
		return err
	}

	l.Info("deleting organizations")
	if err := db.Where("id = id").Delete(&Organization{}).Error; err != nil {
		l.Error("failed to delete organizations", zap.Error(err))
		return err
	}

	l.Info("successfully cleared database")
	return nil
}
