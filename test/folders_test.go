package test

import (
	"context"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestFoldersApi() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var doneFns []func()
	defer func() {
		for i := len(doneFns) - 1; i >= 0; i-- {
			doneFns[i]()
		}
	}()

	cli := s.PublicClient(ctx)

	db := &types.Database{}
	if err := cli.CreateDatabase(ctx, &types.Database{Name: "test-db"}, db); !assert.NoError(s.T(), err) {
		return
	}
	doneFns = append(doneFns, func() { _ = cli.DeleteDatabase(ctx, db.ID) })

	created := &types.Folder{}
	if err := cli.CreateFolder(ctx, &types.Folder{DatabaseID: db.ID, Name: "folder"}, created); !assert.NoError(s.T(), err) {
		return
	}
	doneFns = append(doneFns, func() { _ = cli.DeleteFolder(ctx, created.ID) })

	list := &types.FolderList{}
	if err := cli.ListFolders(ctx, list); !assert.NoError(s.T(), err) {
		return
	}

	assert.True(s.T(), list.HasFolderWithID(created.ID))

}
