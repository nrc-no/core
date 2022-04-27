import { FieldDefinition, FormDefinition, FormType } from 'core-api-client';

export const makeForm = (
  index: number,
  type: FormType,
  fields: FieldDefinition[],
): FormDefinition => ({
  id: `form${index}`,
  code: '',
  name: `form ${index}`,
  databaseId: 'databaseId',
  folderId: 'folderId',
  formType: type,
  fields,
});
