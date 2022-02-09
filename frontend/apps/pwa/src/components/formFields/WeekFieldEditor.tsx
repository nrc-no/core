import React, { FC } from 'react';

import { FieldEditorProps } from './types';

export const WeekFieldEditor: FC<FieldEditorProps> = ({
  field,
  value,
  onChange,
  register,
  errors,
}) => {
  const expectedLength = 8;
  const registerObject = register(`values.${field.id}`, {
    required: { value: field.required, message: 'This field is required' },
    pattern: {
      value: /^(?:19|20|21)\d{2}-W[0-5]\d$/,
      message: 'wrong pattern',
    },
    maxLength: { value: expectedLength, message: 'Value is too long' },
    valueAsDate: { value: true, message: 'not a date' },
  });
  return (
    <input
      className={`form-control bg-dark text-light border-secondary ${
        errors?.values && errors?.values[field.id] ? 'is-invalid' : ''
      }`}
      type="week"
      name={field.name}
      placeholder="2021-W52"
      id={field.id}
      value={value || ''}
      {...registerObject}
      onChange={(event) => {
        onChange(event.target.value);
        return registerObject.onChange(event);
      }}
      aria-describedby="errorMessages"
    />
  );
};
