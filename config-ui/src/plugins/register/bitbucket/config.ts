/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import type { PluginConfigType } from '@/plugins';

import {
  BaseConnectionConfig,
  ConnectionName,
  ConnectionEndpoint,
  ConnectionUsername,
  ConnectionPassword,
  ConnectionProxy,
  ConnectionRatelimit,
} from '../base';

import Icon from './assets/icon.svg';

export const BitBucketConfig: PluginConfigType = {
  ...BaseConnectionConfig,
  plugin: 'bitbucket',
  name: 'BitBucket',
  icon: Icon,
  isBeta: true,
  connection: {
    initialValues: {
      name: 'BitBucket',
      endpoint: 'https://api.bitbucket.org/2.0/',
      rateLimitPerHour: 10000,
    },
    fields: [
      ConnectionName({
        placeholder: 'eg. BitBucket',
      }),
      ConnectionEndpoint({
        placeholder: 'eg. https://api.bitbucket.org/2.0/',
      }),
      ConnectionUsername(),
      ConnectionPassword({
        label: 'App Password',
        placeholder: 'App Password',
      }),
      ConnectionProxy(),
      ConnectionRatelimit(),
    ],
  },
  entities: [],
  transformation: {},
};
