import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

body = {
    "email": "admin@example.com",
    "password": "password123"
}

response = utils.send_and_print(
    url=f"{utils.BASE_URL}/auth/login",
    method="POST",
    body=body,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)

if response.status_code == 200:
    data = response.json()
    root = data.get("data", data)
    
    if "tokens" in root:
        utils.save_config("access_token", root["tokens"]["access"]["token"])
        utils.save_config("refresh_token", root["tokens"]["refresh"]["token"])
        print("\n[INFO] Tokens refreshed via Login")