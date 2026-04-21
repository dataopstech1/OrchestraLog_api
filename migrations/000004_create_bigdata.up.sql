-- Spark
CREATE TABLE spark_clusters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    master_url      TEXT,
    worker_count    INT DEFAULT 4,
    worker_cpu      VARCHAR(20) DEFAULT '4',
    worker_memory   VARCHAR(20) DEFAULT '8Gi',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE spark_applications (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    spark_cluster_id    UUID NOT NULL REFERENCES spark_clusters(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    status              VARCHAR(20) NOT NULL,
    duration            VARCHAR(50),
    started_at          TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    stages_total        INT DEFAULT 0,
    stages_completed    INT DEFAULT 0,
    executor_count      INT DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

-- Flink
CREATE TABLE flink_clusters (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id          UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name                VARCHAR(100) NOT NULL,
    version             VARCHAR(20) NOT NULL,
    namespace           VARCHAR(100) NOT NULL,
    status              VARCHAR(20) NOT NULL,
    jobmanager_url      TEXT,
    taskmanager_count   INT DEFAULT 4,
    slots_per_tm        INT DEFAULT 4,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE flink_jobs (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flink_cluster_id        UUID NOT NULL REFERENCES flink_clusters(id) ON DELETE CASCADE,
    name                    VARCHAR(200) NOT NULL,
    status                  VARCHAR(20) NOT NULL,
    start_time              TIMESTAMPTZ,
    duration                VARCHAR(50),
    parallelism             INT DEFAULT 1,
    checkpoints_completed   INT DEFAULT 0,
    bytes_in                BIGINT DEFAULT 0,
    bytes_out               BIGINT DEFAULT 0,
    records_in              BIGINT DEFAULT 0,
    records_out             BIGINT DEFAULT 0,
    created_at              TIMESTAMPTZ DEFAULT NOW()
);

-- Hive
CREATE TABLE hive_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    metastore_url   TEXT,
    hiveserver2_url TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE hive_tables (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hive_instance_id    UUID NOT NULL REFERENCES hive_instances(id) ON DELETE CASCADE,
    database_name       VARCHAR(100) NOT NULL,
    table_name          VARCHAR(200) NOT NULL,
    table_type          VARCHAR(30),
    format              VARCHAR(30),
    partitions          INT DEFAULT 0,
    total_size          VARCHAR(50),
    rows_count          BIGINT DEFAULT 0,
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- HDFS
CREATE TABLE hdfs_clusters (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id                  UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name                        VARCHAR(100) NOT NULL,
    version                     VARCHAR(20) NOT NULL,
    namespace                   VARCHAR(100) NOT NULL,
    status                      VARCHAR(20) NOT NULL,
    namenode_count              INT DEFAULT 2,
    datanode_count              INT DEFAULT 6,
    journalnode_count           INT DEFAULT 3,
    total_capacity              VARCHAR(50),
    used_capacity               VARCHAR(50),
    remaining_capacity          VARCHAR(50),
    total_blocks                INT DEFAULT 0,
    under_replicated_blocks     INT DEFAULT 0,
    created_at                  TIMESTAMPTZ DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ DEFAULT NOW()
);

-- NiFi
CREATE TABLE nifi_instances (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id      UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    version         VARCHAR(20) NOT NULL,
    namespace       VARCHAR(100) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    ui_url          TEXT,
    node_count      INT DEFAULT 3,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE nifi_process_groups (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nifi_instance_id        UUID NOT NULL REFERENCES nifi_instances(id) ON DELETE CASCADE,
    name                    VARCHAR(200) NOT NULL,
    status                  VARCHAR(20) NOT NULL,
    processors_total        INT DEFAULT 0,
    processors_running      INT DEFAULT 0,
    processors_stopped      INT DEFAULT 0,
    input_bytes_per_sec     VARCHAR(30),
    output_bytes_per_sec    VARCHAR(30),
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW()
);
