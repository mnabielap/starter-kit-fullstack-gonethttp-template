import sys
import os
sys.path.append(os.path.abspath(os.path.dirname(__file__)))
import utils

token = utils.load_config("access_token")
headers = {
    "Authorization": f"Bearer {token}"
}

# Test query params: page 1, limit 5, sort by name ascending
params = "?page=1&limit=5&sortBy=name:asc"

utils.send_and_print(
    url=f"{utils.BASE_URL}/users{params}",
    method="GET",
    headers=headers,
    output_file=f"{os.path.splitext(os.path.basename(__file__))[0]}.json"
)