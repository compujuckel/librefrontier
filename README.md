# LibreFrontier
This aims to be a replacement for the proprietary Internet Radio backend of devices with Silicon Frontier chipsets (e.g. Technisat, Teufel, Hama, Myine, many more).  
Radio data is sourced from http://www.radio-browser.info/ API.

## Why?
Silicon Frontier switched their radio data provider in May 2019 which made the Internet Radio function of these devices unusable for several days ([German news article](https://www.heise.de/newsticker/meldung/Massenhafter-Ausfall-von-Internetradios-4417248.html)).  
**All favorited and custom stations were lost.**

## Features
- [x] Stations by country (pagination still buggy)
- [x] Most popular stations
- [x] Most liked stations
- [x] Search
- [ ] Favorite stations
- [ ] Custom stations
- ... more

## How does this work?
Silicon Frontier devices talk to their vendor-specific backend at wifiradiofrontier.com (e.g. teufel.wifiradiofrontier.com) via a simple, unencrypted HTTP API.  
To use LibreFrontier, create a DNS entry for this domain in your router and point it to the LibreFrontier instance. A public instance will probably come soon(tm)

**This project is work in progress and only tested on a Teufel Radio 3sixty.**
