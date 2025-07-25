class Tag:
    def __init__(self, name: str) -> None:
        self.name = name


class News:
    def __init__(
        self, title: str, content: str, author: str, url: str, tags: list[str] = None
    ) -> None:
        self.title = title
        self.content = content
        self.author = author
        self.url = url
        self.tags = tags or []
