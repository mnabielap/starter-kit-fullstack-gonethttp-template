import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

token = utils.load_config("access_token")
target_id = utils.load_config("target_user_id")

if not target_id:
    print("Target User ID not found. Run B1 first.")
    sys.exit(1)

headers = {
    "Authorization": f"Bearer {token}"
}

utils.send_and_print(
    url=f"{utils.BASE_URL}/users/{target_id}",
    method="GET",
    headers=headers,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)