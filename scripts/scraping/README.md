# scraping audio tours website

## mywowo.net

Need a service account file with bucket write access, then

```
GOOGLE_APPLICATION_CREDENTIALS='./.secrets/talkiewalkie-305117-1a6c005507ac.json' scrapy crawl mywowo -o mywowo-w-upload.jl
```