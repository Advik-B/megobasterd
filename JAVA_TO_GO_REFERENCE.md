# Java to Go Migration Reference - MegaBasterd

This document maps Java classes, libraries, and patterns used in MegaBasterd to their Go equivalents.

---

## Library & Framework Mapping

| Java Library/Framework | Purpose in MegaBasterd | Go Equivalent | Notes |
|------------------------|------------------------|---------------|-------|
| **javax.swing** | GUI framework | `fyne.io/fyne/v2` or `github.com/wailsapp/wails/v2` | Fyne recommended for similarity to Swing |
| **javax.crypto** | Encryption/decryption | `crypto/aes`, `crypto/cipher`, `golang.org/x/crypto` | Go's crypto is excellent |
| **java.sql** / **sqlite-jdbc** | SQLite database | `database/sql` + `github.com/mattn/go-sqlite3` | Similar API |
| **jackson** (JSON) | JSON parsing | `encoding/json` or `github.com/json-iterator/go` | Built-in is good, json-iterator is faster |
| **commons-io** | File operations | `os`, `io`, `io/fs` | Go standard library |
| **java.util.concurrent** | Concurrency | `sync`, `sync.atomic`, `golang.org/x/sync/errgroup` | Go has better primitives |
| **java.net.HttpURLConnection** | HTTP client | `net/http` or `github.com/go-resty/resty/v2` | Standard lib is powerful |
| **java.util.zip.GZIPInputStream** | GZIP compression | `compress/gzip` | Standard library |
| **javax.xml.bind** | Base64 encoding | `encoding/base64` | Standard library |
| **xuggler** | Video processing | `github.com/3d0c/gmf` (FFmpeg bindings) | Alternative video library |

---

## Core Class Mappings

### Main Application Classes

| Java Class | Purpose | Go Package/Type | Implementation Notes |
|------------|---------|-----------------|---------------------|
| `MainPanel.java` | Application entry & main window | `cmd/megobasterd/main.go` + `internal/ui/mainwindow.go` | Split into entry point and UI |
| `Download.java` | Download manager | `internal/downloader/download.go` | Use goroutines instead of ExecutorService |
| `Upload.java` | Upload manager | `internal/uploader/upload.go` | Similar to Download structure |
| `MegaAPI.java` | MEGA API client | `internal/api/client.go` | Use resty for HTTP |
| `CryptTools.java` | Crypto utilities | `internal/crypto/crypto.go` | Go crypto package is well-designed |
| `DBTools.java` | Database operations | `internal/database/db.go` | Use database/sql interface |

### UI Classes (Swing → Fyne)

| Java Class | Swing Component | Fyne Equivalent | Example |
|------------|-----------------|-----------------|---------|
| `MainPanelView.java` | JFrame | `fyne.Window` | `app.NewWindow("Title")` |
| `DownloadView.java` | JPanel with table | `container` + `widget.Table` | `widget.NewTable(...)` |
| `SettingsDialog.java` | JDialog with form | `dialog.ShowForm()` | `dialog.ShowForm("Settings", ...)` |
| `AboutDialog.java` | JDialog | `dialog.ShowInformation()` | `dialog.ShowInformation("About", ...)` |
| Progress bars | JProgressBar | `widget.ProgressBar` | `widget.NewProgressBar()` |
| Tables | JTable | `widget.Table` | `widget.NewTable(...)` |
| Buttons | JButton | `widget.Button` | `widget.NewButton("Text", handler)` |
| Text fields | JTextField | `widget.Entry` | `widget.NewEntry()` |
| Labels | JLabel | `widget.Label` | `widget.NewLabel("Text")` |
| Tabs | JTabbedPane | `container.AppTabs` | `container.NewAppTabs(...)` |
| Menus | JMenu | `fyne.Menu` | `fyne.NewMenu(...)` |
| System Tray | SystemTray | `desktop.App.SetSystemTrayMenu()` | Fyne has built-in support |

---

## Language Feature Mappings

### Data Types

| Java | Go | Notes |
|------|-----|-------|
| `int` | `int` or `int32` | Go's `int` is platform-dependent (32 or 64-bit) |
| `long` | `int64` | |
| `float` | `float32` | |
| `double` | `float64` | |
| `String` | `string` | Go strings are UTF-8 by default |
| `byte[]` | `[]byte` | |
| `ArrayList<T>` | `[]T` | Go slices are built-in |
| `HashMap<K,V>` | `map[K]V` | Go maps are built-in |
| `Object` | `interface{}` or `any` | `any` is alias for `interface{}` (Go 1.18+) |

### Concurrency

