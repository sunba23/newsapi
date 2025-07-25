from abc import ABC, abstractmethod
from datetime import datetime
from dateutil.relativedelta import relativedelta

from config import Config
from models import News, Tag

from newsapi import NewsApiClient as NewsApiClientC


class ApiClient(ABC):
    @abstractmethod
    def get_news(self, tags: list[Tag]) -> list[News]:
        pass


class NewsApiClient(ApiClient):
    def __init__(self, config: Config) -> None:
        self.client = NewsApiClientC(api_key=config.api_key)

    def get_news(self, tags: list[Tag]) -> list[News]:
        news_map = {}

        for tag in tags:
            news_entries = self.get_news_for_tag(tag)
            for news in news_entries:
                key = news.url
                if key in news_map:
                    if tag.name not in news_map[key].tags:
                        news_map[key].tags.append(tag.name)
                else:
                    news.tags = [tag.name]
                    news_map[key] = news

        return list(news_map.values())

    def get_news_for_tag(self, tag: Tag) -> list[News]:
        query = tag.name
        from_param = (datetime.today() - relativedelta(months=1)).strftime("%Y-%m-%d")
        response = self.client.get_everything(
            q=query,
            from_param=from_param,
            language="en",
            sort_by="publishedAt",
            page=2,
        )
        news_entries = response.get("articles", [])
        return self.parse_news(news_entries, tag)

    @staticmethod
    def parse_news(news_entries, tag: Tag) -> list[News]:
        news_list = []
        for news in news_entries:
            title = news.get("title", "")
            content = news.get("content", "") or news.get("description", "")
            author = news.get("author", "") or "Unknown"
            url = news.get("url", "")

            news_item = News(title=title, content=content, author=author, url=url)
            news_list.append(news_item)
        return news_list
