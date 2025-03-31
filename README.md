# 🚀 go-rate-limit: A Simple & Extensible Rate Limiting Library for Go

![Go Version](https://img.shields.io/github/go-mod/go-version/sirius1b/go-rate-limit)
![License](https://img.shields.io/github/license/sirius1b/go-rate-limit)
![Issues](https://img.shields.io/github/issues/sirius1b/go-rate-limit)
![Stars](https://img.shields.io/github/stars/sirius1b/go-rate-limit?style=social)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)
[![Contributions Welcome](https://img.shields.io/badge/Contributions-Welcome-ff69b4.svg)](https://github.com/sirius1b/go-rate-limit/issues)

---

## ⚡ Overview

`go-rate-limit` is a lightweight and easy-to-use rate limiting library for Go.
Currently, it supports traditional rate limiting, with more algorithms coming soon! 🚀

- ✅ **Simple API** – Minimal setup required.
- ✅ **Thread-safe** – Designed for concurrent use.
- ✅ **Extensible** – More rate limiting algorithms will be added.
- ✅ **Extensible** – More Caching Options will be added (REDIS) - Planned
- ✅ **Open for Contributions** – PRs are welcome! 🎉

---

## 📦 Installation

```sh
go get github.com/sirius1b/go-rate-limit
```

## 📦 Usage

#### Fixed Window

```go


	limiter, err := Require(FixedWindow, Options{
		Limit:  10,
		Window: time.Second,
	})
    token_id = "fe8db250-f6fe-4a88-940d-56fbe8892876"  // USER_ID

    limiter.Allow(token_id) // will ALLOW/REJECT request as per rate-limit capacity
    limiter.Wait(token_id) // hold till resource not available

```

Checkout in examples!

## 🛠 Features

- Fixed Window Rate Limiting – Available now!
- Sliding Window Rate Limiting – Available
- Token Bucket – Available !
- Distributed Rate Limiting – Redis-based implementation (Planned).

## 🤝 Contributing

We welcome contributions! 🎉

- Open an Issue for bug reports & feature requests.
- Fork the repo and submit a Pull Request (PR).

###### 💡 New ideas? Let’s discuss! 🗨️

## 📄 License

This project is licensed under the MIT License.

### 🚀 Happy coding! 💙