| Java Pattern | Go Pattern | Example |
|--------------|------------|---------|
| `ExecutorService` | goroutines + channels | `go func() { ... }()` |
| `Future<T>` | channels or `errgroup` | `ch := make(chan Result)` |
| `CompletableFuture` | goroutines + channels | Multiple goroutines communicating via channels |
| `synchronized(obj)` | `sync.Mutex` | `mu.Lock(); defer mu.Unlock()` |
| `AtomicInteger` | `sync/atomic.Int32` | `atomic.AddInt32(&counter, 1)` |
| `CountDownLatch` | `sync.WaitGroup` | `wg.Add(n); wg.Wait()` |
| `Semaphore` | buffered channel | `sem := make(chan struct{}, n)` |
| `ThreadLocal` | goroutine-local storage | Use context or explicit passing |

### Error Handling

```java
// Java
try {
    result = riskyOperation();
} catch (IOException e) {
    log.error("Failed", e);
}
```

```go
// Go
result, err := riskyOperation()
if err != nil {
    log.Error("Failed", err)
}
```

### Null Safety

```java
// Java
if (obj != null) {
    obj.method();
}
```

```go
// Go - use pointer checking
if obj != nil {
    obj.Method()
}

// Or zero value checking for non-pointers
if str != "" {
    // use str
}
```

---

## Common Patterns

### Singleton Pattern

```java
// Java
public class SqliteSingleton {
    private static SqliteSingleton instance;
    
    public static synchronized SqliteSingleton getInstance() {
        if (instance == null) {
            instance = new SqliteSingleton();
        }
        return instance;
    }
}
```

```go
// Go - using sync.Once
package database

import "sync"

var (
    instance *SqliteDB
    once     sync.Once
)

func GetInstance() *SqliteDB {
    once.Do(func() {
        instance = &SqliteDB{}
        instance.Initialize()
    })
    return instance
}
```

### Observer Pattern (for UI updates)

```java
// Java - Swing uses observers extensively
public interface ClipboardChangeObserver {
    void clipboardChange(String newClipboard);
}

public class ClipboardSpy implements ClipboardChangeObservable {
    private List<ClipboardChangeObserver> observers = new ArrayList<>();
    
    public void addObserver(ClipboardChangeObserver observer) {
        observers.add(observer);
    }
    
    private void notifyObservers(String content) {
        for (ClipboardChangeObserver obs : observers) {
            obs.clipboardChange(content);
        }
    }
}
```

```go
// Go - using channels
type ClipboardSpy struct {
    subscribers []chan string
}

func (cs *ClipboardSpy) Subscribe() <-chan string {
    ch := make(chan string, 10)
    cs.subscribers = append(cs.subscribers, ch)
    return ch
}

func (cs *ClipboardSpy) notify(content string) {
    for _, ch := range cs.subscribers {
        select {
        case ch <- content:
        default:
            // Skip if buffer full
        }
    }
}

// Usage
spy := &ClipboardSpy{}
updates := spy.Subscribe()
go func() {
    for content := range updates {
        handleClipboardChange(content)
    }
}()
```

### Builder Pattern

```java
// Java
Download download = new Download.Builder()
    .url("https://mega.nz/...")
    .path("/downloads")
    .workers(6)
    .build();
```

```go
// Go - struct literals often sufficient
download := &Download{
    URL:     "https://mega.nz/...",
    Path:    "/downloads",
    Workers: 6,
}

// Or functional options pattern for complex cases
download := NewDownload(
    WithURL("https://mega.nz/..."),
    WithPath("/downloads"),
    WithWorkers(6),
)
```

---

## Crypto Operations Mapping

### AES Encryption

```java
// Java (from CryptTools.java)
public static Cipher genDecrypter(String algo, String mode, byte[] key, byte[] iv) 
    throws NoSuchAlgorithmException, NoSuchPaddingException, 
           InvalidKeyException, InvalidAlgorithmParameterException {
    SecretKeySpec skeySpec = new SecretKeySpec(key, algo);
    Cipher decryptor = Cipher.getInstance(mode);
    
    if (iv != null) {
        IvParameterSpec ivParameterSpec = new IvParameterSpec(iv);
        decryptor.init(Cipher.DECRYPT_MODE, skeySpec, ivParameterSpec);
    } else {
        decryptor.init(Cipher.DECRYPT_MODE, skeySpec);
    }
    return decryptor;
}
```

```go
// Go
func NewDecrypter(key, iv []byte) (cipher.BlockMode, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    if iv != nil {
        return cipher.NewCBCDecrypter(block, iv), nil
    }
    
    return cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize)), nil
}

// Usage
decrypter, err := NewDecrypter(key, iv)
if err != nil {
    return err
}
decrypter.CryptBlocks(dst, src)
```

### PBKDF2 Key Derivation

