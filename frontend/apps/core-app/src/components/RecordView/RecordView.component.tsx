import * as React from 'react';
import { VStack } from 'native-base';
import { FieldKind } from 'core-api-client';

import { SubformFieldValueComponent } from './SubformFieldValue.component';
import { NormalisedFieldValue } from './RecordView.types';
import { FieldValueComponent } from './FieldValue.component';

type Props = {
  fieldValues: NormalisedFieldValue[];
  hideKeyFields?: boolean;
};

export const RecordViewComponent: React.FC<Props> = ({
  fieldValues,
  hideKeyFields,
}) => {
  return (
    <VStack space={4}>
      {fieldValues.map((f, i) => {
        if (hideKeyFields && f.key) return null;
        if (f.fieldType === FieldKind.SubForm) {
          return (
            <SubformFieldValueComponent
              key={i}
              header={f.header}
              data={f.data}
              columns={f.columns}
            />
          );
        }

        return <FieldValueComponent key={i} item={f} />;
      })}
    </VStack>
  );
};
