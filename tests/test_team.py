import pytest
import httpx

from conftest import get_random_name


def test_team_add_and_get_success(client: httpx.Client):
    team_name = get_random_name()
    members = [
        {"user_id": "u1", "username": "user1", "is_active": True},
        {"user_id": "u2", "username": "user2", "is_active": True},
    ]
    team_data = {"team_name": team_name, "members": members}

    response = client.post("/team/add", json=team_data)
    assert response.status_code == 201
    created_team = response.json().get("team")
    assert created_team["team_name"] == team_name
    assert len(created_team["members"]) == 2

    response = client.get("/team/get", params={"team_name": team_name})
    assert response.status_code == 200
    retrieved_team = response.json()
    assert retrieved_team["team_name"] == team_name
    assert any(m["user_id"] == "u1" for m in retrieved_team["members"])


def test_team_add_team_exists(client: httpx.Client):
    team_name = get_random_name()
    members = [{"user_id": "u1", "username": "user3", "is_active": True}]
    client.post("/team/add", json={"team_name": team_name, "members": members})
    response = client.post(
        "/team/add", json={"team_name": team_name, "members": members})

    assert response.status_code == 400
    assert response.json()["error"]["code"] == "TEAM_EXISTS"


def test_team_get_not_found(client: httpx.Client):
    response = client.get(
        "/team/get", params={"team_name": get_random_name()})
    assert response.status_code == 404
    assert response.json()["error"]["code"] == "NOT_FOUND"
