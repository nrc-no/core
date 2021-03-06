import { View } from 'react-native';
import React from 'react';
import { Text, TextInput as TextInputRNP } from 'react-native-paper';

import { darkTheme } from '../../constants/theme';

import { InputProps } from './FormControl';

const TextInput = ({
  fieldDefinition,
  style,
  value,
  onChange,
  onBlur,
  error,
  invalid,
  isTouched,
  isDirty,
  isMultiple,
  isQuantity,
}: InputProps) => {
  return (
    <View style={style}>
      {fieldDefinition.name && <Text theme={darkTheme}>{fieldDefinition.name}</Text>}
      {fieldDefinition.description && (
        <Text theme={darkTheme} style={{ fontSize: 10 }}>
          {fieldDefinition.description}
        </Text>
      )}
      <TextInputRNP
        onChangeText={onChange}
        value={value}
        onBlur={onBlur}
        error={isTouched && isDirty && error}
        multiline={isMultiple}
        keyboardType={isQuantity ? 'numeric' : 'default'}
      />
      {isTouched && isDirty && error && <Text>{error.message == '' ? 'invalid' : error.message}</Text>}
    </View>
  );
};

export default TextInput;
