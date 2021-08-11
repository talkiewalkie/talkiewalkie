# scraping audio tours website

## mywowo.net

Need a service account file with bucket write access, then

```
GOOGLE_APPLICATION_CREDENTIALS='./.secrets/talkiewalkie-305117-1a6c005507ac.json' scrapy crawl mywowo -o mywowo-w-upload.jl
```

However this command wouldn't need to be run again, as results are available
at `gs://talkiewalkie-dev/scraping/audiotours/mywowo-full-run.jl` and all files are reuploaded with
prefix `gs://talkiewalkie-dev/scraping/audiotours`.

In the `.jl` file are items like this:

```json
{
  "title": "MUSEUM OF WESTERN ART, Rodin Lobby",
  "description": "In the lobby of the museum, you can see the 1882-83 version of The Thinker, perhaps Rodin\u2019s best-known work and one of the most widely admired worldwide. There is more than just one copy of the sculpture; a number of bronze castings were made from the clay model. The one you can admire here dates to 1904.\r\nThe statue depicts the same image of a man deep in thought you already saw on The Gates of Hell. Although the sculpture does not actually look like Dante, Rodin\u2019s intent was to represent the author of the Divine Comedy, and the statue was initially known as \u201cThe Poet\u201d.\r\nThe figure, deep in thought, is nude, because Rodin\u2019s works were inspired by classical Greek sculpture, which idealized a balance between interior and exterior beauty. The specific inspiration for this statue, however, which many see as a representation of philosophy, came from Il Pensieroso, the figure sculpted by Michelangelo on the tomb of Lorenzo de\u2019 Medici in the Basilica of San Lorenzo in Florence.\r\nThe museum is home to a total of fourteen statues by Rodin. You might also like to take a look at the one of Saint John the Baptist. The model for the statue was an Italian peasant who visited the artist\u2019s workshop: at first sight, his striking bearing and physique and the mystical expression on his face prompted Rodin to think of the Saint.\r\nWhat is unusual about this statue, which shows the naturally muscular frame of a man accustomed to toil in the fields, is the position of the legs: spread open as if the figure were walking, yet straight and with both feet on the ground, in an unnatural pose for a man in movement. Rodin explained that his intention was not to show the man walking, but the impulse to do so.\r\nAnother fine work by Rodin exhibited here is the sculpture entitled \"Fugit amor\u201d, Latin for \u201cfugitive love\u201d. Though it is charged with eroticism and sensuality, the statue from 1887 portrays the two lovers in entirely unnatural positions, with the woman lying down and the man slipping off her back, and its aim is to show the fleeting nature and rapid end of carnal love. You may recall that this figure also appears among those on The Gates of Hell. \u00a0You can compare it with a similarly erotic statue entitled \u201cI am beautiful\u201d, showing a nude man raising a woman in his arms, representing the bond between the two sexes.\r\nAn interesting fact: Rodin was so enamored of The Thinker that he wanted a copy placed on his tomb.",
  "audio_length": "2.47",
  "author": null,
  "audio_url": "https://mywowo.net/media/audio/demo/11200_30_baa36e7f63755e4ac8a651976871dd52.mp3",
  "tour_url": "https://mywowo.net/en/japan/tokyo/museum-of-western-art/rodin-lobby",
  "cover_url": "https://mywowo.net/media/images/cache/tokyo_museum_western_art_04_rodin_atrio_jpg_1200_630_cover_85.jpg",
  "file_urls": [
    "https://mywowo.net/media/audio/demo/11200_30_baa36e7f63755e4ac8a651976871dd52.mp3",
    "https://mywowo.net/media/images/cache/tokyo_museum_western_art_04_rodin_atrio_jpg_1200_630_cover_85.jpg"
  ],
  "files": [
    {
      "url": "https://mywowo.net/media/audio/demo/11200_30_baa36e7f63755e4ac8a651976871dd52.mp3",
      "path": "full/4d50825958277a16816975e49b2d4d19e3b991a5.mp3",
      "checksum": "9240addee6cf96bfd986e69ccd840840",
      "status": "downloaded"
    },
    {
      "url": "https://mywowo.net/media/images/cache/tokyo_museum_western_art_04_rodin_atrio_jpg_1200_630_cover_85.jpg",
      "path": "full/31de6883e0279a00aff50e7dfa1086c811bde581.jpg",
      "checksum": "f2b63bf5a3983f78ec74b052c1fb3a05",
      "status": "downloaded"
    }
  ]
}
```