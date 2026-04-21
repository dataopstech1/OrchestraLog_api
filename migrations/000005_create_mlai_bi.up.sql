-- MLFlow
CREATE TABLE mlflow_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    tracking_url    TEXT,
    artifact_store  TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE mlflow_experiments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mlflow_instance_id  UUID NOT NULL REFERENCES mlflow_instances(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    status              VARCHAR(20) DEFAULT 'active',
    runs_total          INT DEFAULT 0,
    best_metric_name    VARCHAR(100),
    best_metric_value   FLOAT,
    created_by          VARCHAR(100),
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE mlflow_models (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mlflow_instance_id  UUID NOT NULL REFERENCES mlflow_instances(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    latest_version      INT DEFAULT 1,
    stage               VARCHAR(30),
    description         TEXT,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- Feast
CREATE TABLE feast_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    online_store    VARCHAR(50),
    offline_store   VARCHAR(50),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE feast_entities (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feast_instance_id   UUID NOT NULL REFERENCES feast_instances(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    value_type          VARCHAR(30) NOT NULL,
    description         TEXT,
    join_keys           JSONB DEFAULT '[]',
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE feast_feature_views (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feast_instance_id   UUID NOT NULL REFERENCES feast_instances(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    entities            JSONB DEFAULT '[]',
    features            JSONB DEFAULT '[]',
    ttl                 VARCHAR(50),
    source              VARCHAR(100),
    online              BOOLEAN DEFAULT true,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- JupyterHub
CREATE TABLE jupyterhub_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    hub_url         TEXT,
    max_users       INT DEFAULT 100,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE jupyterhub_notebooks (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jupyterhub_instance_id  UUID NOT NULL REFERENCES jupyterhub_instances(id) ON DELETE CASCADE,
    name                    VARCHAR(200) NOT NULL,
    owner                   VARCHAR(100) NOT NULL,
    status                  VARCHAR(20) NOT NULL,
    image                   VARCHAR(200),
    cpu_limit               VARCHAR(20),
    memory_limit            VARCHAR(20),
    gpu_limit               INT DEFAULT 0,
    last_activity           TIMESTAMPTZ,
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW()
);

-- LLM Deployments
CREATE TABLE llm_deployments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id          UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name                VARCHAR(100) NOT NULL,
    model_name          VARCHAR(100) NOT NULL,
    model_version       VARCHAR(50),
    namespace           VARCHAR(100) NOT NULL,
    status              VARCHAR(20) NOT NULL,
    endpoint_url        TEXT,
    replicas            INT DEFAULT 1,
    gpu_count           INT DEFAULT 0,
    gpu_type            VARCHAR(50),
    max_tokens          INT,
    context_window      INT,
    requests_per_min    INT DEFAULT 0,
    avg_latency_ms      FLOAT DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- Superset
CREATE TABLE superset_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    url             TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE superset_dashboards (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    superset_instance_id    UUID NOT NULL REFERENCES superset_instances(id) ON DELETE CASCADE,
    name                    VARCHAR(200) NOT NULL,
    slug                    VARCHAR(200),
    status                  VARCHAR(20) DEFAULT 'published',
    charts_count            INT DEFAULT 0,
    owners                  JSONB DEFAULT '[]',
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW()
);

-- Metabase
CREATE TABLE metabase_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    url             TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE metabase_dashboards (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metabase_instance_id    UUID NOT NULL REFERENCES metabase_instances(id) ON DELETE CASCADE,
    name                    VARCHAR(200) NOT NULL,
    collection              VARCHAR(200),
    cards_count             INT DEFAULT 0,
    views_count             INT DEFAULT 0,
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW()
);

-- N8N
CREATE TABLE n8n_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    url             TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE n8n_workflows (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    n8n_instance_id     UUID NOT NULL REFERENCES n8n_instances(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    status              VARCHAR(20) NOT NULL,
    nodes_count         INT DEFAULT 0,
    connections_count   INT DEFAULT 0,
    last_execution_at   TIMESTAMPTZ,
    total_executions    INT DEFAULT 0,
    success_rate        FLOAT DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);
