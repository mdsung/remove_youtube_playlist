# Remove all items in playlist in Youtube

- Author: MinDong Sung
- Date: 2022-11-21

---

## Objective

- There is no way to remove all items in playlist in Youtube.
- To apply zettelkasten, I need to remove all items in the watch-later playlist in Youtube.
- This script is to remove all items in the watch-later playlist in Youtube.

## Pre-requisite

### Playlist Id

    - To remove watch-later playlist in Youtube, you need to know your watch-later playlist id.
    - Youtube Data API does not provide a way to get the watch-later playlist id.
    - You need to get the watch-later playlist id by yourself.
    - To get the watch-later playlist id, you need to download your own takeout data from Google(https://takeout.google.com/).
    - Select youtube data and download it.
    - You can find `Watch later.csv` file in playlist folder.
    - You can find the watch-later playlist id in `Watch later.csv` file. It starts with 'PL\*'.
    - Set your .env file with `PLAYLIST_ID=PL*`.

### Client Secret

    - You need to know the playlist ID of the playlist you want to remove all items.
    - The Youtube Data API v3 is enabled in your Google Cloud Platform.
    - You need to create a service account and download the JSON file.(`client_secret.json`)

### golang

    - You need to install golang.
    - https://golang.org/doc/install

## Process of the project

- Start your project

```
mkdir work/2022_youtube_api
go mod init 2022_youtube_api
```

- You need to install the Google API Client Library for Golang.

```
go get google.golang.org/api/youtube/v3
go get golang.org/x/oauth2/...
```

- Run script

```
go run remove_all_watch_later.go
```
