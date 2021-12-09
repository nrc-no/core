import React from 'react';
import { Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return <Path fill={variant} d="M23.4,20.2l-9-9l1.1-1.1l9.5,9.5l0.5,0.5L25,20.8l-9,9L15,28.7L23.4,20.2z" />;
};
