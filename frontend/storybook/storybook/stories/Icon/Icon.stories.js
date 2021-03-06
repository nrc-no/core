import { storiesOf } from '@storybook/react-native';
import React from 'react';
import { icons, theme, tokens, Icon } from 'core-design-system';
import { select } from '@storybook/addon-knobs';
import { Box } from 'native-base';

import CenterView from '../CenterView';

storiesOf('Icon', module)
  .addDecorator((getStory) => <CenterView>{getStory()}</CenterView>)
  .add('Icon, Overview', () => {
    const IconNameList = Object.keys(icons);

    return (
      <Box
        style={{ flexWrap: 'wrap', flexDirection: 'initial', width: '250px' }}
      >
        {IconNameList.map((name) => {
          return (
            <Icon
              size="6"
              m="2"
              key={name}
              name={name}
              color={select(
                'color',
                tokens.colors.icons,
                theme.colors.icons.dark,
              )}
              customIconProps={{
                innerColor: select(
                  'inner color',
                  tokens.colors.neutral,
                  theme.colors.neutral['100'],
                ),
              }}
            />
          );
        })}
      </Box>
    );
  })
  .add('Icon', () => {
    const IconNameList = Object.keys(icons);

    return (
      <Icon
        name={select('name', IconNameList, IconNameList[0])}
        color={select('color', tokens.colors.icons, theme.colors.icons.dark)}
      />
    );
  });
