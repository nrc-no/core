import React from 'react';
import { Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return (
    <Path
      fill={variant}
      d="M29.8,11.8H10.2c-0.4,0-0.8,0.3-0.8,0.8v15c0,0.4,0.3,0.8,0.8,0.8h19.5c0.4,0,0.8-0.3,0.8-0.8v-15
	C30.5,12.1,30.2,11.8,29.8,11.8z M28.8,26.6H11.2v-0.9l3.2-3.9l3.5,4.2l5.5-6.5l5.4,6.4V26.6z M28.8,23.5l-5.2-6.2
	c-0.1-0.1-0.2-0.1-0.3,0L18,23.6l-3.4-4c-0.1-0.1-0.2-0.1-0.3,0l-3.1,3.7v-9.9h17.6V23.5z M15.1,18.7c0.3,0,0.5-0.1,0.8-0.2
	c0.3-0.1,0.5-0.3,0.7-0.4c0.2-0.2,0.3-0.4,0.4-0.7c0.1-0.3,0.2-0.5,0.2-0.8s-0.1-0.5-0.2-0.8c-0.1-0.3-0.3-0.5-0.4-0.7
	c-0.2-0.2-0.4-0.3-0.7-0.4c-0.3-0.1-0.5-0.2-0.8-0.2c-0.5,0-1.1,0.2-1.5,0.6s-0.6,0.9-0.6,1.5s0.2,1.1,0.6,1.5
	C14.1,18.5,14.6,18.7,15.1,18.7z M15.1,16c0.4,0,0.7,0.3,0.7,0.7s-0.3,0.7-0.7,0.7s-0.7-0.3-0.7-0.7S14.8,16,15.1,16z"
    />
  );
};
