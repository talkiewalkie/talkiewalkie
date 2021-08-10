import re

import scrapy

from audiotours.items import AudiotoursItem


class MywowoSpider(scrapy.Spider):
    name = "mywowo"
    allowed_domains = ["mywowo.net"]
    start_urls = ["https://mywowo.net/en/available-content"]

    def parse(self, response):
        cities = response.css(".cities-list-item::attr(href)").getall()
        for city in cities:
            if city is None or not city.startswith("https://mywowo.net/en"):
                continue
            yield scrapy.Request(city, callback=self.parse_city_page)
            break

    def parse_city_page(self, response):
        tours = response.css(".audioguide-list-item::attr(href)").getall()
        for tour in tours:
            if tour is None or not tour.startswith("https://mywowo.net/en"):
                continue
            yield scrapy.Request(tour, callback=self.parse_tour_page)

    def parse_tour_page(self, response):
        tours = response.css(".audioguide-list-item::attr(href)").getall()
        for tour in tours:
            if tour is None or not tour.startswith("https://mywowo.net/en"):
                continue
            yield scrapy.Request(tour, callback=self.parse_single_tour_page)

    def parse_single_tour_page(self, response):
        title = response.css(".page-small-title::text").get()
        description = "".join(response.css(".audioguide__description *::text").getall()).strip()

        infos_soup = "".join(response.css(".audioguide__info > .row > div *::text").getall())
        lengths = re.findall(r"Audio File length: ([^\n]+)", infos_soup)
        length = lengths[0] if len(lengths) > 0 else None
        authors = re.findall(r"Author: ([^\n]+)", infos_soup)
        author = authors[0] if len(authors) > 0 else None

        audio_url = response.css(".audioguide__player::attr(data-audio)").get()
        cover_url = response.urljoin(response.css(".audioguide__wrapper > img ::attr(src)").get())

        yield AudiotoursItem(
            dict(
                title=title,
                description=description,
                audio_length=length,
                author=author,
                audio_url=audio_url,
                tour_url=response.url,
                cover_url=cover_url,
            )
        )
