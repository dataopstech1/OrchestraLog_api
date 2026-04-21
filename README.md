# OrchestraLog API

OrchestraLog, Kafka, Spark, Flink, MLflow, JupyterHub gibi veri altyapısı servislerini tek bir platformdan yönetmek için geliştirilmiş bir **Data Infrastructure Orchestration** backend API'sidir.

---

## Teknoloji Stack

| Katman | Teknoloji |
|--------|-----------|
| Dil | Go 1.22+ |
| Router | go-chi/chi v5 |
| Veritabanı | PostgreSQL 16 |
| Cache | Redis 7 |
| Auth | JWT (HS256) + Refresh Token |
| ORM | sqlx + lib/pq |
| Container | Docker + Docker Compose |

---

## Gereksinimler

Başlamadan önce aşağıdakilerin kurulu olduğundan emin ol:

- [Go 1.22+](https://go.dev/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- Git

---

## Kurulum

### 1. Repoyu klonla

```bash
git clone https://github.com/dataopstech1/OrchestraLog_api.git
cd OrchestraLog_api
```

### 2. Ortam değişkenlerini ayarla

```bash
cp .env.example .env
```

`.env` dosyasını açıp gerekli değerleri doldur (geliştirme ortamı için varsayılanlar çalışır):

```env
PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=orchestralog
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_ACCESS_SECRET=change-this-in-production
JWT_REFRESH_SECRET=change-this-in-production
JWT_ACCESS_EXP_MINUTES=15
JWT_REFRESH_EXP_DAYS=7
```

> **Not:** `DB_PORT=5433` kullanılıyor çünkü Docker'daki PostgreSQL 5433 portundan yayın yapar (yerel 5432 çakışmasını önlemek için).

### 3. Docker ile PostgreSQL ve Redis'i başlat

```bash
docker compose up -d postgres redis
```

Containerların ayağa kalktığını doğrula:

```bash
docker compose ps
```

Çıktı şöyle görünmeli:

```
NAME                          STATUS
orchestralog-api-postgres-1   Up (healthy)
orchestralog-api-redis-1      Up (healthy)
```

### 4. Go bağımlılıklarını yükle

```bash
go mod tidy
```

### 5. Veritabanı migration'larını çalıştır

```bash
go run ./scripts/migrate.go up
```

Bu komut `migrations/` klasöründeki tüm `.up.sql` dosyalarını sırayla çalıştırır ve tüm tabloları oluşturur.

### 6. Örnek veri (seed) yükle

```bash
go run ./scripts/seed.go
```

Bu komut şu kayıtları oluşturur:

**Kullanıcılar:**
| E-posta | Şifre | Rol |
|---------|-------|-----|
| admin@orchestralog.com | Admin1234! | admin |
| operator@orchestralog.com | Oper1234! | operator |
| viewer@orchestralog.com | View1234! | viewer |

**Cluster'lar:**
| İsim | Region | Nodes |
|------|--------|-------|
| production-cluster-01 | eu-west-1 | 12 |
| staging-cluster-01 | eu-central-1 | 6 |
| dev-cluster-01 | us-east-1 | 3 |

### 7. API'yi başlat

```bash
go run ./cmd/server
```

Başarılı başlangıçta şu çıktıyı görmelisin:

```
2026/04/21 12:56:17 INFO starting OrchestraLog API port=8080 env=development
```

---

## API'yi Test Et

### Health Check

```bash
curl http://localhost:8080/api/v1/health
```

```json
{"status":"ok"}
```

### Giriş Yap

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@orchestralog.com","password":"Admin1234!"}'
```

Dönen `access_token` değerini kopyalayıp sonraki isteklerde kullan:

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

---

## Makefile Komutları

```bash
make run          # API'yi başlat
make build        # Binary derle (bin/server)
make docker-up    # Tüm Docker servislerini başlat
make docker-down  # Docker servislerini durdur
make migrate-up   # Migration'ları uygula
make migrate-down # Migration'ları geri al
make seed         # Örnek veri yükle
make tidy         # go mod tidy
make test         # Testleri çalıştır
make lint         # Linter çalıştır
```

---

## API Endpoint'leri

Tüm endpoint'ler `/api/v1` prefix'i ile başlar. Auth endpoint'leri hariç tüm isteklerde `Authorization: Bearer <token>` header'ı gereklidir.

### Auth
| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | /auth/login | Giriş yap, token al |
| POST | /auth/refresh | Access token yenile |
| POST | /auth/logout | Çıkış yap |
| GET | /auth/me | Oturum bilgilerini getir |

### Clusters
| Method | Endpoint | Rol |
|--------|----------|-----|
| GET | /clusters | Tümü |
| POST | /clusters | Admin/Operator |
| GET | /clusters/:id | Tümü |
| PUT | /clusters/:id | Admin/Operator |
| DELETE | /clusters/:id | Admin |

### Diğer Modüller

Aşağıdaki modüller benzer CRUD yapısına sahiptir:

| Prefix | Modül |
|--------|-------|
| /kafka/clusters | Apache Kafka |
| /spark/clusters | Apache Spark |
| /flink/clusters | Apache Flink |
| /hive/instances | Apache Hive |
| /hdfs/clusters | HDFS |
| /nifi/instances | Apache NiFi |
| /mlflow/instances | MLflow |
| /feast/instances | Feast Feature Store |
| /jupyterhub/instances | JupyterHub |
| /llm/deployments | LLM Deployments |
| /superset/instances | Apache Superset |
| /metabase/instances | Metabase |
| /n8n/instances | N8N Automation |
| /data-flows | Data Flow Pipeline |
| /users | Kullanıcı Yönetimi |
| /dashboard | Dashboard Özet |
| /monitoring/clusters/:id | Cluster Monitoring |

---

## Rol Yetkileri

| İşlem | Admin | Operator | Viewer |
|-------|-------|----------|--------|
| Okuma (GET) | ✅ | ✅ | ✅ |
| Oluşturma (POST) | ✅ | ✅ | ❌ |
| Güncelleme (PUT) | ✅ | ✅ | ❌ |
| Silme (DELETE) | ✅ | ❌ | ❌ |

---

## Proje Yapısı

```
orchestralog-api/
├── cmd/server/          # Uygulama giriş noktası
├── internal/
│   ├── config/          # Ortam değişkenleri
│   ├── handler/         # HTTP handler'ları
│   ├── middleware/       # Auth, CORS, rate limit, RBAC
│   ├── model/           # Veritabanı model struct'ları
│   ├── repository/      # Veritabanı sorguları
│   ├── server/          # Router ve server yapılandırması
│   └── service/         # İş mantığı katmanı
├── migrations/          # SQL migration dosyaları
├── pkg/
│   ├── apierror/        # Standart hata tipleri
│   ├── pagination/      # Sayfalama yardımcıları
│   └── response/        # JSON envelope yapısı
├── scripts/
│   ├── migrate.go       # Migration runner
│   └── seed.go          # Örnek veri yükleyici
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── Makefile
```

---

## Frontend

Bu API ile birlikte çalışan Next.js tabanlı frontend için:
👉 [OrchestraLog UI](https://github.com/dataopstech1/OrchestraLog_ui)

---

## Lisans

MIT
