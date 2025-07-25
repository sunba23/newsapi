from abc import ABC, abstractmethod
from enum import Enum

from models import News, Tag

import psycopg
from psycopg.sql import SQL, Literal
from psycopg import rows


class QueryType(Enum):
    READ = 1
    WRITE = 2


class DbClient(ABC):
    @abstractmethod
    def run_query(self, query: SQL, type: QueryType) -> list[rows.TupleRow] | int:
        pass

    @abstractmethod
    def get_tags(self) -> list[Tag]:
        pass

    @abstractmethod
    def save_news(self, news: list[News]) -> int:
        pass


class PgClient(DbClient):
    def __init__(self, conn_str: str) -> None:
        self.conn_str = conn_str

    def get_connection(self) -> psycopg.Connection:
        return psycopg.connect(conninfo=self.conn_str, row_factory=rows.tuple_row)

    def run_query(self, query: SQL, type: QueryType) -> list[rows.TupleRow] | int:
        with self.get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(query)
                if type is QueryType.READ:
                    return cur.fetchall()
                elif type is QueryType.WRITE:
                    conn.commit()
                    return cur.rowcount

    def get_tags(self) -> list[Tag]:
        row_list: list[rows.TupleRow] = self.run_query(
            SQL("SELECT name FROM tags"), QueryType.READ
        )
        return [Tag(name=str(row[0])) for row in row_list]

    def save_news(self, news: list[News]) -> int:
        if not news:
            return 0

        total_news_added = 0

        for article in news:
            query = SQL(
                """
                WITH inserted_news AS (
                    INSERT INTO news (title, content, author, created_at)
                    VALUES ({title}, {content}, {author}, CURRENT_TIMESTAMP)
                    RETURNING id
                ),
                matched_tags AS (
                    SELECT id FROM tags 
                    WHERE name = ANY({tags}::text[])
                )
                INSERT INTO news_tags (news_id, tag_id)
                SELECT n.id, t.id
                FROM inserted_news n
                CROSS JOIN matched_tags t
                ON CONFLICT DO NOTHING
            """
            ).format(
                title=Literal(article.title),
                content=Literal(article.content),
                author=Literal(article.author),
                tags=Literal(article.tags),
            )

            count = self.run_query(query, QueryType.WRITE)
            if count > 0:
                total_news_added += 1

        return total_news_added
