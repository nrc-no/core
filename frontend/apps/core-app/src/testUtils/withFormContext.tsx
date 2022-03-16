import { Button } from 'native-base';
import React from 'react';
import { FieldValues, FormProvider, useForm } from 'react-hook-form';

export const withFormContext =
  <T,>(
    Component: React.FC<T>,
    defaultValues: FieldValues = {},
    onSubmit: (data: any) => void = () => {},
  ) =>
  // eslint-disable-next-line react/display-name
  (props: T) => {
    const f = useForm({ defaultValues });

    return (
      <FormProvider {...f}>
        <Component {...props} />
        <Button
          testID="with-form-context-submit-button"
          onPress={f.handleSubmit(onSubmit)}
        >
          Submit
        </Button>
      </FormProvider>
    );
  };
