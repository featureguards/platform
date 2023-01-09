/* tslint:disable */
/* eslint-disable */
/**
 *
 *
 *
 *
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { Configuration } from './configuration';
import globalAxios, { AxiosPromise, AxiosInstance, AxiosRequestConfig } from 'axios';
// Some imports not used depending on template conditions
// @ts-ignore
import {
  DUMMY_BASE_URL,
  assertParamExists,
  setApiKeyToObject,
  setBasicAuthToObject,
  setBearerAuthToObject,
  setOAuthToObject,
  setSearchParams,
  serializeDataIfNeeded,
  toPathString,
  createRequestFunction
} from './common';
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, RequestArgs, BaseAPI, RequiredError } from './base';

/**
 *
 * @export
 * @enum {string}
 */

export const DateTimeOperator = {
  BEFORE: 0,
  AFTER: 1
} as const;

export type DateTimeOperator = typeof DateTimeOperator[keyof typeof DateTimeOperator];

/**
 *
 * @export
 * @enum {string}
 */

export const DynamicSettingType = {
  BOOL: 0,
  STRING: 1,
  INTEGER: 2,
  FLOAT: 3,
  SET: 4,
  MAP: 5,
  LIST: 6,
  JSON: 7
} as const;

export type DynamicSettingType = typeof DynamicSettingType[keyof typeof DynamicSettingType];

/**
 *
 * @export
 * @enum {string}
 */

export const FeatureToggleType = {
  ON_OFF: 0,
  PERCENTAGE: 1
} as const;

export type FeatureToggleType = typeof FeatureToggleType[keyof typeof FeatureToggleType];

/**
 *
 * @export
 * @enum {string}
 */

export const FloatOperator = {
  EQ: 0,
  GT: 1,
  LT: 2,
  GTE: 3,
  LTE: 4,
  NEQ: 5,
  IN: 6
} as const;

export type FloatOperator = typeof FloatOperator[keyof typeof FloatOperator];

/**
 *
 * @export
 * @enum {string}
 */

export const IntOperator = {
  EQ: 0,
  GT: 1,
  LT: 2,
  GTE: 3,
  LTE: 4,
  NEQ: 5,
  IN: 6
} as const;

export type IntOperator = typeof IntOperator[keyof typeof IntOperator];

/**
 *
 * @export
 * @enum {string}
 */

export const KeyType = {
  STRING: 0,
  BOOLEAN: 1,
  FLOAT: 2,
  INT: 3,
  DATE_TIME: 4
} as const;

export type KeyType = typeof KeyType[keyof typeof KeyType];

/**
 *
 * @export
 * @enum {string}
 */

export const PlatformTypeType = {
  DEFAULT: 0,
  WEB: 1,
  IOS: 2,
  ANDROID: 3
} as const;

export type PlatformTypeType = typeof PlatformTypeType[keyof typeof PlatformTypeType];

/**
 *
 * @export
 * @enum {string}
 */

export const PrimitiveTypeType = {
  BOOL: 0,
  STRING: 1,
  INTEGER: 2,
  FLOAT: 3
} as const;

export type PrimitiveTypeType = typeof PrimitiveTypeType[keyof typeof PrimitiveTypeType];

/**
 *
 * @export
 * @enum {string}
 */

export const ProjectInviteStatus = {
  UNKNOWN: 0,
  PENDING: 1,
  ACCEPTED: 2,
  EXPIRED: 3
} as const;

export type ProjectInviteStatus = typeof ProjectInviteStatus[keyof typeof ProjectInviteStatus];

/**
 *
 * @export
 * @enum {string}
 */

export const ProjectMemberRole = {
  UNKNOWN: 0,
  ADMIN: 1,
  MEMBER: 2
} as const;

export type ProjectMemberRole = typeof ProjectMemberRole[keyof typeof ProjectMemberRole];

/**
 *
 * @export
 * @enum {string}
 */

export const StickinessType = {
  RANDOM: 0,
  KEYS: 1
} as const;

export type StickinessType = typeof StickinessType[keyof typeof StickinessType];

/**
 *
 * @export
 * @enum {string}
 */

export const StringOperator = {
  EQ: 0,
  CONTAINS: 1,
  IN: 2
} as const;

export type StringOperator = typeof StringOperator[keyof typeof StringOperator];
