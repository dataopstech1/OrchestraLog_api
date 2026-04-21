CREATE TABLE kafka_clusters (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id                  UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name                        VARCHAR(100) NOT NULL,
    version                     VARCHAR(20) NOT NULL,
    environment                 VARCHAR(20) NOT NULL CHECK (environment IN ('development', 'staging', 'production')),
    namespace                   VARCHAR(100) NOT NULL,
    status                      VARCHAR(20) NOT NULL CHECK (status IN (
        'pending', 'creating', 'running', 'failed', 'updating', 'deleting'
    )),
    broker_count                INT NOT NULL DEFAULT 3,
    broker_cpu                  VARCHAR(20) DEFAULT '2',
    broker_memory               VARCHAR(20) DEFAULT '4Gi',
    broker_storage              VARCHAR(20) DEFAULT '100Gi',
    broker_storage_class        VARCHAR(50) DEFAULT 'standard',
    broker_jvm_heap             VARCHAR(20) DEFAULT '2g',
    zk_enabled                  BOOLEAN DEFAULT true,
    zk_external_conn            TEXT,
    zk_replicas                 INT DEFAULT 3,
    zk_cpu                      VARCHAR(20) DEFAULT '0.5',
    zk_memory                   VARCHAR(20) DEFAULT '1Gi',
    zk_storage                  VARCHAR(20) DEFAULT '10Gi',
    auth_enabled                BOOLEAN DEFAULT false,
    auth_mechanism              VARCHAR(30) CHECK (auth_mechanism IN ('PLAIN', 'SCRAM-SHA-256', 'SCRAM-SHA-512')),
    authz_enabled               BOOLEAN DEFAULT false,
    authz_type                  VARCHAR(20) CHECK (authz_type IN ('simple', 'opa')),
    tls_enabled                 BOOLEAN DEFAULT false,
    tls_version                 VARCHAR(10) DEFAULT 'TLSv1.3',
    service_type                VARCHAR(20) DEFAULT 'ClusterIP' CHECK (service_type IN ('ClusterIP', 'NodePort', 'LoadBalancer')),
    ingress_enabled             BOOLEAN DEFAULT false,
    ingress_class               VARCHAR(50),
    ingress_annotations         JSONB,
    schema_registry_enabled     BOOLEAN DEFAULT false,
    schema_registry_replicas    INT DEFAULT 1,
    kafka_connect_enabled       BOOLEAN DEFAULT false,
    kafka_connect_replicas      INT DEFAULT 1,
    ksql_enabled                BOOLEAN DEFAULT false,
    ksql_replicas               INT DEFAULT 1,
    monitoring_enabled          BOOLEAN DEFAULT true,
    prometheus_enabled          BOOLEAN DEFAULT true,
    grafana_enabled             BOOLEAN DEFAULT false,
    jmx_exporter_enabled        BOOLEAN DEFAULT true,
    server_properties           JSONB DEFAULT '{}',
    jvm_options                 JSONB DEFAULT '[]',
    log_retention_hours         INT DEFAULT 168,
    default_replication_factor  INT DEFAULT 3,
    min_insync_replicas         INT DEFAULT 2,
    created_by                  UUID REFERENCES users(id),
    created_at                  TIMESTAMPTZ DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kafka_topics (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kafka_cluster_id    UUID NOT NULL REFERENCES kafka_clusters(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    partitions          INT NOT NULL DEFAULT 1,
    replication_factor  INT NOT NULL DEFAULT 1,
    retention_ms        BIGINT DEFAULT 604800000,
    status              VARCHAR(20) DEFAULT 'active',
    configuration       JSONB DEFAULT '{}',
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(kafka_cluster_id, name)
);

CREATE TABLE kafka_consumer_groups (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kafka_cluster_id    UUID NOT NULL REFERENCES kafka_clusters(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    state               VARCHAR(20),
    members             INT DEFAULT 0,
    topics_subscribed   JSONB DEFAULT '[]',
    total_lag           BIGINT DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kafka_resources (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kafka_cluster_id    UUID NOT NULL REFERENCES kafka_clusters(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    type                VARCHAR(30) NOT NULL CHECK (type IN ('topic', 'user', 'acl', 'connector', 'custom')),
    namespace           VARCHAR(100),
    configuration       JSONB DEFAULT '{}',
    status              VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'pending', 'failed', 'deleting')),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);