```java
// Java
public static byte[] pbkdf2(String password, byte[] salt, int iterations, int keyLength) {
    KeySpec spec = new PBEKeySpec(password.toCharArray(), salt, iterations, keyLength);
    SecretKeyFactory factory = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256");
    return factory.generateSecret(spec).getEncoded();
}
```

```go
// Go
import "golang.org/x/crypto/pbkdf2"

func DeriveKey(password string, salt []byte, iterations, keyLen int) []byte {
    return pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)
}
```

---

## File I/O Mapping

### Reading Files

```java
// Java
byte[] data = Files.readAllBytes(Paths.get("/path/to/file"));
```

```go
// Go
data, err := os.ReadFile("/path/to/file")
if err != nil {
    log.Fatal(err)
}
```

### Writing Files

```java
// Java
Files.write(Paths.get("/path/to/file"), data, StandardOpenOption.CREATE);
```

```go
// Go
err := os.WriteFile("/path/to/file", data, 0644)
if err != nil {
    log.Fatal(err)
}
```

### Copying Files

```java
// Java
Files.copy(source, target, StandardCopyOption.REPLACE_EXISTING);
```

```go
// Go
srcFile, _ := os.Open(sourcePath)
defer srcFile.Close()

dstFile, _ := os.Create(targetPath)
defer dstFile.Close()

io.Copy(dstFile, srcFile)
```

---

## Network Operations

### HTTP GET Request

```java
// Java
URL url = new URL("https://api.mega.co.nz/...");
HttpURLConnection conn = (HttpURLConnection) url.openConnection();
conn.setRequestMethod("GET");
conn.setRequestProperty("User-Agent", "MegaBasterd");

InputStream in = conn.getInputStream();
// Read response...
```

```go
// Go (using standard library)
req, _ := http.NewRequest("GET", "https://api.mega.co.nz/...", nil)
req.Header.Set("User-Agent", "MegaBasterd")

client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

body, _ := io.ReadAll(resp.Body)

// Or using Resty (recommended)
import "github.com/go-resty/resty/v2"

client := resty.New()
resp, err := client.R().
    SetHeader("User-Agent", "MegaBasterd").
    Get("https://api.mega.co.nz/...")
```

### HTTP POST with JSON

```java
// Java
HttpURLConnection conn = (HttpURLConnection) url.openConnection();
conn.setRequestMethod("POST");
conn.setRequestProperty("Content-Type", "application/json");
conn.setDoOutput(true);

String jsonBody = objectMapper.writeValueAsString(data);
try (OutputStream os = conn.getOutputStream()) {
    os.write(jsonBody.getBytes());
}
```

```go
// Go with Resty
type RequestData struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

data := RequestData{Field1: "value", Field2: 42}

resp, err := client.R().
    SetHeader("Content-Type", "application/json").
    SetBody(data).
    Post("https://api.mega.co.nz/...")
```

---

## Database Operations

### SQLite Connection

```java
// Java
Class.forName("org.sqlite.JDBC");
Connection conn = DriverManager.getConnection("jdbc:sqlite:megabasterd.db");
```

```go
// Go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

db, err := sql.Open("sqlite3", "megabasterd.db")
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### Prepared Statements

```java
// Java
PreparedStatement stmt = conn.prepareStatement(
    "INSERT INTO downloads (url, path, size) VALUES (?, ?, ?)"
);
stmt.setString(1, url);
stmt.setString(2, path);
stmt.setLong(3, size);
stmt.executeUpdate();
```

```go
// Go
stmt, err := db.Prepare("INSERT INTO downloads (url, path, size) VALUES (?, ?, ?)")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

_, err = stmt.Exec(url, path, size)
if err != nil {
    log.Fatal(err)
}
```

### Query Results

```java
// Java
ResultSet rs = stmt.executeQuery("SELECT * FROM downloads");
while (rs.next()) {
    String url = rs.getString("url");
    long size = rs.getLong("size");
    // Process...
}
```

```go
// Go
rows, err := db.Query("SELECT * FROM downloads")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

for rows.Next() {
    var url string
    var size int64
    err := rows.Scan(&url, &size)
    if err != nil {
        log.Fatal(err)
    }
    // Process...
}
```

---

## Specific MegaBasterd Components

### ChunkDownloader Pattern

```java
// Java - uses ExecutorService for parallel downloads
ExecutorService executor = Executors.newFixedThreadPool(workers);
for (Chunk chunk : chunks) {
    executor.submit(() -> downloadChunk(chunk));
}
executor.shutdown();
executor.awaitTermination(Long.MAX_VALUE, TimeUnit.SECONDS);
```

```go
// Go - using errgroup for coordinated goroutines
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())

// Limit concurrency
sem := make(chan struct{}, workers)

