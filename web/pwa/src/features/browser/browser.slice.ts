import {Folder, FormDefinition} from "core-js-api-client";
import {RootState} from "../../app/store";
import {selectByFolderOrDBId} from "../../reducers/form";
import {selectByParentId, selectDatabaseRootFolders} from "../../reducers/folder";
import {databaseGlobalSelectors} from "../../reducers/database";

export const selectChildFolders = (dbOrFolderId?: string) => (state: RootState): Folder[] => {
    if (!dbOrFolderId) {
        return []
    }
    const db = databaseGlobalSelectors.selectById(state, dbOrFolderId)
    if (db) {
        return selectDatabaseRootFolders(state, db.id)
    }
    return selectByParentId(state, dbOrFolderId)
}

export const selectChildForms = (dbOrFolderId?: string) => (state: RootState): FormDefinition[] => {
    return selectByFolderOrDBId(state, dbOrFolderId)
}


