import { FieldPath } from 'react-hook-form';
import { FieldDefinition, FieldKind, FormType } from 'core-api-client';
import { RegisterOptions } from 'react-hook-form/dist/types/validator';

import { Form, FormField, ValidationForm } from '../../reducers/Former/types';

export const validationConstants = {
  name: {
    minLength: 3,
    maxLength: 128,
    pattern: /[A-Za-z0-9]/,
  },
  options: {
    min: 2,
    max: 60,
    name: {
      minLength: 2,
      maxLength: 128,
    },
  },
  fields: {
    min: 1,
    max: 100,
  },
};

export const registeredValidation = {
  name: {
    minLength: {
      value: validationConstants.name.minLength,
      message: `Form name needs to be at least ${validationConstants.name.minLength} characters long`,
    },
    maxLength: {
      value: validationConstants.name.maxLength,
      message: `Form name can be at most ${validationConstants.name.maxLength} characters long`,
    },
    pattern: {
      value: validationConstants.name.pattern,
      message: 'Form name contains invalid characters',
    },
    required: { value: true, message: 'Form name is a required field' },
  },
  selectedField: {
    name: {
      minLength: {
        value: validationConstants.name.minLength,
        message: `Field name needs to be at least ${validationConstants.name.minLength} characters long`,
      },
      maxLength: {
        value: validationConstants.name.maxLength,
        message: `Field name can be at most ${validationConstants.name.maxLength} characters long`,
      },
      required: { value: true, message: 'Field name is required' },
    },
    fieldType: {
      reference: {
        databaseId: {
          required: { value: true, message: 'Data base is required' },
        },
        formId: {
          required: { value: true, message: 'Form is required' },
        },
      },
      select: {
        option: {
          name: {
            minLength: {
              value: validationConstants.name.minLength,
              message: `Option name needs to be at least ${validationConstants.name.minLength} characters long`,
            },
            maxLength: {
              value: validationConstants.name.maxLength,
              message: `Option name can be at most ${validationConstants.name.maxLength} characters long`,
            },
            required: { value: true, message: 'Option name is required' },
          },
        },
      },
    },
  },
  values: (field: FieldDefinition): RegisterOptions => {
    const rules: RegisterOptions = {
      required: {
        value: field.required,
        message: 'This field is required',
      },
    };
    if (
      field.fieldType === FieldKind.Date ||
      field.fieldType === FieldKind.Month ||
      field.fieldType === FieldKind.Week
    ) {
      rules.valueAsDate = true;
    }
    if (field.fieldType === FieldKind.Quantity) {
      rules.valueAsNumber = true;
    }
    if (field.fieldType === FieldKind.Month) {
      rules.pattern = {
        value: /^(?:19|20|21)\d{2}-[01]\d$/,
        message: 'wrong pattern',
      };
    }
    if (field.fieldType === FieldKind.Week) {
      rules.pattern = {
        value: /^(?:19|20|21)\d{2}-W[0-5]\d$/,
        message: 'wrong pattern',
      };
    }
    return rules;
  },
};

type CustomError = {
  field: FieldPath<ValidationForm>;
  message: string;
};

export const customValidation = {
  form: (form: Form): CustomError[] => {
    const errors = [];
    if (form.fields.length < validationConstants.fields.min) {
      errors.push({
        field: 'fields' as const,
        message: `Form needs to have at least ${validationConstants.fields.min} field`,
      });
    }
    if (form.fields.length > validationConstants.fields.max) {
      errors.push({
        field: 'fields' as const,
        message: `Form can have at most ${validationConstants.fields.max} fields`,
      });
    }
    if (form.formType === FormType.RecipientFormType) {
      const keyFields = form.fields.filter((field) => {
        return field.key;
      });
      if (keyFields.length !== 1) {
        errors.push({
          field: 'fields' as const,
          message: 'Form needs to have exactly 1 key field',
        });
      } else if (keyFields[0].fieldType !== FieldKind.Reference) {
        errors.push({
          field: 'fields' as const,
          message: 'Key field needs to be a reference',
        });
      }
    }
    return errors;
  },
  selectedField: (field: FormField): CustomError[] => {
    const errors = [];
    if (
      field.fieldType === FieldKind.SingleSelect ||
      field.fieldType === FieldKind.MultiSelect
    ) {
      if (field.options.length < validationConstants.options.min) {
        errors.push({
          field:
            `selectedField.fieldType.${field.fieldType}.options` as FieldPath<ValidationForm>,
          message: `At least ${validationConstants.options.min} options are required`,
        });
      }
      if (field.options.length > validationConstants.options.max) {
        errors.push({
          field:
            `selectedField.fieldType.${field.fieldType}.options` as FieldPath<ValidationForm>,
          message: `At most ${validationConstants.options.max} options are allowed`,
        });
      }
    }
    if (field.fieldType === FieldKind.SubForm) {
      if (field.required) {
        errors.push({
          field: 'selectedField.fieldType.subForm' as const,
          message: 'Subforms cannot be required',
        });
      }
    }
    return errors;
  },
};
