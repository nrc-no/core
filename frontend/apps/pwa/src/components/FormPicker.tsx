import React, { FC, useCallback } from 'react';
import { FormDefinition, FormType } from 'core-api-client';
import { FieldErrors } from 'react-hook-form';
import { ErrorMessage } from '@hookform/error-message';

import { useDatabase, useForms } from '../app/hooks';
import { ValidationForm } from '../reducers/Former/types';
import { registeredValidation } from '../features/former/validation';

export type FormPickerProps = {
  forms: FormDefinition[];
  formId: string | undefined;
  disabled?: boolean;
  setFormId: (formId: string) => void;
  register: any;
  errors?: FieldErrors<ValidationForm>;
};

export const FormPicker: FC<FormPickerProps> = ({
  forms,
  formId,
  setFormId,
  disabled,
  register,
  errors,
}) => {
  const hasForms = forms.length > 0;

  const registerSelectedFieldReference = register(
    'selectedField.fieldType.reference.formId',
    registeredValidation.selectedField.fieldType.reference.formId,
  );
  return (
    <div>
      <select
        disabled={disabled || !hasForms}
        {...registerSelectedFieldReference}
        onChange={(e) => {
          setFormId(e.target.value);
          return registerSelectedFieldReference.onChange(e);
        }}
        value={formId || ''}
        className={`form-select ${
          errors?.selectedField?.fieldType?.reference?.formId && !disabled
            ? 'is-invalid'
            : ''
        }`}
        aria-label="Select Form"
      >
        <option disabled value="">
          {hasForms ? 'Select Form' : 'No Forms'}
        </option>
        {forms.map((f) => {
          return (
            <option value={f.id} key={f.id}>
              {f.name}
            </option>
          );
        })}
      </select>
      <div className="invalid-feedback" id="errorMessages">
        <ErrorMessage
          errors={errors}
          name="selectedField.fieldType.reference.formId"
        />
      </div>
    </div>
  );
};

export type FormPickerContainerProps = {
  databaseId: string | undefined;
  formId: string | undefined;
  isRecipientKey: boolean;
  // setForm: (form: FormDefinition | undefined) => void;
  setFormId: (formId: string) => void;
  register: any;
  errors?: FieldErrors<ValidationForm>;
};

const FormPickerContainer: FC<FormPickerContainerProps> = ({
  databaseId,
  formId,
  setFormId,
  // setForm,
  isRecipientKey,
  register,
  errors,
}) => {
  const database = useDatabase(databaseId);
  const allForms = useForms({ databaseId });
  const forms = isRecipientKey
    ? allForms.filter((f) => f.formType === FormType.RecipientFormType)
    : allForms;

  const callback = useCallback(
    (fID: string) => {
      setFormId(fID);
      // const form = forms.find((f) => f.id === fID);
      // setForm(form);
    },
    [forms, setFormId],
  );

  return (
    <FormPicker
      disabled={!database}
      setFormId={callback}
      forms={forms}
      formId={formId}
      register={register}
      errors={errors}
    />
  );
};

export default FormPickerContainer;
