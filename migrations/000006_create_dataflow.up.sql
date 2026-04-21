CREATE TABLE data_flows (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    status          VARCHAR(20) NOT NULL CHECK (status IN (
        'draft', 'deployed', 'running', 'stopped', 'failed'
    )),
    template_id     VARCHAR(50),
    layout          VARCHAR(10) DEFAULT 'auto' CHECK (layout IN ('auto', 'manual')),
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE flow_nodes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_flow_id    UUID NOT NULL REFERENCES data_flows(id) ON DELETE CASCADE,
    node_id         VARCHAR(100) NOT NULL,
    label           VARCHAR(200) NOT NULL,
    resource_type   VARCHAR(30) NOT NULL CHECK (resource_type IN (
        'kafkaTopic', 'kafkaBroker', 'kafkaConnect',
        'sparkJob', 'flinkJob',
        'hdfsStorage', 'elasticsearch',
        'database', 'apiSource', 'fileStorage',
        'monitoring', 'dashboard', 'mlModel',
        'dataWarehouse', 'streamProcessor'
    )),
    cluster         VARCHAR(100),
    namespace       VARCHAR(100),
    status          VARCHAR(20) DEFAULT 'healthy',
    position_x      FLOAT NOT NULL,
    position_y      FLOAT NOT NULL,
    config          JSONB DEFAULT '{}',
    metrics         JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(data_flow_id, node_id)
);

CREATE TABLE flow_edges (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_flow_id    UUID NOT NULL REFERENCES data_flows(id) ON DELETE CASCADE,
    edge_id         VARCHAR(100) NOT NULL,
    source_node_id  VARCHAR(100) NOT NULL,
    target_node_id  VARCHAR(100) NOT NULL,
    flow_type       VARCHAR(20) NOT NULL CHECK (flow_type IN (
        'dataStream', 'batchTransfer', 'monitoring', 'errorFlow', 'controlFlow'
    )),
    label           VARCHAR(200),
    animated        BOOLEAN DEFAULT false,
    config          JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(data_flow_id, edge_id)
);

CREATE TABLE flow_templates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    category        VARCHAR(50),
    nodes           JSONB NOT NULL,
    edges           JSONB NOT NULL,
    thumbnail_url   TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
