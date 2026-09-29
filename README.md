# Salon Management System

![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-16-black?logo=next.js&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)
![JWT](https://img.shields.io/badge/Auth-JWT-000000?logo=jsonwebtokens&logoColor=white)
![Cloudflare](https://img.shields.io/badge/Storage-Cloudflare_R2-F38020?logo=cloudflare&logoColor=white)

## Mô tả ngắn

**Salon Management System** là nền tảng dành cho salon tóc và spa, hợp nhất lịch hẹn, đơn hàng, tồn kho sản phẩm, vật tư dịch vụ và vai trò nhân viên trong một hệ thống duy nhất. Dự án thay thế các file bảng tính rời rạc bằng mô hình dữ liệu theo tổ chức, API Go theo Clean Architecture và dashboard Next.js hiện đại. Hệ thống vận hành toàn bộ quy trình hằng ngày của salon — từ đặt lịch đến phục vụ, bán hàng, kiểm soát kho và thanh toán.

## Tính năng chính

- **Xử lý đơn hàng robust** — đơn hàng với tổng tiền do server tự tính, tự động trừ kho cho sản phẩm bán lẻ và trừ định mức vật tư cho dịch vụ; vòng đời đầy đủ (`serve`, `complete`, `cancel`, `refund`) kèm thanh toán theo từng đơn.
- **Quản lý danh mục theo tổ chức** — danh mục thống nhất (`service` / `product` / `material`), dịch vụ kèm định mức vật tư, giá sản phẩm theo từng tổ chức.
- **Phân quyền truy cập theo vai trò (RBAC)** — xác thực JWT (access token ngắn hạn + refresh cookie httpOnly) với phân quyền theo phòng ban và vai trò nhân viên.
- **Tải tệp lớn** — cơ chế upload chia nhỏ (chunked) cho video/tệp lớn qua object storage tương thích S3 (Cloudflare R2 ) giúp giảm tải server.
- **Quản lý lịch hẹn** — nhận diện khách hàng qua số điện thoại, đặt lịch an toàn race-condition, tự sinh mã lịch hẹn, đổi/hủy lịch linh hoạt.
- **Kiểm soát tồn kho** — phiếu nhập/xuất kho, sổ cái giao dịch và cảnh báo sắp hết hàng.

## Công nghệ sử dụng

- **Backend:** Golang (Gin), Clean Architecture, xác thực JWT
- **Frontend:** Next.js, TypeScript, TanStack React Query, Redux Toolkit, Axios
- **Database:** PostgreSQL 16 (driver pgx, migration có version kèm rollback + seed data)
- **DevOps & Hạ tầng:** Docker, Docker Compose, Cloudflare R2 / Azure Storage, tài liệu API Swagger (swaggo)

## Kiến trúc hệ thống

Backend tuân thủ nghiêm ngặt **Clean Architecture** với các nguyên tắc **Domain-Driven Design (DDD)**:

```
Delivery (Gin handlers / routers / middleware)
   → Usecase (luật nghiệp vụ ứng dụng, map lỗi tập trung một điểm)
      → Domain (entities, DTOs, repository interfaces — Go thuần, không framework)
         → Infrastructure (PostgreSQL repositories qua pgx)
```

- **Tầng Domain** chứa business logic thuần túy: entities, validation và các hợp đồng repository.
- **Tầng Usecase** điều phối luồng nghiệp vụ (ví dụ: tạo đơn kiểm tra mọi tham chiếu thuộc tổ chức, tự tính tổng tiền, gói các ghi đa bảng trong transaction ngắn).
- **Tầng Infrastructure** triển khai persistence trên PostgreSQL với truy vấn theo phạm vi tổ chức.
- **Tầng Delivery** expose HTTP API qua Gin dưới `/api/v1`, JWT middleware nạp organization context.

## Bắt đầu (Cài đặt local)

### Yêu cầu

- Docker + Docker Compose
- Go ≥ 1.27
- Node.js ≥ 20

### Chạy dự án

```bash
# 1. Clone repository
git clone https://github.com/daidat02/salon-management.git
cd salon-management-system

# 2. Cấu hình môi trường
cp server/.env.example server/.env
cp client/.env.example client/.env
# Chỉnh sửa 2 file trên (thông tin database, JWT secret, API base URL, storage keys)

# 3. Khởi chạy database
cd server
docker compose up -d              # Postgres :5432, Adminer :8081

# 4. Migration (dùng DB_URL trong server/.env)
make migrate-up                   # migrate up toàn bộ
# make migrate-down               # rollback 1 bản khi cần

# 5. Chạy backend
make tidy                         # go mod tidy
make run                          # API → http://localhost:8080 (Swagger: /swagger/index.html)
make dev                        # chạy bằng air (hot-reload)

# 6. Chạy frontend (terminal mới)
cd ../client
npm install
npm run dev                       # Web → http://localhost:3000
```

Lệnh hữu ích khác (trong `server/`): `make test` (test usecase), `make lint` (`go vet`), `make migrate-create name=ten_bang` (tạo migration mới).

## Liên hệ / Tác giả

- **Tên:** Đại Đạt
- **Email:** [daidat01@gmail.com]
- **GitHub:** [[github.com/daidat02](https://github.com/daidat02/)]
