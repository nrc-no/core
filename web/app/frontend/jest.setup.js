import 'react-native-gesture-handler/jestSetup';
import mockAsyncStorage from '@react-native-async-storage/async-storage/jest/async-storage-mock';
import 'core-js-api-client';

jest.mock('@react-native-async-storage/async-storage', () => mockAsyncStorage);

jest.mock('expo-secure-store');

jest.mock('react-native-crypto-js', () => {
    const CryptoJS = jest.requireActual('react-native-crypto-js');
    return {
        ...CryptoJS,
        AES: {
            ...CryptoJS.AES,
            encrypt: jest.fn(CryptoJS.AES.encrypt),
            decrypt: jest.fn(CryptoJS.AES.decrypt),
        }
    };
});

jest.mock('react-native-reanimated', () => {
    const Reanimated = require('react-native-reanimated/mock');

    // The mock for `call` immediately calls the callback which is incorrect
    // So we override it with a no-op
    Reanimated.default.call = () => {
    };

    return Reanimated;
});

// Silence the warning: Animated: `useNativeDriver` is not supported because the native animated module is missing
jest.mock('react-native/Libraries/Animated/NativeAnimatedHelper');

jest.mock('core-js-api-client', () => {
    const {FormDefinition} = jest.requireActual('core-js-api-client/lib/types/types');
    const formsList = {response: {items: new Array(10).map(() => new FormDefinition())}};

    function client(host) {
        this.host = host;
        this.listForms = () => new Promise((resolve) => {
            resolve(formsList);
        });

    }

    return { client };
});
