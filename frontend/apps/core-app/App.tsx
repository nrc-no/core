import * as React from 'react';
import * as WebBrowser from 'expo-web-browser';
import { NativeBaseProvider } from 'native-base';
import { theme } from 'core-design-system';
import { NavigationContainer } from '@react-navigation/native';
import {
  // eslint-disable-next-line camelcase
  Roboto_400Regular,
  // eslint-disable-next-line camelcase
  Roboto_400Regular_Italic,
  // eslint-disable-next-line camelcase
  Roboto_500Medium,
  // eslint-disable-next-line camelcase
  Roboto_700Bold,
  useFonts,
} from '@expo-google-fonts/roboto';
import Constants from 'expo-constants';

import { AuthWrapper } from './src/components/AuthWrapper';
import { formsClient } from './src/clients/formsClient';
import { RootNavigator } from './src/navigators';

WebBrowser.maybeCompleteAuthSession();

const App: React.FC = () => {
  const [fontsLoaded] = useFonts({
    Roboto_400Regular,
    Roboto_400Regular_Italic,
    Roboto_500Medium,
    Roboto_700Bold,
  });

  const linking = {
    prefixes: [
      `https://${Constants.manifest?.scheme}.com`,
      `${Constants.manifest?.scheme}://`,
    ],
    config: {
      screens: {
        Recipients: 'recipients',
        RecipientList: '/',
        RecipientRegistration: 'recipients/register',
        RecipientProfile: {
          path: '/:id',
          parse: {
            id: String,
          },
        },
      },
    },
  };

  return (
    fontsLoaded && (
      <NativeBaseProvider theme={theme}>
        <NavigationContainer linking={linking}>
          <AuthWrapper onTokenChange={formsClient.setToken}>
            <RootNavigator />
          </AuthWrapper>
        </NavigationContainer>
      </NativeBaseProvider>
    )
  );
};

export default App;
