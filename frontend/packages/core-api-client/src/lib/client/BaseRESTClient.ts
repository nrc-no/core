import axios, {
  AxiosError,
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  Method,
} from 'axios';

import { RequestOptions, Response } from '../types/client/utils';
import { clientResponse } from '../utils/responses';

export class BaseRESTClient {
  protected readonly axiosInstance: AxiosInstance;

  private token = '';

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({
      baseURL,
    });
    this.axiosInstance.interceptors.request.use((value: AxiosRequestConfig) => {
      if (!this.getToken()) {
        return value;
      }
      return {
        ...value,
        headers: {
          ...value.headers,
          Authorization: `Bearer ${this.getToken()}`,
        },
      };
    });
  }

  public getToken = (): string => this.token;

  public setToken = (token: string): void => {
    this.token = token;
  };

  // Deprecated
  public async do<TRequest, TBody>(
    request: TRequest,
    url: string,
    method: Method,
    data: unknown,
    expectStatusCode: number,
    options?: RequestOptions,
  ): Promise<Response<TRequest, TBody>> {
    const headers: { [key: string]: string } = options?.headers ?? {
      Accept: 'application/json',
    };
    let value: AxiosResponse<TBody> | AxiosError<TBody>;
    try {
      value = await this.axiosInstance.request<TBody>({
        responseType: 'json',
        method,
        url,
        data,
        headers,
        withCredentials: true,
      });
    } catch (err) {
      value = err as AxiosError<TBody>;
    }
    return clientResponse<TRequest, TBody>(value, request, expectStatusCode);
  }

  private async makeRequest<TRequestBody, TResponseBody>(
    url: string,
    method: Method,
    data: TRequestBody,
    expectedStatusCode: number,
    options?: RequestOptions,
  ): Promise<Response<TRequestBody, TResponseBody>> {
    let value: AxiosResponse<TResponseBody> | AxiosError<TResponseBody>;
    try {
      const headers: { [key: string]: string } = options?.headers ?? {
        Accept: 'application/json',
      };
      value = await this.axiosInstance.request<TResponseBody>({
        responseType: 'json',
        method,
        url,
        data,
        headers,
        withCredentials: true,
      });
    } catch (err) {
      value = err as AxiosError<TResponseBody>;
    }
    return clientResponse<TRequestBody, TResponseBody>(
      value,
      data,
      expectedStatusCode,
    );
  }

  public async get<TResponseBody>(
    url: string,
  ): Promise<Response<undefined, TResponseBody>> {
    return this.makeRequest(url, 'GET', undefined, 200);
  }

  public async post<TRequestBody, TResponseBody>(
    url: string,
    data: TRequestBody,
  ): Promise<Response<TRequestBody, TResponseBody>> {
    return this.makeRequest(url, 'POST', data, 201);
  }

  public async put<TRequestBody, TResponseBody>(
    url: string,
    data: TRequestBody,
  ): Promise<Response<TRequestBody, TResponseBody>> {
    return this.makeRequest(url, 'PUT', data, 200);
  }

  public async delete(url: string): Promise<Response<undefined, undefined>> {
    return this.makeRequest(url, 'DELETE', undefined, 204);
  }
}
