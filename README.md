# CleanCodeFiberPostgreSQL

## Adım 1: Go Kurulumu

Golang'ı en az 1.19 sürümüne yükseltin. Aşağıdaki adımları takip edebilirsiniz:

- Önceki Golang kurulumunu temizleyin:

```
sudo rm -rf /usr/local/go
```
- Yeni Golang sürümünü indirin (örn. Go 1.19):

```
wget https://go.dev/dl/go1.19.12.linux-amd64.tar.gz
```
- İndirdiğiniz arşivi /usr/local klasörüne çıkarın:

```
sudo tar -C /usr/local -xzf go1.19.12.linux-amd64.tar.gz
```
- Ortam değişkenlerini güncelleyin:

```
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
```
- go version komutunu çalıştırarak doğru Go sürümünün kurulu olup olmadığını kontrol edin:

```
go version
```

## Adım 2: Postgresql Kurulumu:

- PostgreSQL veritabanını kurun:

```
sudo apt install -y postgresql postgresql-contrib
```

- PostgreSQL servisinin çalıştığından emin olun:

```
sudo systemctl start postgresql
sudo systemctl status postgresql
```

- Varsayılan PostgreSQL kullanıcısı olarak oturum açın:

```
sudo -i -u postgres
```

- Bir veritabanı oluşturun ve bir kullanıcıya yetki verin:

```
psql -c "CREATE DATABASE mydatabase;"
psql -c "CREATE USER myuser WITH ENCRYPTED PASSWORD 'mypassword';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE mydatabase TO myuser;"
```

- PostgreSQL sunucusunun güvenlik duvarı ayarlarında 5432 portunun açık olduğundan emin olun:

```
sudo ufw allow 5432/tcp
```

- PostgreSQL'in bağlantı izinlerini düzenlemek için pg_hba.conf dosyasını açın. Bu dosya genellikle /etc/postgresql/{sürüm numarası}/main/ dizinindedir.
Uzak IP adreslerinden bağlantıya izin vermek için bir satır ekleyin veya mevcut yapılandırmayı kontrol edin:

```
host    all             all             0.0.0.0/24        md5
```
- PostgreSQL'in harici bağlantılara izin vermesini sağlamak için yapılandırma dosyasını (postgresql.conf) açın ve listen_addresses parametresini '*' olarak ayarlayın:

```
listen_addresses = '*'
```

- Ardından PostgreSQL hizmetini yeniden başlatın:

```
sudo systemctl restart postgresql
```

- PostgreSQL veritabanı yöneticisi (postgres kullanıcısı) olarak oturum açın:

```
sudo -i -u postgres
```
- Sonrasında PostgreSQL kabuğuna girin:

```
psql
```

- mydatabase gibi veritabanınıza bağlanın:

```
\c mydatabase
```

- İzin vermek için aşağıdaki SQL komutu kullanın:

```
GRANT ALL PRIVILEGES ON TABLE items TO myuser;
```
- İzinlerin doğru verildiğinden emin olun:

```
\z items
```
Bu komut, tabloda hangi kullanıcının hangi tür izinlere sahip olduğunu gösterir.

- Sıralamaya (sequence) erişim izni vermek için aşağıdaki komutu çalıştırın. items_id_seq yerine kendi sıralama adınızı ekleyin:

```
GRANT ALL PRIVILEGES ON SEQUENCE items_id_seq TO myuser;
```
- Eğer yalnızca belirli izinler vermek istiyorsanız:
```
GRANT USAGE, SELECT ON SEQUENCE items_id_seq TO myuser;
```
- İlgili izinlerin doğru bir şekilde uygulandığını doğrulayın:

```
\dp items_id_seq
```

- Kabuğu kapatmak için:

```
\q
```
- Ardından postgres kullanıcısından çıkın:

```
exit
```

## Adım 3: Go Projesini Başlatın ve Fiber ile Gerekli Paketleri Kurun:

- Go modülünü başlatın:

```
mkdir simple-web-server
cd simple-web-server
go mod init simple-web-server
```
- Gerekli Go paketlerini ekleyin:

```
go get github.com/gofiber/fiber/v2
```
- PostgreSQL veritabanı için pgx veya gorm paketini ekleyin:

```
go get github.com/jackc/pgx/v5
```
## Adım 4: Fiber ve PostgreSQL ile Web Sunucusunu Oluşturun:

**main.go:**

```bash
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Veritabanı bağlantısını başlat
    initializeDatabaseConnection()
    defer pool.Close()

    // Fiber uygulamasını başlat
    app := fiber.New()

    // Rotaları tanımla
    app.Get("/", getRoot)
    app.Get("/items", getAllItems)
    app.Post("/items", addItem)

    // Sunucuyu başlat
    log.Fatal(app.Listen(":3000"))
}
```

**database.go:**

```bash
package main

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "log"
)

var pool *pgxpool.Pool

// Veritabanı bağlantısını kuran fonksiyon
func initializeDatabaseConnection() {
    dbUrl := "postgresql://myuser:mypassword@10.151.231.133:5432/mydatabase"
    var err error
    pool, err = pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatalf("Veritabanına bağlanılamadı: %v", err)
    }
}
```
**handlers.go:**

```bash
package main

import (
    "github.com/gofiber/fiber/v2"
)

// GET kök rotası
func getRoot(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
}

// Tüm öğeleri veritabanından alıp JSON formatında döndüren fonksiyon
func getAllItems(c *fiber.Ctx) error {
    rows, err := pool.Query(c.Context(), "SELECT id, name, price FROM items")
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }
    defer rows.Close()

    var items []Item
    for rows.Next() {
        var item Item
        rows.Scan(&item.ID, &item.Name, &item.Price)
        items = append(items, item)
    }

    return c.JSON(items)
}

// Yeni bir öğeyi ekleyen POST rotası
func addItem(c *fiber.Ctx) error {
    var newItem Item
    if err := c.BodyParser(&newItem); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    _, err := pool.Exec(c.Context(), "INSERT INTO items (name, price) VALUES ($1, $2)", newItem.Name, newItem.Price)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusCreated).JSON(newItem)
}
```
**models.go:**

```bash
package main

// Yapı (struct) tanımları
type Item struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

## Adım 5: PostgreSQL Tablosunu Oluşturun:

- Veritabanında bir items tablosu oluşturun. PostgreSQL oturumuna girip aşağıdaki komutları çalıştırabilirsiniz:

```
sudo -i -u postgres
psql
\c mydatabase
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    price NUMERIC(10, 2)
);
```

- Tabloyu kontrol edin:

```
\dt
```
Oluşturduğunuz items tablosunu görmelisiniz.

## Adım 6: Projeyi Çalıştırma:

- Son olarak, main.go dosyanızdaki kodu çalıştırarak web sunucusunu başlatın:

```
go run main.go
```
Artık uygulama, Fiber ve PostgreSQL ile basit bir web sunucusu olarak çalışacak.

## Adım 7: Postman'de Test Etme:

- GET isteği için URL:

```
http://localhost:3000/items
```

- POST isteği için URL:

```
http://localhost:3000/items
```
- Body kısmında JSON formatında veri ekleyin:
         
```
{
  "name": "new item",
  "price": 350.0
}
```
