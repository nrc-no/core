import React, { FC } from 'react';
import { useFormContext } from 'react-hook-form';

import { FieldEditorProps } from './types';

export const DateFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  errors,
}) => {
  const { register } = useFormContext();

  const registerObject = register(`values.${field.id}`, {
    required: { value: field.required, message: 'This field is required' },
    valueAsDate: true,
  });

  return (
    <input
      className={`form-control bg-dark text-light border-secondary ${
        errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
      }`}
      type="date"
      id={field.id}
      value={value || ''}
      {...registerObject}
      onChange={(event) => {
        onChange(event.target.value);
        return registerObject.onChange(event);
      }}
      aria-describedby={`errorMessages description-${field.id}`}
    />
  );
};
