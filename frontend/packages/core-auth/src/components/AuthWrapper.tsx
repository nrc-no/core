import React, { useCallback, useEffect, useMemo, useState } from 'react';

import Browser from '../types/browser';
import exchangeCodeAsync from '../utils/exchangeCodeAsync';
import useDiscovery from '../hooks/useDiscovery';
import useAuthRequest from '../hooks/useAuthRequest';
import {
  AuthWrapperProps,
  CodeChallengeMethod,
  ResponseType,
} from '../types/types';
import { TokenResponse } from '../types/response';

const AuthWrapper: React.FC<AuthWrapperProps> = ({
  children,
  scopes = [],
  clientId,
  issuer,
  redirectUri,
  customLoginComponent,
  handleLoginErr = () => {},
  onTokenChange,
  injectToken = 'access_token',
}) => {
  const browser = useMemo(() => new Browser(), []);

  browser.maybeCompleteAuthSession();

  const discovery = useDiscovery(issuer);

  const [tokenResponse, setTokenResponse] = useState<TokenResponse>();

  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const [request, response, promptAsync] = useAuthRequest(
    {
      clientId,
      usePKCE: true,
      responseType: ResponseType.Code,
      codeChallengeMethod: CodeChallengeMethod.S256,
      scopes,
      redirectUri,
    },
    discovery,
    browser,
  );

  // Make initial token request
  useEffect(() => {
    (async () => {
      if (
        !discovery ||
        request?.codeVerifier == null ||
        !response ||
        response?.type !== 'success'
      )
        return;

      const exchangeConfig = {
        code: response.params.code,
        clientId,
        redirectUri,
        extraParams: {
          code_verifier: request?.codeVerifier,
        },
      };

      try {
        const tr = await exchangeCodeAsync(exchangeConfig, discovery);
        setTokenResponse(tr);
      } catch {
        setTokenResponse(undefined);
      }
    })();
  }, [request?.codeVerifier, response, discovery]);

  // Trigger onTokenChange callback when TokenResponse is received
  useEffect(() => {
    if (!tokenResponse) return;

    const token = (() => {
      switch (injectToken) {
        case 'access_token':
          return tokenResponse?.accessToken ?? '';
        case 'id_token':
          return tokenResponse?.idToken ?? '';
        default:
          return '';
      }
    })();

    onTokenChange(token);
  }, [tokenResponse?.accessToken, tokenResponse?.idToken]);

  // Update logged in status accordingly
  useEffect(() => {
    if (tokenResponse) {
      if (!isLoggedIn) {
        setIsLoggedIn(true);
      }
    } else if (isLoggedIn) {
      setIsLoggedIn(false);
    }
  }, [tokenResponse, isLoggedIn]);

  // Each second, check if token is fresh
  // If not, refresh the token
  const refreshTokenInterval = React.useRef<number | null>(null);
  React.useEffect(() => {
    const refreshToken = async () => {
      if (!discovery) return;
      if (!tokenResponse) return;
      if (!tokenResponse.shouldRefresh()) return;

      const refreshConfig = {
        clientId,
        scopes,
        extraParams: {},
      };

      try {
        const resp = await tokenResponse.refreshAsync(refreshConfig, discovery);
        setTokenResponse(resp);
      } catch (err) {
        setTokenResponse(undefined);
        throw err;
      }
    };

    if (refreshTokenInterval.current)
      window.clearInterval(refreshTokenInterval.current);

    if (
      tokenResponse &&
      tokenResponse.expiresIn &&
      tokenResponse.expiresIn > 0
    ) {
      refreshTokenInterval.current = window.setInterval(refreshToken, 1000);
    }

    return () => {
      if (refreshTokenInterval.current)
        clearInterval(refreshTokenInterval.current);
    };
  }, [tokenResponse?.refreshToken, tokenResponse?.expiresIn]);

  const handleLogin = useCallback(() => {
    promptAsync().catch((err) => {
      handleLoginErr(err);
    });
  }, [discovery, request, promptAsync]);

  if (!isLoggedIn) {
    return (
      <>
        {customLoginComponent ? (
          customLoginComponent({ login: handleLogin })
        ) : (
          <button type="button" onClick={handleLogin}>
            Login
          </button>
        )}
      </>
    );
  }

  return <>{children}</>;
};

export default AuthWrapper;
