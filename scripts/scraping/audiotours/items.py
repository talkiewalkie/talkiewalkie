# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class AudiotoursItem(scrapy.Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    title = scrapy.Field()
    description = scrapy.Field()
    author = scrapy.Field()
    audio_length = scrapy.Field()
    audio_url = scrapy.Field()
    tour_url = scrapy.Field()
