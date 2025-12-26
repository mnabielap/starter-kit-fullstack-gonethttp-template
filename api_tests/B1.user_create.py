import sys
import os
import time
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

# Ensure we are logged in as Admin (usually the seed user 'admin@example.com' is admin)
# You might need to run A2.auth_login.py with admin creds first if A1 created a normal user.
token = utils.load_config("access_token")

unique = int(time.time())
body = {
    "name": f"New User {unique}",
    "email": f"newuser_{unique}@example.com",
    "password": "password123",
    "role": "user"
}

headers = {
    "Authorization": f"Bearer {token}"
}

response = utils.send_and_print(
    url=f"{utils.BASE_URL}/users",
    method="POST",
    headers=headers,
    body=body,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)

if response.status_code == 201:
    data = response.json()
    root = data.get("data", data)
    # Save this user ID for B3, B4, B5
    utils.save_config("target_user_id", root["id"])
    print(f"\n[INFO] Created user ID: {root['id']}")