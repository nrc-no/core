import * as React from 'react';
import {
  Pressable,
  FormControl,
  Input,
  Modal,
  Button,
  HStack,
  Checkbox,
  ChevronDownIcon,
  usePropsResolution,
} from 'native-base';
import { useFormContext, useController } from 'react-hook-form';
import { FormDefinition, Validation } from 'core-api-client';

type Props = {
  formId: string;
  field: FormDefinition['fields'][number];
};

export const MultiSelectFieldInput: React.FC<Props> = ({ formId, field }) => {
  const [open, setOpen] = React.useState(false);
  const [internalValues, setInternalValues] = React.useState<string[]>([]);

  const { control } = useFormContext();

  const {
    field: { onChange, onBlur, value, ref },
    fieldState: { error, invalid },
  } = useController({
    name: `${formId}.${field.id}`,
    control,
    rules: Validation.Record.formValidationRules.field.multiSelect(field),
  });

  React.useEffect(() => {
    if (value) {
      setInternalValues(value);
    }
  }, [JSON.stringify(value)]);

  const handleOpenModal = () => setOpen(true);

  const handleCloseModal = (reset: boolean) => () => {
    setOpen(false);
    if (reset) {
      setInternalValues(value || []);
    }
  };

  const handleAdd = () => {
    onChange(internalValues);
    handleCloseModal(false)();
  };

  const { customDropdownIconProps } = usePropsResolution(
    'Select',
    {},
    {
      isDisabled: false,
      isHovered: false,
      isFocused: false,
      isFocusVisible: false,
    },
    undefined,
  );

  const valueString = value
    ? field.fieldType.multiSelect?.options
        .filter((o) => value.includes(o.id))
        .map((o) => o.name)
        .join(', ')
    : '';

  return (
    <>
      <FormControl isRequired={field.required} isInvalid={invalid}>
        <FormControl.Label>{field.name}</FormControl.Label>
        <Pressable
          testID="multi-select-field-input-modal-toggle-button"
          onPress={handleOpenModal}
        >
          <Input
            testID="multi-select-field-input-value"
            ref={ref}
            editable={false}
            onBlur={onBlur}
            onChangeText={onChange}
            value={valueString}
            autoCompleteType="off"
            InputRightElement={<ChevronDownIcon {...customDropdownIconProps} />}
          />
        </Pressable>
        <FormControl.HelperText>{field.description}</FormControl.HelperText>
        <FormControl.ErrorMessage>{error?.message}</FormControl.ErrorMessage>
      </FormControl>

      <Modal isOpen={open} onClose={handleCloseModal(true)}>
        <Modal.Content>
          <Modal.Header>{field.name}</Modal.Header>
          <Modal.Body>
            <Checkbox.Group onChange={setInternalValues}>
              {field.fieldType.multiSelect?.options.map((option, i) => (
                <Checkbox
                  testID={`multi-select-field-input-option-${i}`}
                  key={option.id}
                  value={option.id}
                >
                  {option.name}
                </Checkbox>
              ))}
            </Checkbox.Group>
          </Modal.Body>
          <Modal.Footer>
            <HStack space={4}>
              <Button
                testID="multi-select-field-input-modal-cancel"
                onPress={handleCloseModal(true)}
                colorScheme="secondary"
                variant="minor"
              >
                Cancel
              </Button>
              <Button
                testID="multi-select-field-input-modal-submit"
                onPress={handleAdd}
                colorScheme="secondary"
                variant="major"
              >
                Add
              </Button>
            </HStack>
          </Modal.Footer>
        </Modal.Content>
      </Modal>
    </>
  );
};
