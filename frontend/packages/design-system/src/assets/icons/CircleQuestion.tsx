import React from 'react';
import { Path, Circle, G } from 'react-native-svg';

import { IconProps } from '../../types/icons';

export const CircleQuestion: React.FC<IconProps['customIconProps']> = ({
  innerColor,
}) => (
  <G>
    <Circle cx="10" cy="10" r="10" />
    <Path
      d="M8.77547 11.9496C8.77547 11.2485 8.8606 10.6901 9.03087 10.2744C9.20115 9.85875 9.51164 9.4506 9.96237 9.04996C10.4181 8.64431 10.7211 8.31628 10.8713 8.06588C11.0216 7.81047 11.0967 7.54254 11.0967 7.26209C11.0967 6.41574 10.7061 5.99256 9.9248 5.99256C9.55421 5.99256 9.25624 6.10775 9.03087 6.33812C8.81052 6.56348 8.69534 6.87648 8.68532 7.27712H6.50684C6.51685 6.32059 6.82484 5.57189 7.43081 5.03102C8.04179 4.49016 8.87312 4.21973 9.9248 4.21973C10.9865 4.21973 11.8103 4.47764 12.3963 4.99346C12.9822 5.50428 13.2752 6.22794 13.2752 7.16444C13.2752 7.59012 13.18 7.99326 12.9897 8.37387C12.7994 8.74947 12.4664 9.16764 11.9906 9.62838L11.3821 10.2068C11.0015 10.5724 10.7837 11.0006 10.7286 11.4914L10.6985 11.9496H8.77547ZM8.55762 14.2558C8.55762 13.9202 8.6703 13.6448 8.89566 13.4295C9.12603 13.2091 9.419 13.0989 9.77456 13.0989C10.1301 13.0989 10.4206 13.2091 10.646 13.4295C10.8763 13.6448 10.9915 13.9202 10.9915 14.2558C10.9915 14.5863 10.8788 14.8593 10.6535 15.0746C10.4331 15.2899 10.1401 15.3976 9.77456 15.3976C9.40898 15.3976 9.11351 15.2899 8.88815 15.0746C8.66779 14.8593 8.55762 14.5863 8.55762 14.2558Z"
      fill={innerColor || 'white'}
    />
  </G>
);
