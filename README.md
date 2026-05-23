# tuigram

A terminal UI client for Telegram. This is just an early version for my personal use only. Feel free to fork and tweak however you like to suit your needs.

## Installation

Download the latest release for your platform from the [releases page](https://github.com/sangnd2x/tuigram/releases).

```bash
# Linux amd64
curl -L https://github.com/sangnd2x/tuigram/releases/latest/download/tuigram_<version>_linux_amd64.tar.gz | tar -xz
sudo mv tui-telegram /usr/local/bin/
```

```bash
# Linux arm64
curl -L https://github.com/sangnd2x/tuigram/releases/latest/download/tuigram_<version>_linux_arm64.tar.gz | tar -xz
sudo mv tui-telegram /usr/local/bin/
```

Replace `<version>` with the latest version number (e.g. `0.1.0`).

## Setup

tuigram uses the Telegram API directly and requires your own API credentials.

**Get your API credentials:**

1. Go to [https://my.telegram.org](https://my.telegram.org) and log in with your phone number
2. Click **API development tools**
3. Create a new application — the name and description can be anything
4. Copy your **App api_id** and **App api_hash**

**First run:**

```bash
tui-telegram
```

On first launch, tuigram will prompt you for your API ID and hash, then save them to `~/.config/tui-telegram/config.toml`. You won't be asked again.

## Configuration

Config is stored at `~/.config/tui-telegram/config.toml`:

```toml
api_id = 12345678
api_hash = "your_api_hash_here"
```

You can edit this file directly if you need to update your credentials.

## Usage

```bash
tui-telegram
```

## License

MIT
