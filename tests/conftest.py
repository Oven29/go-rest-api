import uuid
import pytest
import httpx


@pytest.fixture(scope="session")
def base_url():
    return "http://localhost:8080"


@pytest.fixture()
def client(base_url):
    with httpx.Client(base_url=base_url) as c:
        yield c


def get_random_name():
    return str(uuid.uuid4())
