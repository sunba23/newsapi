from clients.api_client import ApiClient
from clients.db_client import DbClient
from config import Config

from clients.api_client import NewsApiClient
from clients.db_client import PgClient
from models import News

from dotenv import load_dotenv


class App:
    def __init__(
        self, config: Config, db_client: DbClient, api_clients: list[ApiClient]
    ) -> None:
        self.config = config
        self.db_client = db_client
        self.api_clients = api_clients

    def fetch(self) -> list[News]:
        tags = self.db_client.get_tags()
        news: list[News] = []
        for api_client in self.api_clients:
            news.extend(api_client.get_news(tags))
        return news

    def save(self, news) -> bool:
        return self.db_client.save_news(news)


def main() -> None:
    load_dotenv()
    config = Config()
    config.set_defaults()

    api_clients: list[ApiClient] = [NewsApiClient(config=config)]
    pg_client = PgClient(conn_str=config.conn_str)
    app = App(config=config, db_client=pg_client, api_clients=api_clients)
    news: list[News] = app.fetch()
    app.save(news)


if __name__ == "__main__":
    main()
