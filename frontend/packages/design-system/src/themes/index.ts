import { extendTheme } from 'native-base';

import tokens from '../tokens';

import Button from './Button';
import Icon from './Icon';
import Text from './Text';

export default extendTheme({
  colors: tokens.colors,
  fontConfig: tokens.fontConfig,
  fontSizes: tokens.fontSizes,
  fonts: tokens.fonts,
  lineHeights: tokens.lineHeights,
  letterSpacings: tokens.letterSpacings,
  components: {
    Button,
    Icon,
    Text,
  },
  space: {
    '1': tokens.spacing.spacing5,
    '2': tokens.spacing.spacing10,
    '3': tokens.spacing.spacing15,
    '4': tokens.spacing.spacing20,
    '5': tokens.spacing.spacing25,
    '6': tokens.spacing.spacing30,
    '8': tokens.spacing.spacing40,
  },
  sizes: {
    '1': tokens.spacing.spacing5,
    '2': tokens.spacing.spacing10,
    '3': tokens.spacing.spacing15,
    '4': tokens.spacing.spacing20,
    '5': tokens.spacing.spacing25,
    '6': tokens.spacing.spacing30,
    '8': tokens.spacing.spacing40,
  },
});