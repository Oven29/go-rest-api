import pytest
import httpx

from conftest import get_random_name


def set_user_activity(client: httpx.Client, user_id: str, is_active: bool):
    return client.post("/users/setIsActive", json={
        "user_id": user_id,
        "is_active": is_active
    })


def test_user_set_is_active_success(client: httpx.Client):
    user_id = "u1"
    members = [{"user_id": user_id, "username": "user1", "is_active": True}]
    client.post("/team/create",
                json={"team_name": get_random_name(), "members": members})

    response = set_user_activity(client, user_id, False)
    assert response.status_code == 200
    updated_user = response.json().get("user")
    assert updated_user["user_id"] == user_id
    assert updated_user["is_active"] is False

    response = set_user_activity(client, user_id, True)
    assert response.status_code == 200
    updated_user = response.json().get("user")
    assert updated_user["is_active"] is True


def test_user_set_is_active_not_found(client: httpx.Client):
    response = set_user_activity(client, "u666", False)
    assert response.status_code == 404
    assert response.json()["error"]["code"] == "NOT_FOUND"
