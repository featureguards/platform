import { AxiosError } from 'axios';

import { SerializedError } from '@reduxjs/toolkit';

export const SerializeError = (err: AxiosError): SerializedError => {
  return {
    name: err.name,
    // Envoy converts grpc-message into message field in response data.
    message: err.response?.data.message || err.message,
    code: String(err.response?.status),
    stack: err.stack
  };
};
