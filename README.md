# Crypto Orderbook

Gerçek zamanlı kripto para sipariş defteri uygulaması. WebSocket ile canlı güncellemeler, kullanıcı girişi ve PostgreSQL ile kalıcı veri saklama.

## Ne yapıyor?

Basit bir crypto exchange orderbook simülasyonu. Kullanıcılar kayıt olup giriş yapabiliyor, alım/satım siparişi verebiliyor. Siparişler gerçek zamanlı olarak tüm kullanıcılara gösteriliyor.

## Teknolojiler

- Backend: Go + Fiber framework
- Frontend: React + TypeScript  
- Database: PostgreSQL
- Real-time: WebSocket
- Deployment: Docker Compose

## Kurulum

Docker yüklü olması yeterli:
```bash
git clone <repo-url>
cd crypto-orderbook
docker-compose up --build
```

3-4 dakika bekle, sonra `http://localhost:3000` adresinden açabilirsin.

## Kullanım

1. Register sayfasından hesap oluştur
2. Login yap
3. Sağdaki formdan sipariş ver (buy ya da sell)
4. Sol tarafta siparişlerin görünüyor
5. Başka bir tarayıcıdan farklı kullanıcıyla gir, sipariş ver - ilk kullanıcı da anlık görüyor

## Ayarlar

`docker-compose.yml` dosyasında environment variable'lar var:

- `PORT`: Backend portu (8080)
- `DB_PASSWORD`: Veritabanı şifresi (varsayılan: 123456)
- `JWT_SECRET`: Token için secret key (production'da mutlaka değiştir)

## Database

PostgreSQL kullanıyor. İki tane tablo var:

**users**: Kullanıcı bilgileri (email, username, şifre hash'i)  
**orders**: Sipariş bilgileri (user_id, type, price, amount, status)

Migration'lar backend başlarken otomatik çalışıyor, bir şey yapman gerekmiyor.

## Proje yapısı
```
backend/
  cmd/server/main.go          # ana dosya
  internal/
    handlers/                 # API endpoint'ler
    repository/               # database işlemleri
    websocket/                # websocket hub
    middleware/               # JWT kontrolü
    models/                   # struct'lar
    
frontend/
  src/
    components/               # React component'ler
    hooks/                    # custom hook'lar
    services/                 # API çağrıları
```

## API

**Auth:**
- `POST /api/auth/register` - Kayıt ol
- `POST /api/auth/login` - Giriş yap

**Orders:** (token gerekli)
- `GET /api/orders` - Tüm siparişleri getir
- `POST /api/orders` - Yeni sipariş oluştur
- `GET /api/orders/my` - Kendi siparişlerimi getir

**WebSocket:**
- `WS /ws` - Canlı güncellemeler için

## Local development

Docker kullanmadan da çalıştırabilirsin:

Backend:
```bash
cd backend
go run cmd/server/main.go
```

Frontend:
```bash
cd frontend
npm install
npm run dev
```

PostgreSQL'in çalışıyor olması lazım tabii.

## Mimari

Repository pattern kullandım backend'de. Handler'lar direkt database'e yazmıyor, repository üzerinden gidiyor. WebSocket için hub pattern var - tüm client'ları merkezi bir yerde yönetiyor.

Frontend'de React Context ile auth state'i tutuyorum. WebSocket bağlantısı custom hook ile yönetiliyor.

Real-time kısım şöyle çalışıyor: Sipariş oluşturulunca backend hem database'e yazıyor hem de websocket hub'a gönderiyor. Hub da bağlı tüm client'lara broadcast ediyor.

## Güvenlik

- Şifreler bcrypt ile hash'leniyor
- JWT token kullanıyorum (header'da Bearer token)
- SQL injection'a karşı prepared statement'lar var
- CORS sadece belirli origin'lere açık

## Notlar

Production'da JWT_SECRET ve DB_PASSWORD'ü mutlaka değiştir. Şu an development için basit değerler kullanıyorum.

Frontend Nginx ile serve ediliyor container'da. Backend port'u 8080, frontend 3000'de.

## Lisans

MIT