for _, chunk := range chunks {
    chunk := chunk // Capture loop variable
    
    g.Go(func() error {
        sem <- struct{}{}        // Acquire
        defer func() { <-sem }() // Release
        
        return downloadChunk(ctx, chunk)
    })
}

if err := g.Wait(); err != nil {
    log.Fatal(err)
}
```

### Progress Tracking

```java
// Java - typically uses SwingWorker or similar
SwingUtilities.invokeLater(() -> {
    progressBar.setValue(progress);
    speedLabel.setText(speed + " MB/s");
});
```

```go
// Go - using channels for progress updates
type Progress struct {
    BytesDownloaded int64
    TotalBytes      int64
    Speed           float64
}

progressChan := make(chan Progress, 10)

// Producer (download goroutine)
go func() {
    for {
        // ... download work ...
        progressChan <- Progress{
            BytesDownloaded: downloaded,
            TotalBytes:      total,
            Speed:           currentSpeed,
        }
    }
}()

// Consumer (UI goroutine)
go func() {
    for p := range progressChan {
        // Update UI
        progressBar.SetValue(float64(p.BytesDownloaded) / float64(p.TotalBytes))
        speedLabel.SetText(fmt.Sprintf("%.2f MB/s", p.Speed))
    }
}()
```

---

## Configuration Management

### Java Properties

```java
// Java
Properties props = new Properties();
props.load(new FileInputStream("config.properties"));
String downloadPath = props.getProperty("download.path", "/Downloads");
```

### Go Viper (Recommended)

```go
// Go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.SetDefault("download.path", "/Downloads")

if err := viper.ReadInConfig(); err != nil {
    log.Fatal(err)
}

downloadPath := viper.GetString("download.path")
```

---

## Testing Comparison

### JUnit → Go Testing

```java
// Java (JUnit)
@Test
public void testEncryption() {
    byte[] plaintext = "hello".getBytes();
    byte[] key = generateKey();
    
    byte[] encrypted = CryptTools.encrypt(plaintext, key);
    byte[] decrypted = CryptTools.decrypt(encrypted, key);
    
    assertArrayEquals(plaintext, decrypted);
}
```

```go
// Go (testing package)
func TestEncryption(t *testing.T) {
    plaintext := []byte("hello")
    key := generateKey()
    
    encrypted, err := Encrypt(plaintext, key)
    if err != nil {
        t.Fatal(err)
    }
    
    decrypted, err := Decrypt(encrypted, key)
    if err != nil {
        t.Fatal(err)
    }
    
    if !bytes.Equal(plaintext, decrypted) {
        t.Errorf("got %v, want %v", decrypted, plaintext)
    }
}

// Or with testify
import "github.com/stretchr/testify/assert"

func TestEncryption(t *testing.T) {
    plaintext := []byte("hello")
    key := generateKey()
    
    encrypted, err := Encrypt(plaintext, key)
    assert.NoError(t, err)
    
    decrypted, err := Decrypt(encrypted, key)
    assert.NoError(t, err)
    assert.Equal(t, plaintext, decrypted)
}
```

---

## Key Differences to Remember

1. **Error Handling**: Go uses explicit error returns, not exceptions
2. **Nil vs Null**: Go uses `nil` for pointers, zero values for value types
3. **Interfaces**: Go interfaces are implicit, not explicit
4. **Goroutines**: Much lighter than Java threads
5. **Channels**: Built-in, type-safe communication between goroutines
6. **Defer**: Go's defer executes at function return (like try-finally)
7. **Pointers**: Go has explicit pointers with & and * operators
8. **No Classes**: Go uses structs with methods, not classes
9. **Composition**: Go uses embedding, not inheritance
10. **Package System**: Go uses packages, not class-based namespaces

---

## Common Pitfalls

### 1. Loop Variable Capture

```go
// WRONG
for _, item := range items {
    go func() {
        process(item) // Bug: all goroutines see last item
    }()
}

// CORRECT
for _, item := range items {
    item := item // Capture loop variable
    go func() {
        process(item)
    }()
}
```

### 2. Nil Pointer Dereference

```go
// Check before use
if obj != nil {
    obj.Method()
}
```

### 3. Closing Channels

```go
// Only close channels from sender side
ch := make(chan int)
go func() {
    for i := 0; i < 10; i++ {
        ch <- i
    }
    close(ch) // Close when done sending
}()

for v := range ch { // Range automatically handles closed channel
    fmt.Println(v)
}
```

---

## Resources for Java Developers

- [Go for Java Programmers](https://go.dev/talks/2014/go4java.slide)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go Proverbs](https://go-proverbs.github.io/)

---

This reference should help you translate MegaBasterd's Java code to Go efficiently. When in doubt, consult the Go standard library documentation – it's excellent!
