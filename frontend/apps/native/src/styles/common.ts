import { StyleSheet } from 'react-native';

import theme from '../constants/theme';

export default StyleSheet.create({
  textCentered: {
    alignSelf: 'center',
  },
  darkBackground: {
    backgroundColor: theme.colors.backdrop,
    color: theme.colors.text,
  },
});
