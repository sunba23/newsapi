import os
from datetime import timedelta


class Config:
    def __init__(self) -> None:
        self.conn_str = os.getenv("POSTGRES_CONN_STR", "")
        self.api_key = os.getenv("NEWSAPI_KEY", "")

        if not self.conn_str:
            raise ValueError("POSTGRES_CONN_STR environment variable is not set")
        if not self.api_key:
            raise ValueError("NEWSAPI_KEY environment variable is not set")

        self.set_defaults()

    def set_defaults(self) -> None:
        self.timeout: timedelta = timedelta(minutes=10)
