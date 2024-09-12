package main

import (
    "bufio"
    "fmt"
    "log"
    "net/http"
    "os"
    "io/ioutil"
    "strings"
)

func main() {
    file, err := os.Open("message.txt")
    if err != nil {
        log.Fatalf("Dosya bulunamadı %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        token := scanner.Text()
        if token == "" {
            continue
        }

        req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
        if err != nil {
            log.Fatalf("istek gonderilemedi %v", err)
        }
        req.Header.Set("Authorization", "Bearer "+token)

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            log.Fatalf("istek gonderilemedi %v", err)
        }
        defer resp.Body.Close()

        body, _ := ioutil.ReadAll(resp.Body)
        responseText := string(body)

        switch resp.StatusCode {
        case 200:
            fmt.Printf("hesaba giris yapildi %s\n", token)
        case 401:
            fmt.Printf("hesaba giris yapilamadi %s - Geçersiz token\n", token)
        case 403:
            if strings.Contains(responseText, "2FA") {
                fmt.Printf("giris basarisiz %s - İki aşamalı doğrulama gerekli\n", token)
            } else if strings.Contains(responseText, "robot") {
                fmt.Printf("giris basarisiz %s - Robot doğrulaması gerekli\n", token)
            } else {
                fmt.Printf("giris basarisiz %s - Erişim reddedildi\n", token)
            }
        default:
            fmt.Printf("Giriş başarısız: %s - Durum Kodu: %d\n", token, resp.StatusCode)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("Dosya okuma hatası: %v", err)
    }
}
