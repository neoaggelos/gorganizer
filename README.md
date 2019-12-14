# gorganizer

Group up old photos and videos by year and month, removing duplicates.

## Requirements

- Go

## Build

`gorganizer` uses only the Go standard library, so building is as simple as:

```bash
$ go build gorganizer.go -o gorganizer
```

## What

- Automatically organize old photos and videos by year and month, based on
  their modification timestamps.
- Automatically catch and drop duplicate files (using MD5 hash)
- Preserve the modification timestamps of all files.
- `gorganizer` will only copy files, it will not remove anything.

## Why

- Because I had a huge folder of old photographs and videos lying around, and
  I wanted to finally fix that mess.
- In order to practice my skills with Go.
- Because I can.

## Example

Assume you have all your old photos at `/mnt/HD/OLD-PHOTOS`. You bought a shiny
new HD, and want to move your organized files at `/mnt/MY-SHINY-NEW-HD`

Do this with:

```bash
$ ./gorganizer -source /mnt/media/HD/OLD_PHOTOS/ -dest /mnt/media/MY-SHINY-NEW-HD
```


~ Aggelos Kolaitis
