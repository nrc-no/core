import getPkce from 'oauth-pkce';

import { getQueryParams } from '../utils/queryparams';
import { buildQueryString } from '../utils/buildQueryString';

import Browser from './browser';
import {
  AuthError,
  AuthRequestConfig,
  AuthSessionResult,
  CodeChallengeMethod,
  DiscoveryDocument,
  Prompt,
  ResponseType,
  WebBrowserAuthSessionResult,
  WebBrowserResultType,
} from './types';
import { TokenResponse } from './response';
import { CodedError } from './error';

let _authLock = false;

export class AuthRequest {
  public state: string;

  public url = '';

  public codeVerifier = '';

  public codeChallenge = '';

  public responseType: ResponseType;

  public clientId = '';

  public extraParams: Record<string, string>;

  public usePKCE: boolean;

  public codeChallengeMethod: CodeChallengeMethod;

  public redirectUri = '';

  public scopes: string[];

  public prompt: Prompt;

  public browser: Browser;

  constructor(request: AuthRequestConfig, browser: Browser) {
    this.responseType = request.responseType ?? ResponseType.Code;
    this.clientId = request.clientId;
    this.redirectUri = request.redirectUri;
    this.scopes = request.scopes ? request.scopes : [];
    this.prompt = request.prompt ? request.prompt : Prompt.Default;
    this.browser = browser;
    this.state = request.state ?? this.browser.generateRandom(10);
    this.extraParams = request.extraParams ?? {};
    this.codeChallengeMethod = request.codeChallengeMethod ?? CodeChallengeMethod.S256;
    this.usePKCE = request.usePKCE ?? true;
  }

  public async promptAsync(discoveryDocument: DiscoveryDocument): Promise<AuthSessionResult> {
    const url = await this.makeAuthUrlAsync(discoveryDocument);
    let result: WebBrowserAuthSessionResult;
    try {
      _authLock = true;
      result = await this.browser.openAuthSessionAsync({ url });
    } finally {
      _authLock = false;
    }

    if (result?.type !== WebBrowserResultType.SUCCESS) {
      return { type: result?.type } as AuthSessionResult;
    }

    if ('url' in result) {
      return this.parseReturnUrl((result as any)?.url);
    }
    throw new CodedError('ERR_NO_RESULT_URL', 'url not found in auth result');
  }

  parseReturnUrl(url: string): AuthSessionResult {
    const { params, errorCode } = getQueryParams(url);
    const { state, error = errorCode } = params;

    let parsedError: AuthError | null = null;
    let authentication: TokenResponse | null = null;
    if (state !== this.state) {
      // This is a non-standard error
      parsedError = new AuthError({
        error: 'state_mismatch',
        error_description: 'Cross-Site request verification failed. Cached state and returned state do not match.',
      });
    } else if (error) {
      parsedError = new AuthError({ error, ...params });
    }
    if (params.access_token) {
      authentication = TokenResponse.fromQueryParams(params);
    }

    const result: AuthSessionResult = {
      type: parsedError ? 'error' : 'success',
      error: parsedError,
      url,
      params,
      authentication,

      // Return errorCode for legacy
      errorCode,
    };

    return result;
  }

  async makeAuthUrlAsync(discovery: DiscoveryDocument): Promise<string> {
    const request = await this.getAuthRequestConfigAsync();
    if (!request.state) throw new Error('Cannot make request URL without a valid `state` loaded');

    // Create a query string
    const params: Record<string, string> = {};

    if (request.codeChallenge) {
      params.code_challenge = request.codeChallenge;
    }

    // copy over extra params
    for (const extra in request.extraParams) {
      if (extra in request.extraParams) {
        params[extra] = request.extraParams[extra];
      }
    }

    if (request.usePKCE && request.codeChallengeMethod) {
      params.code_challenge_method = request.codeChallengeMethod;
    }

    if (request.prompt) {
      params.prompt = request.prompt;
    }

    // These overwrite any extra params
    params.redirect_uri = request.redirectUri;
    params.client_id = request.clientId;
    params.response_type = request.responseType!;
    params.state = request.state;

    if (request.scopes?.length) {
      params.scope = request.scopes.join(' ');
    }

    const query = buildQueryString(params);
    // Store the URL for later
    this.url = `${discovery.authorization_endpoint}?${query}`;
    return this.url;
  }

  private async ensureCodeIsSetupAsync(): Promise<void> {
    if (this.codeVerifier) {
      return;
    }
    return new Promise((resolve, reject) => {
      getPkce(43, (error, value) => {
        if (error) {
          reject(error);
        }
        this.codeVerifier = value.verifier;
        this.codeChallenge = value.challenge;
        resolve();
      });
    });
  }

  private async getAuthRequestConfigAsync(): Promise<AuthRequestConfig> {
    if (this.usePKCE) {
      await this.ensureCodeIsSetupAsync();
    }
    return {
      clientId: this.clientId,
      responseType: this.responseType,
      redirectUri: this.redirectUri,
      scopes: this.scopes,
      usePKCE: this.usePKCE,
      codeChallengeMethod: this.codeChallengeMethod,
      codeChallenge: this.codeChallenge,
      prompt: this.prompt,
      extraParams: this.extraParams,
      state: this.state,
    };
  }
}
