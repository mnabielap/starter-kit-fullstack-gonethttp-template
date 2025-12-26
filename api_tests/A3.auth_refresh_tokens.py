import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

refresh_token = utils.load_config("refresh_token")

if not refresh_token:
    print("No refresh token found. Run A1 or A2 first.")
    sys.exit(1)

body = {
    "refreshToken": refresh_token
}

response = utils.send_and_print(
    url=f"{utils.BASE_URL}/auth/refresh-tokens",
    method="POST",
    body=body,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)

if response.status_code == 200:
    data = response.json()
    root = data.get("data", data)
    
    if "access" in root:
        utils.save_config("access_token", root["access"]["token"])
        utils.save_config("refresh_token", root["refresh"]["token"])
        print("\n[INFO] Tokens refreshed successfully")