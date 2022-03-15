import { edgeConfig } from '@ory/integrations/next';
import { Configuration, V0alpha2Api } from '@ory/kratos-client';

// Initialize the Ory Kratos SDK which will connect to the
edgeConfig.basePath = '/identity';

export default new V0alpha2Api(new Configuration(edgeConfig));
