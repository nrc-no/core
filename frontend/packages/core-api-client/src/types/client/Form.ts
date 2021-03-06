import { FormDefinition, FormDefinitionList } from '../model/Form';

import { DataOperation, PartialObjectWrapper, Response } from './utils';

export type FormListRequest = Record<string, never> | undefined;
export type FormListResponse = Response<FormListRequest, FormDefinitionList>;

export type FormGetRequest = { id: string };
export type FormGetResponse = Response<FormGetRequest, FormDefinition>;

export type FormCreateRequest = PartialObjectWrapper<FormDefinition>;
export type FormCreateResponse = Response<FormCreateRequest, FormDefinition>;

export type FormLookup = { databaseId: string; formId: string };

export interface FormClientDefinition {
  create: DataOperation<FormCreateRequest, FormCreateResponse>;
  list: DataOperation<FormListRequest, FormListResponse>;
  get: DataOperation<FormGetRequest, FormGetResponse>;
}
