import React from 'react';
import { useForm } from 'react-hook-form';
import { FormDefinition } from 'core-api-client';

import ViewRecordScreen, {
  ViewRecordScreenProps,
} from '../screens/ViewRecordScreen';
import { ViewRecordScreenContainerProps } from '../../types/screens';
import useApiClient from '../../utils/clients';

export const ViewRecordScreenContainer = ({
  route,
  state,
}: ViewRecordScreenContainerProps) => {
  const { formId, recordId } = route.params;

  const [isLoading, setIsLoading] = React.useState(true);
  const [form, setForm] = React.useState<FormDefinition>();

  const { control, reset } = useForm();

  const apiClient = useApiClient();

  React.useEffect(() => {
    const getForm = async () => {
      try {
        const { response } = await apiClient.Form.get({ id: formId });
        setForm(response);
      } catch (err) {
        console.error(err);
      }
    };
    getForm();
  }, [formId]);

  React.useEffect(() => {
    if (form) {
      reset(state.formsById[formId].recordsById[recordId].values);
      setIsLoading(false);
    }
  }, [form]);

  const viewRecordScreenProps: ViewRecordScreenProps = {
    isLoading,
    form,
    control,
  };
  return <ViewRecordScreen {...viewRecordScreenProps} />;
};
