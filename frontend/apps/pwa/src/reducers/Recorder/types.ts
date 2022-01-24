import { Record } from 'core-api-client';
import { EntityState } from '@reduxjs/toolkit';

import { ErrorMessage } from '../../types/errors';

export type RecordMap = { [key: string]: FormValue[] };

export interface FormValue extends Omit<Record, 'databaseId'> {
  // records the sub form field that the record belongs to, if any
  ownerFieldId?: string;
  errors: ErrorMessage | undefined;
}

export interface RecorderState extends EntityState<FormValue> {
  saveError: any;
  selectedRecordId: string;
  baseFormId: string;
  editingValues: { [recordId: string]: { [key: string]: any } };
}
