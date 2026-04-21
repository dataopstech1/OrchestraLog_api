CREATE TABLE clusters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL CHECK (status IN ('healthy', 'warning', 'critical', 'offline')),
    region          VARCHAR(50) NOT NULL,
    k8s_version     VARCHAR(20) NOT NULL,
    api_server_url  TEXT,
    kubeconfig      TEXT,
    nodes_total     INT DEFAULT 0,
    nodes_ready     INT DEFAULT 0,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE namespaces (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL CHECK (status IN ('Active', 'Terminating')),
    cpu_quota       VARCHAR(20),
    memory_quota    VARCHAR(20),
    pods_quota      VARCHAR(20),
    cpu_used        VARCHAR(20),
    memory_used     VARCHAR(20),
    pods_used       INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(cluster_id, name)
);

CREATE TABLE resources (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    namespace_id    UUID NOT NULL REFERENCES namespaces(id) ON DELETE CASCADE,
    name            VARCHAR(200) NOT NULL,
    type            VARCHAR(30) NOT NULL CHECK (type IN (
        'Deployment', 'StatefulSet', 'Service', 'ConfigMap',
        'Secret', 'PVC', 'Job', 'CronJob', 'Pod'
    )),
    status          VARCHAR(20) NOT NULL CHECK (status IN (
        'Running', 'Pending', 'Failed', 'Succeeded', 'Unknown'
    )),
    replicas_ready  INT,
    replicas_total  INT,
    cpu_usage       VARCHAR(20),
    memory_usage    VARCHAR(20),
    age             VARCHAR(50),
    labels          JSONB,
    annotations     JSONB,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE cluster_metrics (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    metric_type     VARCHAR(30) NOT NULL,
    value           FLOAT NOT NULL,
    unit            VARCHAR(20),
    node_name       VARCHAR(100),
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cluster_metrics_time ON cluster_metrics (cluster_id, metric_type, recorded_at DESC);
CREATE INDEX idx_resources_namespace ON resources (namespace_id, type);
