# 🚀 RateLimit-Go: A Simple & Extensible Rate Limiting Library for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/ratelimit-go)](https://goreportcard.com/report/github.com/yourusername/ratelimit-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/ratelimit-go.svg)](https://pkg.go.dev/github.com/yourusername/ratelimit-go)
[![Contributions Welcome](https://img.shields.io/badge/Contributions-Welcome-ff69b4.svg)](https://github.com/yourusername/ratelimit-go/issues)

---

## ⚡ Overview

`go-rate-limit` is a lightweight and easy-to-use rate limiting library for Go.
Currently, it supports **Fixed Window** rate limiting, with more algorithms coming soon! 🚀

- ✅ **Simple API** – Minimal setup required.
- ✅ **Thread-safe** – Designed for concurrent use.
- ✅ **Extensible** – More rate limiting algorithms will be added.
- ✅ **Extensible** – More Caching Options will be added (REDIS)
- ✅ **Open for Contributions** – PRs are welcome! 🎉

---

## 📦 Installation

```sh
go get github.com/sirius1b/go-rate-limit
```

## 📦 Usage

```go


	limiter, err := Require(FixedWindow, Options{
		Limit:  10,
		Window: time.Second,
	})
    token_id = "fe8db250-f6fe-4a88-940d-56fbe8892876"  // USER_ID

    limiter.Allow(token_id) // will ALLOW/REJECT request as per rate-limit capacity
    limiter.Wait(token_id) // hold till resource not available

```

## 🛠 Features

- Fixed Window Rate Limiting – Available now!
- Sliding Window Rate Limiting – Coming soon!
- Token Bucket – Smooth request handling (Planned).
- Leaky Bucket – Even request distribution (Planned).
- Distributed Rate Limiting – Redis-based implementation (Planned).

## 🤝 Contributing

We welcome contributions! 🎉

- Open an Issue for bug reports & feature requests.
- Fork the repo and submit a Pull Request (PR).

###### 💡 New ideas? Let’s discuss! 🗨️

## 📄 License

This project is licensed under the MIT License.

### 🚀 Happy coding! 💙
