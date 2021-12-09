import React from 'react';
import { Path } from 'react-native-svg';

import { IconVariants } from '../../types/icons';

export default (variant: IconVariants) => {
  return (
    <Path
      fill={variant}
      d="M25.6,23.6c-0.7,0-1.3,0.2-1.8,0.6L19,20.7c0.1-0.4,0.1-0.9,0-1.3l4.9-3.5c0.5,0.4,1.1,0.6,1.8,0.6c1.6,0,2.8-1.3,2.8-2.8
	s-1.3-2.8-2.8-2.8s-2.8,1.3-2.8,2.8c0,0.3,0,0.5,0.1,0.8l-4.6,3.3c-0.7-0.9-1.8-1.5-3-1.5c-2.1,0-3.8,1.7-3.8,3.8s1.7,3.8,3.8,3.8
	c1.2,0,2.3-0.6,3-1.5l4.6,3.3c-0.1,0.2-0.1,0.5-0.1,0.8c0,1.6,1.3,2.8,2.8,2.8s2.8-1.3,2.8-2.8S27.2,23.6,25.6,23.6z M25.6,12.4
	c0.7,0,1.2,0.5,1.2,1.2s-0.5,1.2-1.2,1.2s-1.2-0.5-1.2-1.2S25,12.4,25.6,12.4z M15.3,22.1c-1.1,0-2.1-0.9-2.1-2.1s0.9-2.1,2.1-2.1
	s2.1,0.9,2.1,2.1S16.4,22.1,15.3,22.1z M25.6,27.6c-0.7,0-1.2-0.5-1.2-1.2s0.5-1.2,1.2-1.2s1.2,0.5,1.2,1.2S26.3,27.6,25.6,27.6z"
    />
  );
};
