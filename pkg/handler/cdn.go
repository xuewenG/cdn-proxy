package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xuewenG/cdn-proxy/pkg/config"
	"golang.org/x/net/proxy"
)

type FileLock struct {
	mu    sync.Mutex
	group sync.WaitGroup
}

var fileLocks sync.Map

func getFileLock(path string) *FileLock {
	lock, _ := fileLocks.LoadOrStore(path, &FileLock{})
	return lock.(*FileLock)
}

func CdnProxy(c *gin.Context) {
	cdnName := c.Param("cdnName")
	resourcePath := c.Param("resourcePath")
	log.Printf("Receive request, cdnName: %s, resourcePath: %s", cdnName, resourcePath)

	if strings.Contains(resourcePath, "..") {
		log.Printf("Unsafe resourcePath detected: %s", resourcePath)
		c.Status(http.StatusBadRequest)
		return
	}

	var cdnUrl string
	for _, cdn := range config.Config.Cdn {
		if cdn.Name == cdnName {
			cdnUrl = cdn.Url
			break
		}
	}

	log.Printf("Found cdnUrl: %s", cdnUrl)
	if cdnUrl == "" {
		log.Printf("CDN not found: %s", cdnName)
		c.Status(http.StatusNotFound)
		return
	}

	cacheFilePath := filepath.Join(config.Config.CacheDir, cdnName, resourcePath)
	cacheFileLock := getFileLock(cacheFilePath)
	if _, err := os.Stat(cacheFilePath); err == nil {
		log.Printf("Serving from cache: %s", cacheFilePath)
		c.File(cacheFilePath)
		return
	}

	if !cacheFileLock.mu.TryLock() {
		log.Printf("Waiting for file lock: %s", cacheFilePath)
		cacheFileLock.group.Wait()
		c.File(cacheFilePath)
		return
	}
	defer cacheFileLock.mu.Unlock()

	cacheFileLock.group.Add(1)
	defer cacheFileLock.group.Done()

	targetUrl := fmt.Sprintf("%s%s", cdnUrl, resourcePath)
	log.Printf("Requesting Original CDN: %s", targetUrl)

	var client *http.Client
	if config.Config.SocksProxyUrl == "" {
		client = &http.Client{}
	} else {
		log.Printf("Using SOCKS proxy: %s", config.Config.SocksProxyUrl)

		proxyURL, err := url.Parse(config.Config.SocksProxyUrl)
		if err != nil {
			log.Printf("Error parsing SOCKS proxy URL: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			log.Printf("Error creating SOCKS dialer: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		transport := &http.Transport{
			Dial: dialer.Dial,
		}
		client = &http.Client{Transport: transport}
	}

	resp, err := client.Get(targetUrl)
	if err != nil {
		log.Printf("Error requesting Url %s: %v", targetUrl, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Response status: %d for Url: %s", resp.StatusCode, targetUrl)
	if resp.StatusCode != http.StatusOK {
		c.Status(resp.StatusCode)
		return
	}

	os.MkdirAll(filepath.Dir(cacheFilePath), 0755)

	file, err := os.Create(cacheFilePath)
	if err != nil {
		log.Printf("Error creating cache file %s: %v", cacheFilePath, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Printf("Error copying response to cache file %s: %v", cacheFilePath, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully cached file: %s", cacheFilePath)
	c.File(cacheFilePath)
}
