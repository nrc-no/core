package store

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type Suite struct {
	suite.Suite
	db          *gorm.DB
	dbFactory   *factory
	dbStore     *databaseStore
	folderStore *folderStore
}

func (d *Suite) SetupSuite() {
	db, err := gorm.Open(sqlite.Dialector{DSN: "file::memory:?cache=shared&_foreign_keys=1"}, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	db = db.Debug()
	if err != nil {
		d.FailNow(err.Error())
	}
	d.db = db
	if err := db.AutoMigrate(&Database{}, &Form{}, &Field{}); !assert.NoError(d.T(), err) {
		d.FailNow(err.Error())
	}
	dbFactory := &factory{
		db: db,
	}
	d.dbFactory = dbFactory
}

func (d *Suite) SetupTest() {
	d.dbStore = &databaseStore{
		createDatabase: func(db *gorm.DB, database *types.Database) error {
			return nil
		},
		deleteDatabase: func(db *gorm.DB, databaseId string) error {
			return nil
		},
		db: d.dbFactory,
	}
	d.folderStore = &folderStore{
		db: d.dbFactory,
	}
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) mustCreateDatabase(ctx context.Context) *types.Database {
	db := &types.Database{
		ID:   uuid.NewV4().String(),
		Name: "my-db",
	}
	if _, err := s.dbStore.Create(ctx, db); !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	return db
}

func (s *Suite) mustCreateFolder(ctx context.Context, databaseId string) *types.Folder {
	return s.mustCreateFolderWithParent(ctx, databaseId, "")
}

func (s *Suite) createFolder(ctx context.Context, databaseId string) (*types.Folder, error) {
	return s.createFolderWithParent(ctx, databaseId, "")
}

func (s *Suite) mustCreateFolderWithParent(ctx context.Context, databaseId, parentId string) *types.Folder {
	folder, err := s.createFolderWithParent(ctx, databaseId, parentId)
	if !assert.NoError(s.T(), err) {
		s.T().FailNow()
	}
	return folder
}

func (s *Suite) createFolderWithParent(ctx context.Context, databaseId, parentId string) (*types.Folder, error) {
	folder := &types.Folder{
		ID:         uuid.NewV4().String(),
		DatabaseID: databaseId,
		Name:       "my-folder",
		ParentID:   parentId,
	}
	return s.folderStore.Create(ctx, folder)
}
